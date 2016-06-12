package nanairoishi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/bitly/go-simplejson"
	"github.com/spf13/viper"
)

const APPHOME = "/.nanairoishi/"
const CONFIG = APPHOME + "config.yaml"
const HISTORY = APPHOME + "history.json"

type SGConfig struct {
	Name    string
	Profile string
	Region  string
	ID      string
	Port    int64
	IP      string
}

type SGConfigs []SGConfig

func init() {
	home, homeErr := HomeDir()
	if homeErr != nil {
		panic("Can't init User Home directory")
	}
	// apphomeがあれば初期化
	_, err := os.Stat(home + APPHOME)
	if err == nil {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home + APPHOME)
	}
}

func Initialization() string {
	home, err := HomeDir()
	if err != nil {
		panic("Can't init User Home directory")
	}
	// apphomeがなければ作る
	if _, apErr := os.Stat(home + APPHOME); os.IsNotExist(apErr) {
		mkdirErr := os.Mkdir(home+APPHOME, 0700)
		if mkdirErr != nil {
			panic("Can't init app home directory")
		}
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(home + APPHOME)
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		initErr := initConfig(home)
		if initErr != nil {
			panic("Can't init configuration files")
		}
	}

	msg := `application home : '%v'
config file : '%v'
history file : '%v'
`
	return fmt.Sprintf(msg, home+APPHOME, home+CONFIG, home+HISTORY)
}

func HomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func initConfig(home string) error {
	content := []byte("# default config\n")
	err := ioutil.WriteFile(home+CONFIG, content, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(home+HISTORY, []byte("{}"), os.ModePerm)
}

func LoadConfigs() (SGConfigs, error) {
	viperErr := viper.ReadInConfig()
	if viperErr != nil {
		return nil, viperErr
	}

	var configs SGConfigs
	err := viper.UnmarshalKey("Configs", &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func loadHistoryJson() (*simplejson.Json, string, error) {
	home, err := HomeDir()
	if err != nil {
		return nil, "", err
	}
	rf, ioErr := ioutil.ReadFile(home + HISTORY)
	if ioErr != nil {
		return nil, "", ioErr
	}
	root, jsonErr := simplejson.NewJson(rf)
	if jsonErr != nil {
		return nil, "", jsonErr
	}
	return root, home, nil
}

func GetHistory(name string) (SGConfig, error) {
	var config SGConfig
	root, _, jsonErr := loadHistoryJson()
	if jsonErr != nil {
		return config, jsonErr
	}

	body := root.Get(name).MustString("")
	// 空文字なら履歴はない
	if body == "" {
		config.Name = ""
		return config, nil
	}

	unErr := json.Unmarshal([]byte(body), &config)
	return config, unErr
}

func SaveHistory(config SGConfig) error {
	root, home, jsonErr := loadHistoryJson()
	if jsonErr != nil {
		return jsonErr
	}

	// 値を上書き or 新規作成
	history, _ := json.Marshal(config)
	root.Set(config.Name, string(history))

	w, openErr := os.Create(home + HISTORY)
	if openErr != nil {
		return openErr
	}
	defer w.Close()
	o, encErr := root.EncodePretty()
	if encErr != nil {
		return encErr
	}
	w.Write(o)

	return nil
}
