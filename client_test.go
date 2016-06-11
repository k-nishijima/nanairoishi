package nanairoishi

import (
	"testing"

	"github.com/k-nishijima/nanairoishi"
)

func TestGetMyIP(t *testing.T) {
	ip, err := nanairoishi.GetMyIP()
	if err != nil {
		t.Error(err)
	}
	t.Log(ip)
}

/*
func TestAddRule(t *testing.T) {
	var c nanairoishi.SGConfig
	c.Profile = "kumogata"
	c.Region = "us-west-2"
	c.ID = "sg-aaaaee2c"
	c.Port = 22
	ip, _ := nanairoishi.GetMyIP()
	c.IP = ip

	err := nanairoishi.AddRule(true, c)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveRule(t *testing.T) {
	var c nanairoishi.SGConfig
	c.Profile = "gkumogata"
	c.Region = "us-west-2"
	c.ID = "sg-aaaaee2c"
	c.Port = 22
	ip, _ := nanairoishi.GetMyIP()
	c.IP = ip

	err := nanairoishi.RemoveRule(false, c)
	if err != nil {
		t.Error(err)
	}
}
*/
