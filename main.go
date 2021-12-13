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
				fmt.Println("|", region.Name, "|", ins.Name, "|", ins.ID, "|", ins.PublicIP, "|", ins.PrivateIP, "|", ins.State, "|")
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

	fmt.Println("详细信息可以通过 lhbin ins desc --region region --insids lhins-xxxxx,lhins-yyyyy 命令进行查看")

}

func describeInstances() {

	batchBaseOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
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

func batchBaseOperatorInstances(secondConfirm bool, checkCallback func(region string, insids string) error, callback func(cdriver driver.Driver, region, name string, insid string, args ...interface{})) {
	var driverName string
	var account string
	var region string
	var insids string
	var insid string
	var force bool

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "实例所在地域，不填则默认为所有可用区")
	flag.StringVar(&insid, "insid", "", "实例ID，如果设置此值，则会忽略insids参数")
	flag.StringVar(&insids, "insids", "", "实例ID，多个请用逗号隔开。如果不填则默认为所选择可用区下的所有实例")
	flag.BoolVar(&force, "f", false, "强制执行，忽略二次确认")

	flag.CommandLine.Parse(os.Args[3:])

	if insid != "" {
		insids = insid
	}

	err := checkCallback(region, insids)
	if err != nil {
		log.Panic(err)
	}

	if secondConfirm && !force {

		var needConfirm = false

		if region == "" && insids == "" {
			fmt.Println("未设置地域和实例ID，操作将对所有的实例生效")
			needConfirm = true
		}

		if region != "" && insids == "" {
			fmt.Println("设置了地域，但未设置实例ID，操作将对此地域下的所有的实例生效")
			needConfirm = true
		}

		if needConfirm {
			fmt.Println("如果不希望出现此确认步骤，请加上-f参数来强制运行")
			fmt.Print("请输入Y来确认是否进行下一步操作（不区分大小写，输入其他任意字符取消操作）:")
			var confirm string
			fmt.Scan(&confirm)
			fmt.Println("")
			if strings.ToLower(confirm) != "y" {
				log.Println("操作已取消")
				return
			}
		}

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
						fmt.Printf("查询%s地域下的%s信息失败 \n", region.Region, instanceID)
					} else {
						callback(cdriver, region.Region, insinfo.Name, instanceID)
					}

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
					fmt.Printf("查询%s地域下的%s信息失败 \n", region, instanceID)
				} else {
					callback(cdriver, region, insinfo.Name, instanceID)
				}

			}
		}

	}

}

func batchOperatorInstances(operator string, secondConfirm bool, checkCallback func(region string, insids string) error, callback func(cdriver driver.Driver, region, name string, insid string, args ...interface{}) error) {

	batchBaseOperatorInstances(secondConfirm, checkCallback, func(cdriver driver.Driver, region, name string, insid string, args ...interface{}) {

		err := callback(cdriver, region, name, insid, args...)
		if err != nil {
			fmt.Printf("%s地域的实例%s(%s)%s失败，原因是:%s \n", region, name, insid, operator, err.Error())
		} else {
			fmt.Printf("%s地域的实例(%s)%s%s成功\n", region, name, insid, operator)
		}
	})

}

func stopInstances() {

	batchOperatorInstances("停止", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.StopInstances(region, []string{insid})
	})

}

func startInstances() {

	batchOperatorInstances("启动", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.StartInstances(region, []string{insid})
	})
}

func reStartInstances() {

	batchOperatorInstances("重启", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.RestartInstances(region, []string{insid})
	})

}

func terminateInstances() {

	batchOperatorInstances("销毁", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.TerminateInstances(region, []string{insid})
	})

}

func resetInstances() {

	var blueprintId string
	flag.StringVar(&blueprintId, "imageid", "", "镜像ID，可以用过lhbin image list 查询可以使用的镜像")

	batchOperatorInstances("重置镜像", true, func(region string, insids string) error {
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

	batchOperatorInstances("重置密码", true, func(region string, insids string) error {
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
		describeInstances()
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
	batchBaseOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
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

func listSnapshots() {
	fmt.Println("------------------------------------------")
	fmt.Println("| 地域 | 实例名称 | 实例ID | 快照名称 | 快照ID | 创建时间 |状态 |")
	fmt.Println("------------------------------------------")

	batchBaseOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
		snapshots, err := cdriver.ListSnapshots(region, insid)
		if err != nil {
			fmt.Printf("|%s|%s(%s)|查询失败，原因:%s|\n", region, name, insid, err.Error())
			fmt.Println("------------------------------------------")
		} else {
			for _, snapshot := range snapshots {
				fmt.Println("|", region, "|", name, "|", insid, "|", snapshot.Name, "|", snapshot.SnapShot, "|", snapshot.CreatedTime.Format("2006-01-02 15:04:05"), "|", snapshot.State, "|")
				fmt.Println("------------------------------------------")
			}
		}
	})

	fmt.Println("详细信息可以通过 lhbin ss desc --region region --ssid lhsnap-xxxxxxxxx 命令进行查看")

}

func describeSnapshots() {

	var driverName string
	var account string
	var region string
	var snapshotID string
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域，必须填写")
	flag.StringVar(&snapshotID, "ssid", "", "快照ID")
	flag.CommandLine.Parse(os.Args[3:])

	checkArg(&region, "地域不能为空")
	checkArg(&snapshotID, "快照ID不能为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	snapshot, err := cdriver.SnapshotInfo(region, snapshotID)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("-------------------------------")
	fmt.Println("| 地域 | ", region, "|")
	fmt.Println("| 快照ID | ", snapshot.SnapShot, "|")
	fmt.Println("| 快照名称 | ", snapshot.Name, "|")
	fmt.Println("| 状态 | ", snapshot.State, "|")
	fmt.Println("| 进度 | ", snapshot.Percent, "|")
	fmt.Println("| 创建时间 | ", snapshot.CreatedTime.Format("2006-01-02 15:04:05"), "|")
	fmt.Println("-------------------------------")

}

func deleteSnapshots() {

	var driverName string
	var account string
	var region string
	var snapshotID string
	var snapshotIDs string
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域，必须填写")
	flag.StringVar(&snapshotID, "ssid", "", "快照ID,如果填写此项，则忽略ssids参数的值")
	flag.StringVar(&snapshotIDs, "ssids", "", "快照ID列表，用逗号隔开")
	flag.CommandLine.Parse(os.Args[3:])

	if snapshotID != "" {
		snapshotIDs = snapshotID
	}
	checkArg(&region, "地域不能为空")
	checkArg(&snapshotIDs, "快照ID不能为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	ssids := strings.Split(snapshotIDs, ",")

	for _, ssid := range ssids {
		err := cdriver.DeleteSnapshots(region, []string{ssid})
		if err != nil {
			fmt.Printf("%s地域的快照%s删除失败，原因是:%s \n", region, ssid, err.Error())
		} else {
			fmt.Printf("%s地域的快照%s删除成功 \n", region, ssid)
		}
	}
}

func createSnapshot() {

	var ssname string
	flag.StringVar(&ssname, "name", "", "快照名称")

	batchOperatorInstances("创建快站", true, func(region string, insids string) error {
		checkArg(&ssname, "快照名称不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		_, err := cdriver.CreateSnapshot(region, insid, ssname)
		return err
	})

}

func applySnapshot() {

	var driverName string
	var account string
	var region string
	var insid string
	var snapshotID string
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域，必须填写")
	flag.StringVar(&snapshotID, "ssid", "", "快照ID")
	flag.StringVar(&insid, "insid", "", "实例ID")
	flag.CommandLine.Parse(os.Args[3:])

	checkArg(&region, "地域不能为空")
	checkArg(&snapshotID, "快照ID不能为空")
	checkArg(&insid, "实例ID不能为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	err = cdriver.ApplySnapshot(region, insid, snapshotID)
	if err != nil {
		fmt.Printf("%s地域的实例%s恢复快照%s失败，原因:%s\n", region, insid, snapshotID, err.Error())
	} else {
		fmt.Printf("%s地域的实例%s恢复快照%s成功\n", region, insid, snapshotID)
	}

}

func snapshot() {
	operator := os.Args[2]
	if operator == "list" {
		listSnapshots()
	}

	if operator == "desc" || operator == "describe" {
		describeSnapshots()
	}

	if operator == "del" || operator == "delete" {
		RiskOperation("", deleteSnapshots)
	}

	if operator == "create" {
		createSnapshot()
	}

	if operator == "apply" {
		DangerOperation("恢复快照会丢失创建快照以后的数据", applySnapshot)
	}

}

func listBlueprints() {

	var driverName string
	var account string
	var region string
	var platform string
	var imageType string
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域")
	flag.StringVar(&platform, "platform", "all", "操作系统，可选值为all、linux、win")
	flag.StringVar(&imageType, "type", "ALL", "镜像类型，可选值为all、app(应用镜像)、system(系统镜像)、private(私有镜像)、shared(共享镜像)")
	flag.CommandLine.Parse(os.Args[3:])

	checkArg(&region, "地域不能为空")

	platformType := driver.AllPlatform
	if platform == "linux" {
		platformType = driver.LinuxPlatform
	}
	if platform == "win" || platform == "window" || platform == "windows" {
		platformType = driver.WinPlatform
	}

	blueprintType := driver.AllBlueprint
	if imageType == "app" {
		blueprintType = driver.AppBlueprint
	}

	if imageType == "system" || imageType == "pure" {
		blueprintType = driver.PureBlueprint
	}

	if imageType == "private" {
		blueprintType = driver.PrivateBlueprint
	}

	if imageType == "shared" {
		blueprintType = driver.SharedBlueprint
	}

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	bps, err := cdriver.ListBlueprints(region, platformType, blueprintType)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("------------------------------------------")
	fmt.Println("| 镜像名称 | 镜像ID | 系统 | 磁盘 | 内存 | 状态 |")
	fmt.Println("------------------------------------------")

	for _, bp := range bps {
		fmt.Println("|", bp.Name, "|", bp.Blueprint, "|", bp.OsName, "|", bp.RequiredDiskSize, "GB |", bp.RequiredMemory, "GB |", bp.State, "|")
		fmt.Println("----------------------------------------------------")
	}

	fmt.Println("详细信息可以通过 lhbin image desc --region region --imageid lhsnap-xxxxxxxxx 命令进行查看")

}

func describeBlueprint() {

	var driverName string
	var account string
	var region string
	var blurprintID string
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域")
	flag.StringVar(&blurprintID, "imageid", "", "镜像ID")
	flag.CommandLine.Parse(os.Args[3:])

	checkArg(&region, "地域不能为空")
	checkArg(&blurprintID, "镜像ID不能为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	blueprint, err := cdriver.BlueprintInfo(region, blurprintID)
	if err != nil {
		fmt.Printf("查询%s地域下的镜像%s信息失败，原因是:%s \n", region, blurprintID, err.Error())
	}

	fmt.Println("-------------------------------")
	fmt.Println("| 镜像名称 | ", blueprint.Name, "|")
	fmt.Println("| 镜像ID | ", blueprint.Blueprint, "|")
	fmt.Println("| 操作系统 | ", blueprint.OsName, "|")
	fmt.Println("| 最小磁盘要求 | ", blueprint.RequiredDiskSize, "GB |")
	fmt.Println("| 最小内存要求 | ", blueprint.RequiredMemory, "GB |")
	fmt.Println("| 描述 | ", blueprint.Description, "|")
	fmt.Println("-------------------------------")

}

func deleteBlueprints() {

	var driverName string
	var account string
	var region string
	var blurprintID string
	var blurprintIDs string

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域")
	flag.StringVar(&blurprintID, "imageid", "", "镜像ID，如果填写此项，则忽略imageids参数")
	flag.StringVar(&blurprintIDs, "imageids", "", "镜像ID列表，用逗号隔开")

	flag.CommandLine.Parse(os.Args[3:])

	if blurprintID != "" {
		blurprintIDs = blurprintID
	}

	checkArg(&region, "地域不能为空")
	checkArg(&blurprintIDs, "镜像ID或者镜像ID列表不能都为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	bpids := strings.Split(blurprintIDs, ",")

	for _, bpid := range bpids {
		err := cdriver.DeleteBlueprints(region, []string{bpid})
		if err != nil {
			fmt.Printf("%s地域下的镜像%s删除失败，原因是:%s \n", region, bpid, err.Error())
		} else {
			fmt.Printf("%s地域下的镜像%s删除成功， \n", region, bpid)
		}
	}

}

func createBlueprint() {

	var bpname string
	var desc string
	flag.StringVar(&bpname, "name", "", "镜像名称")
	flag.StringVar(&desc, "desc", "", "镜像描述")

	batchOperatorInstances("创建镜像", true, func(region string, insids string) error {
		checkArg(&bpname, "镜像名称不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		_, err := cdriver.CreateBlueprint(region, insid, bpname, desc)
		return err
	})

}

func blueprint() {
	operator := os.Args[2]
	if operator == "list" {
		listBlueprints()
	}

	if operator == "desc" || operator == "describe" {
		describeBlueprint()
	}

	if operator == "del" || operator == "delete" {
		DangerOperation("镜像删除后不能恢复", deleteBlueprints)
	}

	if operator == "create" {
		createBlueprint()
	}

}

func listFirewalls() {
	batchBaseOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
		rules, err := cdriver.ListFirewallRules(region, insid)
		if err != nil {
			fmt.Printf("%s地域下的%s(%s)防火墙规则查询失败，原因是:%s \n", region, name, insid, err.Error())
		} else {
			fmt.Println("---------------------------------------")
			fmt.Println("| 地域 |", region, "| 实例名称 |", name, "| 实例ID |", insid, "|")
			fmt.Println("---------------------------------------")
			fmt.Println("| 协议 | 端口 | 来源 | 策略 | 描述 |")
			fmt.Println("---------------------------------------")
			for _, rule := range rules {
				fmt.Println("|", rule.Protocol, "|", rule.Port, "|", rule.CidrBlock, "|", rule.Action, "|", rule.Description, "|")
				fmt.Println("---------------------------------------")
			}
			fmt.Println()

		}
	})
}

func deleteFirewallRules() {
	fmt.Println("防火墙的规则写法如下:")
	fmt.Println("protocol|port|cidr|action")
	fmt.Println("其中，protocal的取值是TCP、UDP、ICMP、ALL")
	fmt.Println("port可以是端口1,端口2这样的形式，也可以是起始端口-结束端口，端口的范围是1-65535，也可以输入ALL代表1-65535")
	fmt.Println("cidr为ip或者ip/port，例如0.0.0.0/0或者0.0.0.0，可以不填，默认为0.0.0.0/0")
	fmt.Println("action的取值为ACCEPT或者DROP，可以不填，默认为ACCEPT")
	fmt.Println("示例1 TCP|8080-8090|0.0.0.0/0|ACCEPT")
	fmt.Println("示例2 TCP|8080-8090")
	fmt.Println()
	fmt.Println("请输入你要删除的防火墙规则:")
	var rule string
	fmt.Scan(&rule)
	attrs := strings.Split(rule, "|")
	if len(attrs) < 2 {
		fmt.Println("至少需要输入协议和端口")
		return
	}

	deleteRule := &driver.FirewallRule{
		Protocol:  driver.FirewallRuleProtocol(attrs[0]),
		Port:      attrs[1],
		CidrBlock: "0.0.0.0/0",
		Action:    "ACCEPT",
	}

	if len(attrs) >= 3 {
		deleteRule.CidrBlock = attrs[2]
	}

	if len(attrs) >= 4 {
		deleteRule.Action = driver.FirewallRuleAction(attrs[3])
	}

	batchOperatorInstances("删除防火墙规则", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.DeleteFirewallRules(region, insid, []*driver.FirewallRule{deleteRule})
	})
}

func addFirewallRules() {
	fmt.Println("防火墙的规则写法如下:")
	fmt.Println("protocol|port|cidr|action|desc")
	fmt.Println("其中，protocal的取值是TCP、UDP、ICMP、ALL")
	fmt.Println("port可以是端口1,端口2这样的形式，也可以是起始端口-结束端口，端口的范围是1-65535，也可以输入ALL代表1-65535")
	fmt.Println("cidr为ip或者ip/port，例如0.0.0.0/0或者0.0.0.0，可以不填，默认为0.0.0.0/0")
	fmt.Println("action的取值为ACCEPT或者DROP，可以不填，默认为ACCEPT")
	fmt.Println("desc可以不填，默认为空")

	fmt.Println("示例1 TCP|8080-8090|0.0.0.0/0|ACCEPT|测试")
	fmt.Println("示例2 TCP|8080-8090")
	fmt.Println()
	fmt.Println("请输入你要添加的防火墙规则:")
	var rule string
	fmt.Scan(&rule)
	attrs := strings.Split(rule, "|")
	if len(attrs) < 2 {
		fmt.Println("至少需要输入协议和端口")
		return
	}

	addRule := &driver.FirewallRule{
		Protocol:    driver.FirewallRuleProtocol(attrs[0]),
		Port:        attrs[1],
		CidrBlock:   "0.0.0.0/0",
		Action:      "ACCEPT",
		Description: "",
	}

	if len(attrs) >= 3 {
		addRule.CidrBlock = attrs[2]
	}

	if len(attrs) >= 4 {
		addRule.Action = driver.FirewallRuleAction(attrs[3])
	}

	if len(attrs) >= 5 {
		addRule.Description = attrs[4]
	}

	batchOperatorInstances("添加防火墙规则", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.AddFirewallRules(region, insid, []*driver.FirewallRule{addRule})
	})
}

func updateFirewallRules() {
	fmt.Println("防火墙的规则写法如下:")
	fmt.Println("protocol|port|cidr|action|desc")
	fmt.Println("其中，protocal的取值是TCP、UDP、ICMP、ALL")
	fmt.Println("port可以是端口1,端口2这样的形式，也可以是起始端口-结束端口，端口的范围是1-65535，也可以输入ALL代表1-65535")
	fmt.Println("cidr为ip或者ip/port，例如0.0.0.0/0或者0.0.0.0，可以不填，默认为0.0.0.0/0")
	fmt.Println("action的取值为ACCEPT或者DROP，可以不填，默认为ACCEPT")
	fmt.Println("desc可以不填，默认为空")
	fmt.Println("示例1 TCP|8080-8090|0.0.0.0/0|ACCEPT|测试")
	fmt.Println("示例2 TCP|8080-8090")

	fmt.Println()
	fmt.Print("是否添加默认规则，默认规则包含Ping、80、443、22、3389 (Y/N):")
	newrules := []*driver.FirewallRule{}
	var defalutRule string
	fmt.Scan(&defalutRule)
	fmt.Println("")
	if strings.ToLower(defalutRule) == "y" {
		newrules = append(newrules, &driver.FirewallRule{
			Protocol:    driver.IcmpRuleProtocol,
			Port:        "ALL",
			CidrBlock:   "0.0.0.0/0",
			Action:      driver.AcceptRuleAction,
			Description: "放通Ping",
		})

		newrules = append(newrules, &driver.FirewallRule{
			Protocol:    driver.TcpRuleProtocol,
			Port:        "22",
			CidrBlock:   "0.0.0.0/0",
			Action:      driver.AcceptRuleAction,
			Description: "放通Linux SSH登录",
		})

		newrules = append(newrules, &driver.FirewallRule{
			Protocol:    driver.TcpRuleProtocol,
			Port:        "3389",
			CidrBlock:   "0.0.0.0/0",
			Action:      driver.AcceptRuleAction,
			Description: "放通Windows远程桌面登录",
		})

		newrules = append(newrules, &driver.FirewallRule{
			Protocol:    driver.TcpRuleProtocol,
			Port:        "80",
			CidrBlock:   "0.0.0.0/0",
			Action:      driver.AcceptRuleAction,
			Description: "放通Web服务HTTP(80)，如 Apache、Nginx",
		})

		newrules = append(newrules, &driver.FirewallRule{
			Protocol:    driver.TcpRuleProtocol,
			Port:        "443",
			CidrBlock:   "0.0.0.0/0",
			Action:      driver.AcceptRuleAction,
			Description: "放通Web服务HTTPS(443)，如 Apache、Nginx",
		})

	}

	for {
		fmt.Print("是否继续输入(Y/N):")
		var beContinue string
		fmt.Scan(&beContinue)
		fmt.Println("")
		if strings.ToLower(beContinue) != "y" {
			break
		}

		fmt.Println("请输入你要添加的防火墙规则:")
		var rule string
		fmt.Scan(&rule)
		attrs := strings.Split(rule, "|")
		if len(attrs) < 2 {
			fmt.Println("至少需要输入协议和端口")
			continue
		}

		addRule := &driver.FirewallRule{
			Protocol:    driver.FirewallRuleProtocol(attrs[0]),
			Port:        attrs[1],
			CidrBlock:   "0.0.0.0/0",
			Action:      "ACCEPT",
			Description: "",
		}

		if len(attrs) >= 3 {
			addRule.CidrBlock = attrs[2]
		}

		if len(attrs) >= 4 {
			addRule.Action = driver.FirewallRuleAction(attrs[3])
		}

		if len(attrs) >= 5 {
			addRule.Description = attrs[4]
		}

		newrules = append(newrules, addRule)
	}

	batchOperatorInstances("更新防火墙规则", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.UpdateFirewallRules(region, insid, newrules)
	})

}

func firewall() {
	operator := os.Args[2]
	if operator == "list" {
		listFirewalls()
	}

	if operator == "del" || operator == "delete" {
		deleteFirewallRules()
	}

	if operator == "add" || operator == "create" {
		addFirewallRules()
	}

	if operator == "update" {
		RiskOperation("此操作会删除原有的所有的防火墙规则", updateFirewallRules)
	}
}

func createKeyPair() {

	var driverName string
	var account string
	var region string
	var keyName string

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域")
	flag.StringVar(&keyName, "keyname", "", "密钥对名称")

	flag.CommandLine.Parse(os.Args[3:])

	checkArg(&region, "地域不能为空")
	checkArg(&keyName, "密钥对名称不能为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	keypair, err := cdriver.CreateKeyPair(region, keyName)
	if err != nil {
		fmt.Printf("密钥对创建失败，原因是:%s\n", err.Error())
	} else {
		fmt.Printf("密钥对%s(%s)创建成功，私钥后续无法查询，请注意保存私钥 \n", keypair.KeyName, keypair.KeyId)
		fmt.Println("公钥为")
		fmt.Println(keypair.PublicKey)
		fmt.Println("私钥为")
		fmt.Println(keypair.PrivateKey)
	}

}

func listKeyPairs() {
	var driverName string
	var account string
	var region string

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域，不填则默认为所有可用区")

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
	fmt.Println("| 地域 | 密钥名称 | 密钥ID | 绑定实例 | 创建时间 |")
	fmt.Println("------------------------------------------")

	if region == "" {
		// 查询所有的
		regions, err := cdriver.ListRegions()
		if err != nil {
			log.Panic(err)
		}

		for _, region := range regions {
			kps, err := cdriver.ListKeyPair(region.Region)
			if err != nil {
				log.Panic(err)
			}
			for _, kp := range kps {
				fmt.Println("|", region.Region, "|", kp.KeyName, "|", kp.KeyId, "|", strings.Join(kp.AssociatedInstanceIds, ","), "|", kp.CreatedTime.Format("2006-01-02 15:04:05"), "|")
				fmt.Println("------------------------------------------")
			}
		}

	} else {
		kps, err := cdriver.ListKeyPair(region)
		if err != nil {
			log.Panic(err)
		}
		for _, kp := range kps {
			fmt.Println("|", region, "|", kp.KeyName, "|", kp.KeyId, "|", strings.Join(kp.AssociatedInstanceIds, ","), "|", kp.CreatedTime.Format("2006-01-02 15:04:05"), "|")
			fmt.Println("------------------------------------------")
		}
	}

}

func deleteKeyPairs() {

	var driverName string
	var account string
	var region string
	var keyId string
	var keyIds string

	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.StringVar(&region, "region", "", "地域")
	flag.StringVar(&keyId, "keyid", "", "密钥对ID，如果设置此值会忽略keyIds参数，如果和keyids都为空则会删除该区域下全部密钥对")
	flag.StringVar(&keyIds, "keyids", "", "密钥对ID列表，用逗号隔开，如果和keyid都为空则会删除该区域下全部密钥对")

	flag.CommandLine.Parse(os.Args[3:])

	if keyId != "" {
		keyIds = keyId
	}

	checkArg(&region, "地域不能为空")

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		log.Panic(err)
	}

	cdriver, err := driver.GetDriver(acc)
	if err != nil {
		log.Panic(err)
	}

	if keyIds != "" {
		kpids := strings.Split(keyIds, ",")

		for _, kpid := range kpids {
			err := cdriver.DeleteKeyPair(region, []string{kpid})
			if err != nil {
				fmt.Printf("%s地域的密钥对%s删除失败，原因是:%s \n", region, kpid, err.Error())
			} else {
				fmt.Printf("%s地域的密钥对%s删除成功 \n", region, kpid)
			}

		}
	} else {
		kps, err := cdriver.ListKeyPair(region)
		if err != nil {
			log.Panic(err)
		}
		for _, kp := range kps {
			err := cdriver.DeleteKeyPair(region, []string{kp.KeyId})
			if err != nil {
				fmt.Printf("%s地域的密钥对%s删除失败，原因是:%s \n", region, kp.KeyId, err.Error())
			} else {
				fmt.Printf("%s地域的密钥对%s删除成功 \n", region, kp.KeyId)
			}
		}
	}

}

func bindKeyPairs() {

	var keyId string

	flag.StringVar(&keyId, "keyid", "", "密钥对ID")

	batchBaseOperatorInstances(true, func(region string, insids string) error {
		checkArg(&region, "地域不能为空")
		checkArg(&keyId, "密钥对ID不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
		err := cdriver.BindKeyPairs(region, []string{keyId}, []string{insid})
		if err != nil {
			fmt.Printf("%s地域的密钥对%s绑定到实例%s(%s)失败，原因是:%s \n", region, keyId, name, insid, err.Error())
		} else {
			fmt.Printf("%s地域的密钥对%s绑定到实例%s(%s)成功 \n", region, keyId, name, insid)
		}
	})

}

func unBindKeyPairs() {

	var keyId string

	flag.StringVar(&keyId, "keyid", "", "密钥对ID")

	batchBaseOperatorInstances(true, func(region string, insids string) error {
		checkArg(&region, "地域不能为空")
		checkArg(&keyId, "密钥对ID不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
		err := cdriver.UnBindKeyPairs(region, []string{keyId}, []string{insid})
		if err != nil {
			fmt.Printf("%s地域的密钥对%s从实例%s(%s)解绑失败，原因是:%s \n", region, keyId, name, insid, err.Error())
		} else {
			fmt.Printf("%s地域的密钥对%s从实例%s(%s)解绑成功 \n", region, keyId, name, insid)
		}
	})

}

func keypair() {
	operator := os.Args[2]
	if operator == "list" {
		listKeyPairs()
	}
	if operator == "create" {
		createKeyPair()
	}

	if operator == "del" || operator == "delete" {
		RiskOperation("", deleteKeyPairs)
	}

	if operator == "bind" {
		RiskOperation("绑定过程会重启服务器，请注意保存好应用数据", bindKeyPairs)
	}

	if operator == "unbind" {
		RiskOperation("解绑过程会重启服务器，请注意保存好应用数据", unBindKeyPairs)
	}
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

	if childCommand == "ss" || childCommand == "snapshot" {
		snapshot()
	}

	if childCommand == "image" || childCommand == "images" || childCommand == "bp" || childCommand == "bps" || childCommand == "blueprint" || childCommand == "blueprints" {
		blueprint()
	}

	if childCommand == "firewall" || childCommand == "iptable" {
		firewall()
	}

	if childCommand == "kp" || childCommand == "keypair" {
		keypair()
	}

	printHelp()

}
