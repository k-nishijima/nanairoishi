package nanairoishi

import (
	"testing"

	"github.com/k-nishijima/nanairoishi"
)

func TestLoadConfig(t *testing.T) {
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
	config.Profile = "foobarProf"
	config.Region = "us-west-2"
	config.ID = "sg-aaaaaaaa"
	config.IP = "127.0.0.1"
	config.Port = 8080

	err := nanairoishi.SaveHistory(config)
	if err != nil {
		t.Error(err)
	}
}

func TestGetHistory(t *testing.T) {
	c, err := nanairoishi.GetHistory("testName")
	if err != nil {
		t.Error(err)
	}
	t.Log(c.Profile)
	t.Log(c.ID)
	t.Log(c.IP)
}
