package nanairoishi

import (
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
	json, jsonErr := simplejson.NewJson(rf)
	if jsonErr != nil {
		return nil, "", jsonErr
	}
	return json, home, nil
}

func GetHistory(name string) (string, error) {
	json, _, jsonErr := loadHistoryJson()
	if jsonErr != nil {
		return "", jsonErr
	}
	return json.Get(name).String()
}

func SaveHistory(config SGConfig) error {
	json, home, jsonErr := loadHistoryJson()
	if jsonErr != nil {
		return jsonErr
	}

	// 値を上書き or 新規作成
	json.Set(config.Name, config.IP)

	w, openErr := os.Create(home + HISTORY)
	if openErr != nil {
		return openErr
	}
	defer w.Close()
	o, encErr := json.EncodePretty()
	if encErr != nil {
		return encErr
	}
	w.Write(o)

	return nil
}
