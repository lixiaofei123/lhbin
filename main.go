package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lixiaofei123/lhbin/config"
)

func checkArg(arg *string, errText string) {
	if *arg == "" {
		log.Panicln(errText)
	}
}

func addAccount() {

	var driver string
	var account string // 账号
	var akid string
	var aksecret string

	flag.StringVar(&driver, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定")
	flag.StringVar(&akid, "id", "", "密钥ID")
	flag.StringVar(&aksecret, "key", "", "密钥key")

	flag.CommandLine.Parse(os.Args[3:])

	if driver != string(config.QQCloud) {
		log.Fatalln("目前仅支持腾讯云轻量服务器")
	}

	checkArg(&account, "账号名称不能为空")
	checkArg(&akid, "密钥ID不能为空")
	checkArg(&aksecret, "密钥Key不能为空")

	config.AddAccount(&config.AccountConfig{
		Driver:   config.DriverName(driver),
		Account:  account,
		AKID:     akid,
		AKSecret: aksecret,
	})

	fmt.Println("配置成功")

}

func deleteAccount() {
	var driver string
	var account string // 账号
	flag.StringVar(&driver, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定")
	flag.CommandLine.Parse(os.Args[3:])
	if driver != string(config.QQCloud) {
		log.Fatalln("目前仅支持腾讯云轻量服务器")
	}

	checkArg(&account, "账号名称不能为空")
	config.DeleteAccount(config.DriverName(driver), account)
	fmt.Println("删除配置成功")
}

func configAccountAk() {
	operator := os.Args[2]
	if operator == "add" {
		addAccount()
	}
	if operator == "del" || operator == "delete" {
		deleteAccount()
	}
}

func printHelp() {

}

func main() {

	if len(os.Args) <= 1 {
		printHelp()
	}

	childCommand := os.Args[1]

	if childCommand == "config" {
		// 配置ak
		configAccountAk()
	}

	printHelp()

}
