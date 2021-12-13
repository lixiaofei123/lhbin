package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/lixiaofei123/lhbin/config"
	"github.com/lixiaofei123/lhbin/driver"
)

func checkArg(arg *string, errText string) {
	if *arg == "" {
		log.Panicln(errText)
	}
}

func addAccount() {

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
		log.Fatalln("目前仅支持腾讯云轻量服务器")
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

}

func deleteAccount() {
	var driverName string
	var account string // 账号
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定")
	flag.CommandLine.Parse(os.Args[3:])
	if driverName != string(config.QQCloud) {
		log.Fatalln("目前仅支持腾讯云轻量服务器")
	}

	checkArg(&account, "账号名称不能为空")
	config.DeleteAccount(config.DriverName(driverName), account)
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

func queryRegions() {
	var driverName string
	var account string // 账号
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.CommandLine.Parse(os.Args[2:])

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}
	regions, err := cdriver.ListRegions()
	if err != nil {
		log.Panic(err)
	}

	for _, region := range regions {
		fmt.Println(region.Name, region.Region)
	}

}

func listInstances() {

	var driverName string
	var account string
	var region string
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域，不填写则默认为所有地域")
	flag.CommandLine.Parse(os.Args[3:])

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("------------------------------------------")
	fmt.Println("| 地域 | 实例名称 | 实例ID | 公网IP | 内网IP | 状态 |")
	fmt.Println("------------------------------------------")

	if region == "" {

		regions, err := cdriver.ListRegions()
		if err != nil {
			log.Panic(err)
		}

		for _, region := range regions {
			inss, err := cdriver.ListInstances(region.Region)
			if err != nil {
				log.Panic(err)
			}
			for _, ins := range inss {
				fmt.Println("|", ins.Region, "|", ins.Name, "|", ins.ID, "|", ins.PublicIP, "|", ins.PrivateIP, "|", ins.State, "|")
				fmt.Println("------------------------------------------")
			}
		}

	} else {
		inss, err := cdriver.ListInstances(region)
		if err != nil {
			log.Panic(err)
		}
		for _, ins := range inss {
			fmt.Println("|", ins.Region, "|", ins.Name, "|", ins.ID, "|", ins.PublicIP, "|", ins.PrivateIP, "|", ins.State, "|")
			fmt.Println("------------------------------------------")
		}

	}

	fmt.Print("详细信息可以通过 lhbin ins desc --region region --insids lhins-xxxxx,lhins-yyyyy 命令进行查看")

}

func describeInstance() {

	batchBaseOperatorInstances(func() error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
		insinfo, err := cdriver.InstanceInfo(region, insid)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("-------------------------------")
		fmt.Println("| 地域 | ", insinfo.Region, "|")
		fmt.Println("| ID | ", insinfo.ID, "|")
		fmt.Println("| 名称 | ", insinfo.Name, "|")
		fmt.Println("| 可用区 | ", insinfo.Zone, "|")
		fmt.Println("| 公网IP | ", insinfo.PublicIP, "|")
		fmt.Println("| 内网IP | ", insinfo.PrivateIP, "|")
		fmt.Println("| 带宽 | ", insinfo.Bandwidth, " Mbits|")
		fmt.Println("| CPU核心数 | ", insinfo.Cpu, "|")
		fmt.Println("| 内存 | ", insinfo.Memory, "GB |")
		fmt.Println("| 操作系统 | ", insinfo.OSName, "|")
		fmt.Println("| 状态 | ", insinfo.State, "|")
		fmt.Println("| 创建时间 | ", insinfo.CreatedTime.Format("2006-01-02 15:04:05"), "|")
		fmt.Println("| 过期时间 | ", insinfo.ExpiredTime.Format("2006-01-02 15:04:05"), "|")
		fmt.Println("-------------------------------")
	})

}

func batchBaseOperatorInstances(checkCallback func() error, callback func(cdriver driver.Driver, region, name string, insid string, args ...interface{})) {
	var driverName string
	var account string
	var region string
	var insids string
	var insid string

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "实例所在地域，不填则默认为所有可用区")
	flag.StringVar(&insid, "insid", "", "实例ID，如果设置此值，则会忽略insids参数")
	flag.StringVar(&insids, "insids", "", "实例ID，多个请用逗号隔开。如果不填则默认为所选择可用区下的所有实例")

	if insid != "" {
		insids = insid
	}

	flag.CommandLine.Parse(os.Args[3:])

	err := checkCallback()
	if err != nil {
		log.Panic(err)
	}

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	if region == "" {
		regions, err := cdriver.ListRegions()
		if err != nil {
			log.Panic(err)
		}

		for _, region := range regions {
			if insids == "" {
				inss, err := cdriver.ListInstances(region.Region)
				if err != nil {
					log.Fatal(err)
				}
				for _, ins := range inss {
					callback(cdriver, region.Region, ins.Name, ins.ID)

				}
			} else {
				instanceIDs := strings.Split(insids, ",")

				for _, instanceID := range instanceIDs {
					insinfo, err := cdriver.InstanceInfo(region.Region, instanceID)
					if err != nil {
						fmt.Printf("查询%s地域下的%s信息失败", region.Region, instanceID)
					}
					callback(cdriver, region.Region, insinfo.Name, instanceID)
				}
			}
		}
	} else {

		if insids == "" {
			inss, err := cdriver.ListInstances(region)
			if err != nil {
				log.Fatal(err)
			}
			for _, ins := range inss {
				callback(cdriver, region, ins.Name, ins.ID)

			}
		} else {
			instanceIDs := strings.Split(insids, ",")
			for _, instanceID := range instanceIDs {
				insinfo, err := cdriver.InstanceInfo(region, instanceID)
				if err != nil {
					fmt.Printf("查询%s地域下的%s信息失败", region, instanceID)
				}
				callback(cdriver, region, insinfo.Name, instanceID)
			}
		}

	}

}

func batchOperatorInstances(operator string, checkCallback func() error, callback func(cdriver driver.Driver, region, name string, insid string, args ...interface{}) error) {

	batchBaseOperatorInstances(checkCallback, func(cdriver driver.Driver, region, name string, insid string, args ...interface{}) {

		err := callback(cdriver, region, name, insid, args...)
		if err != nil {
			fmt.Printf("%s地域的实例%s(%s)%s失败，原因是:%s \n", region, name, insid, operator, err.Error())
		} else {
			fmt.Printf("%s地域的实例%(%s)s%s成功\n", region, name, insid, operator)
		}
	})

}

func stopInstances() {

	batchOperatorInstances("停止", func() error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.StopInstances(region, []string{insid})
	})

}

func startInstances() {

	batchOperatorInstances("启动", func() error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.StartInstances(region, []string{insid})
	})
}

func reStartInstances() {

	batchOperatorInstances("重启", func() error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.RestartInstances(region, []string{insid})
	})

}

func terminateInstances() {

	batchOperatorInstances("销毁", func() error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.TerminateInstances(region, []string{insid})
	})

}

func resetInstances() {

	var blueprintId string
	flag.StringVar(&blueprintId, "imageid", "", "镜像ID，可以用过lhbin image list 查询可以使用的镜像")

	batchOperatorInstances("重置镜像", func() error {
		checkArg(&blueprintId, "镜像ID不能为空，可以用过lhbin image list 查询可以使用的镜像")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.ResetInstances(region, []string{insid}, blueprintId)
	})

}

func resetInstancesPassword() {
	var username string
	var password string

	flag.StringVar(&username, "username", "", "用户名")
	flag.StringVar(&password, "password", "", "密码")

	batchOperatorInstances("重置密码", func() error {
		checkArg(&username, "用户名不能为空")
		checkArg(&password, "密码不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.ResetPassword(region, []string{insid}, username, password)
	})
}

func instance() {
	operator := os.Args[2]
	if operator == "list" {
		listInstances()
	}

	if operator == "desc" || operator == "describe" {
		describeInstance()
	}

	if operator == "stop" {
		RiskOperation("请确认已经保存好相关的工作", stopInstances)
	}

	if operator == "start" {
		startInstances()
	}

	if operator == "restart" {
		RiskOperation("请确认已经保存好相关的工作", reStartInstances)
	}

	if operator == "reset" {
		DangerOperation("重置服务器后无法恢复，请注意备份好相关数据", resetInstances)
	}

	if operator == "terminate" {
		DangerOperation("销毁服务器后无法恢复，请注意备份好相关数据。是否退款以腾讯云官方为准。", terminateInstances)
	}

	if operator == "passwd" {
		RiskOperation("重置密码会强制关机，请确认已经保存好相关的工作", resetInstancesPassword)
	}
}

func RiskOperation(tips string, callback func()) {
	fmt.Println("警告，下面的操作具有一定的风险性，请谨慎操作:")
	if tips != "" {
		fmt.Println(tips)
	}
	fmt.Print("请输入Y来确认是否进行下一步操作（不区分大小写，输入其他任意字符取消操作）:")
	var confirm string
	fmt.Scan(&confirm)
	fmt.Println("")
	if strings.ToLower(confirm) == "y" {
		callback()
	} else {
		fmt.Println("操作已经取消")
	}

}

func DangerOperation(tips string, callback func()) {
	fmt.Println("警告，下面的操作十分具备危险性，如非必要，强烈建议到控制台操作:")
	if tips != "" {
		fmt.Println(tips)
	}
	randStr := uuid.NewString()[:5]
	fmt.Printf("请输入%s来确认是否进行下一步操作（输入错误会取消操作）:", randStr)
	var confirmStr string
	fmt.Scan(&confirmStr)
	fmt.Println("")
	if confirmStr == randStr {
		callback()
	} else {
		fmt.Println("输入错误，操作已经取消")
	}

}

func queryTps() {
	temp := append(os.Args[0:1], "list")
	os.Args = append(temp, os.Args[1:]...)
	fmt.Println("--------------------------------------")
	batchBaseOperatorInstances(func() error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
		tps, err := cdriver.InstancesTrafficPackages(region, []string{insid})
		if err != nil {
			fmt.Printf("|%s|%s(%s)|查询失败，原因:%s|\n", region, name, insid, err.Error())
		} else {
			tp := tps[0]
			fmt.Printf("|%s|%s(%s)|%s|%s|%s|\n", region, name, insid, wellSize(tp.Total), wellSize(tp.Used), wellSize(tp.Remaining))
		}
		fmt.Println("---------------------------------------")
	})
}

func printHelp() {

}

func wellSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d Byte", size)
	}

	if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(1024))
	}

	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(1024*1024))
	}

	if size < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(1024*1024*1024))
	}

	return fmt.Sprintf("%.2f TB", float64(size)/float64(1024*1024*1024*1024))
}

func main() {

	if len(os.Args) <= 1 {
		printHelp()
	}

	childCommand := os.Args[1]

	childCommand = strings.ToLower(childCommand)

	if childCommand == "config" {
		// 配置ak
		configAccountAk()
	}

	if childCommand == "region" {
		queryRegions()
	}

	if childCommand == "ins" || childCommand == "instance" {
		instance()
	}

	if childCommand == "tp" || childCommand == "trafficpackage" || childCommand == "trafficpackages" {
		queryTps()
	}

	printHelp()

}
