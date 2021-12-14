package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/lixiaofei123/lhbin/config"
)

const ConfigCommandName string = "config"

func init() {

	RegisterChildCommand(ConfigCommandName, "管理账户信息配置，支持多账户", []string{})
	RegisterChildCommandOperator(ConfigCommandName, "add", "添加新的账户", []string{}, SafeOperation(AddAccount))
	RegisterChildCommandOperator(ConfigCommandName, "del", "删除指定账户", []string{"delete"}, SafeOperation(DeleteAccount))
}

func AddAccount() error {

	var driverName string
	var account string // 账号
	var akid string
	var aksecret string

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定")
	flag.StringVar(&akid, "id", "", "密钥ID")
	flag.StringVar(&aksecret, "key", "", "密钥key")

	flag.CommandLine.Parse(os.Args[3:])

	if driverName != string(config.QQCloud) {
		return errors.New("目前仅支持腾讯云轻量服务器")
	}

	checkArg(&account, "账号名称不能为空")
	checkArg(&akid, "密钥ID不能为空")
	checkArg(&aksecret, "密钥Key不能为空")

	config.AddAccount(&config.AccountConfig{
		Driver:   config.DriverName(driverName),
		Account:  account,
		AKID:     akid,
		AKSecret: aksecret,
	})

	fmt.Println("配置成功")
	return nil
}

func DeleteAccount() error {
	var driverName string
	var account string // 账号
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定")
	flag.CommandLine.Parse(os.Args[3:])
	if driverName != string(config.QQCloud) {
		return errors.New("目前仅支持腾讯云轻量服务器")
	}

	checkArg(&account, "账号名称不能为空")
	config.DeleteAccount(config.DriverName(driverName), account)
	fmt.Println("删除配置成功")

	return nil
}
