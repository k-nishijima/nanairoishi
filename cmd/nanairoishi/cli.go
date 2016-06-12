package main

import (
	"fmt"
	"io"

	"github.com/jteeuwen/go-pkg-optarg"
	"github.com/k-nishijima/nanairoishi"
)

const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
	HELP = `
description:

	nanairoishi is maintenance tool for AWS security group.

commands are:

	init   init config files
	update try update security group (you must specify the config name)

config file example:

----
Configs:
- Name: myProject1
  Profile: aws_profile_name
  Region: us-west-2
  ID: sg-aaaabbbb
  Port: 8080

- Name: myProject2
  Profile: aws_profile_name
  Region: us-west-2
  ID: sg-ccccdddd
  Port: 8081
----
`
)

type CLI struct {
	outStream, errStream io.Writer
}

func (c *CLI) Run() int {
	optarg.Header("General options")
	optarg.Add("h", "help", "Displays this help.", false)
	optarg.Add("v", "version", "Displays version information.", false)

	optarg.Header("init : initialize application directory and config files")
	optarg.Header("update : security group by argument name")
	optarg.Add("n", "name", "you must specify the config name", "")
	optarg.Add("D", "dryRun", "dryRun option", false)

	var version, help, dryRun bool
	var configName string

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "h":
			help = opt.Bool()
		case "v":
			version = opt.Bool()
		case "n":
			configName = opt.String()
		case "D":
			dryRun = opt.Bool()
		}
	}

	if version {
		fmt.Fprintf(c.errStream, "nanairoish version %s\n", Version)
		return ExitCodeOK
	}
	if help {
		optarg.Usage()
		fmt.Fprintln(c.outStream, HELP)
		return ExitCodeOK
	}

	if len(optarg.Remainder) == 1 {
		cmd := optarg.Remainder[0]
		switch cmd {
		case "init":
			fmt.Fprintln(c.outStream, "current config files are:")
			fmt.Fprint(c.outStream, nanairoishi.Initialization())
			return ExitCodeOK
		case "update":
			return cmdUpdate(c, dryRun, configName)
		}
	} else {
		optarg.Usage()
	}

	return ExitCodeOK
}

func cmdUpdate(c *CLI, dryRun bool, configName string) int {
	fmt.Fprintf(c.outStream, "you specified the security group name called '%v'\n", configName)

	// // ヒストリがあればそれを削除
	// historyIP, getHistoryErr := nanairoishi.GetHistory(configName)
	// if getHistoryErr != nil {
	// 	fmt.Fprintf(c.errStream, "GetHistory failed\n")
	// 	return ExitCodeParseFlagError
	// }
	// if historyIP != "" {
	// 	// var c nanairoishi.SGConfig
	// 	// c.Profile = "gkumogata"
	// 	// c.Region = "us-west-2"
	// 	// c.ID = "sg-aaaaee2c"
	// 	// c.Port = 22
	// 	// ip, _ := nanairoishi.GetMyIP()
	// 	// c.IP = ip
	// 	//
	// 	// nanairoishi.RemoveRule(dryRun, config)
	// }

	return ExitCodeOK
}
