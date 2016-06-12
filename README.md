nanairoishi
====

nanairoishi is maintenance tool for AWS security group.

Simply, Add My Global IP to security group and remove old MyIP.

(This tool is my golang learning demo program.)

## Install

```
go get github.com/k-nishijima/nanairoishi/cmd/nanairoishi/
```

## Usage

```
$ nanairoishi -h
Usage: nanairoishi [options]:

[General options]
    --help, -h: Displays this help.
 --version, -v: Displays version information.

[init : initialize application directory and config files]

[update : security group by argument name]
    --name, -n: you must specify the config name
  --dryRun, -D: dryRun option

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
```

## Configuration

AWS Credentials (~/.aws/credentials) and config.yaml.

files path is here.

```
$ nanairoishi init
current config files are:
application home : '/Users/username/.nanairoishi/'
config file : '/Users/username/.nanairoishi/config.yaml'
history file : '/Users/username/.nanairoishi/history.json'
```

Required IAM actions are here.

```
"ec2:AuthorizeSecurityGroupIngress",
"ec2:DescribeSecurityGroups",
"ec2:RevokeSecurityGroupIngress"
```

## Licence

MIT

## Author

[Koichiro Nishijima](https://github.com/k-nishijima/)
