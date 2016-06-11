package nanairoishi

import (
	"testing"

	"github.com/k-nishijima/nanairoishi"
)

func TestLoadConfig(t *testing.T) {
	// h, err := nanairoishi.HomeDir()
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(h)

	// c := nanairoishi.LoadConfig()
	// t.Log(c)
	// // map[string][]string
	// t.Log(c["Name"][0])
	// t.Log(c["Name"][1])

	configs, err := nanairoishi.LoadConfigs()
	if err != nil {
		t.Error(err)
	}

	t.Log(configs)
	for _, v := range configs {
		t.Log(v.Name + " is " + v.ID)
	}

}

func TestSaveHistory(t *testing.T) {
	var config nanairoishi.SGConfig
	config.Name = "testName"
	config.IP = "127.0.0.1"

	err := nanairoishi.SaveHistory(config)
	if err != nil {
		t.Error(err)
	}
}

func TestGetHistory(t *testing.T) {
	name, err := nanairoishi.GetHistory("testName")
	if err != nil {
		t.Error(err)
	}
	t.Log(name)
}
