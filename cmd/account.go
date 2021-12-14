package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/lixiaofei123/lhbin/config"
)

const AccountCommandName string = "account"

func init() {

	RegisterChildCommand(AccountCommandName, "管理账户信息配置，支持多账户", []string{})
	RegisterChildCommandOperator(AccountCommandName, "add", "添加新的账户", []string{}, SafeOperation(AddAccount))
	RegisterChildCommandOperator(AccountCommandName, "del", "删除指定账户", []string{"delete"}, SafeOperation(DeleteAccount))
	RegisterChildCommandOperator(AccountCommandName, "list", "列出所有账户", []string{}, SafeOperation(ListAccounts))
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

	fmt.Printf("配置账户%s成功", account)
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
	fmt.Printf("删除账户%s成功\n", account)

	return nil
}

func ListAccounts() error {

	fmt.Println("------------------------------------------")
	fmt.Println("| 驱动 | 账户名称 | AKID | AKSecret |")
	fmt.Println("------------------------------------------")

	for _, account := range config.GlobalConfig.Accounts {
		fmt.Println("|", account.Driver, "|", account.Account, "|", account.AKID, "|", account.AKSecret, "|")
		fmt.Println("------------------------------------------")
	}

	return nil
}
