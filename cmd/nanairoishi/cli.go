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
	fmt.Fprintf(c.outStream, "you specified the config name called '%v'\n", configName)

	// ヒストリがあればそれを削除
	old, getHistoryErr := nanairoishi.GetHistory(configName)
	if getHistoryErr != nil {
		fmt.Fprintln(c.errStream, getHistoryErr)
		fmt.Fprintf(c.errStream, "GetHistory failed\n")
		return ExitCodeParseFlagError
	} else if old.Name != "" {
		fmt.Fprintf(c.outStream, "try RevokeSecurityGroup : %s\n", configName)
		rmErr := nanairoishi.RemoveRule(dryRun, old)
		if rmErr != nil {
			fmt.Fprintf(c.errStream, "RevokeSecurityGroup failed : %s\n", configName)
			return ExitCodeParseFlagError
		}
		fmt.Fprintln(c.outStream, "OK.")
	} else {
		fmt.Fprintln(c.outStream, "history not found. skip RevokeSecurityGroup")
	}

	// 設定を読み込み
	exist := false
	var config nanairoishi.SGConfig
	configs, loadErr := nanairoishi.LoadConfigs()
	if loadErr != nil {
		fmt.Fprintln(c.errStream, loadErr)
		fmt.Fprintf(c.errStream, "LoadConfigs failed\n")
		return ExitCodeParseFlagError
	}
	for _, v := range configs {
		if configName == v.Name {
			exist = true
			config = v
			break
		}
	}
	if !exist {
		fmt.Fprintf(c.errStream, "not found on your config : %s\n", configName)
		return ExitCodeParseFlagError
	}

	// 現在のIPに書き換えて、ルール追加
	ip, _ := nanairoishi.GetMyIP()
	config.IP = ip
	fmt.Fprintf(c.outStream, "try AuthorizeSecurityGroupIngress : %s\n", configName)
	addErr := nanairoishi.AddRule(dryRun, config)
	if addErr != nil {
		fmt.Fprintf(c.errStream, "AuthorizeSecurityGroupIngress failed : %s\n", configName)
		return ExitCodeParseFlagError
	}
	fmt.Fprintln(c.outStream, "OK.")

	// 履歴として保存
	saveErr := nanairoishi.SaveHistory(config)
	if saveErr != nil {
		fmt.Fprintln(c.errStream, saveErr)
		fmt.Fprintf(c.errStream, "SaveHistory failed\n")
		return ExitCodeParseFlagError
	}

	fmt.Fprintf(c.outStream, "#-----------------------------------------\nupdate successful : %v\n", configName)
	return ExitCodeOK
}
