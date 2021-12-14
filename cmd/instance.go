package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lixiaofei123/lhbin/driver"
)

const InstanceCommandName string = "ins"

func init() {
	RegisterChildCommand(InstanceCommandName, "管理轻量服务器实例信息", []string{"instance"})

	RegisterChildCommandOperator(InstanceCommandName, "list", "列出指定条件的轻量实例列表", []string{}, SafeOperation(ListInstances))
	RegisterChildCommandOperator(InstanceCommandName, "desc", "查看指定条件的轻量实例详情", []string{"describe"}, SafeOperation(DescribeInstances))
	RegisterChildCommandOperator(InstanceCommandName, "stop", "停止指定条件的轻量实例", []string{}, RiskOperation("请确认已经保存好相关的工作", StopInstances))
	RegisterChildCommandOperator(InstanceCommandName, "start", "启动指定条件的轻量实例", []string{}, SafeOperation(StartInstances))
	RegisterChildCommandOperator(InstanceCommandName, "restart", "重启指定条件的轻量实例", []string{"reboot"}, RiskOperation("请确认已经保存好相关的工作", RebootInstances))
	RegisterChildCommandOperator(InstanceCommandName, "reset", "重置指定条件的轻量服务器的镜像", []string{}, DangerOperation("重置服务器后无法恢复，请注意备份好相关数据", ResetInstances))
	RegisterChildCommandOperator(InstanceCommandName, "terminate", "销毁指定条件的轻量实例", []string{"destory"}, DangerOperation("销毁服务器后无法恢复，请注意备份好相关数据。是否退款以腾讯云官方为准。", TerminateInstances))
}

func baseBatchOperatorInstances(secondConfirm bool, checkCallback func(region string, insids string) error, callback func(cdriver driver.Driver, region, name string, insid string, args ...interface{})) error {

	var region string
	var insids string
	var insid string
	var force bool
	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "实例所在地域，不填则默认为所有可用区")
		flag.StringVar(&insid, "insid", "", "实例ID，如果设置此值，则会忽略insids参数")
		flag.StringVar(&insids, "insids", "", "实例ID，多个请用逗号隔开。如果不填则默认为所选择可用区下的所有实例")
		flag.BoolVar(&force, "f", false, "强制执行，忽略二次确认")
	}, func() error {
		if insid != "" {
			insids = insid
		}
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
	}

	err = checkCallback(region, insids)
	if err != nil {
		return err
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
				return nil
			}
		}

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

	return nil

}

func batchOperatorInstances(operator string, secondConfirm bool, checkCallback func(region string, insids string) error, callback func(cdriver driver.Driver, region, name string, insid string, args ...interface{}) error) error {

	return baseBatchOperatorInstances(secondConfirm, checkCallback, func(cdriver driver.Driver, region, name string, insid string, args ...interface{}) {

		err := callback(cdriver, region, name, insid, args...)
		if err != nil {
			fmt.Printf("%s地域的实例%s(%s)%s失败，原因是:%s \n", region, name, insid, operator, err.Error())
		} else {
			fmt.Printf("%s地域的实例(%s)%s%s成功\n", region, name, insid, operator)
		}
	})

}

func ListInstances() error {

	var region string

	cdriver, err := parseAndGetDriverWithS(func() {
		flag.StringVar(&region, "region", "", "地域，不填写则默认为所有地域")

	}, os.Args[3:])

	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------")
	fmt.Println("| 地域 | 实例名称 | 实例ID | 公网IP | 内网IP | 状态 |")
	fmt.Println("------------------------------------------")

	if region == "" {

		regions, err := cdriver.ListRegions()
		if err != nil {
			return err
		}

		for _, region := range regions {
			inss, err := cdriver.ListInstances(region.Region)
			if err != nil {
				return err
			}
			for _, ins := range inss {
				fmt.Println("|", region.Name, "|", ins.Name, "|", ins.ID, "|", ins.PublicIP, "|", ins.PrivateIP, "|", ins.State, "|")
				fmt.Println("------------------------------------------")
			}
		}

	} else {
		inss, err := cdriver.ListInstances(region)
		if err != nil {
			return err
		}
		for _, ins := range inss {
			fmt.Println("|", ins.Region, "|", ins.Name, "|", ins.ID, "|", ins.PublicIP, "|", ins.PrivateIP, "|", ins.State, "|")
			fmt.Println("------------------------------------------")
		}

	}

	fmt.Println("详细信息可以通过 lhbin ins desc --region region --insids lhins-xxxxx,lhins-yyyyy 命令进行查看")
	return nil
}

func DescribeInstances() error {

	return baseBatchOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
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

func StopInstances() error {

	return batchOperatorInstances("停止", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.StopInstances(region, []string{insid})
	})

}

func StartInstances() error {

	return batchOperatorInstances("启动", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.StartInstances(region, []string{insid})
	})
}

func RebootInstances() error {

	return batchOperatorInstances("重启", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.RestartInstances(region, []string{insid})
	})

}

func TerminateInstances() error {

	return batchOperatorInstances("销毁", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {

		return cdriver.TerminateInstances(region, []string{insid})
	})

}

func ResetInstances() error {

	var blueprintId string
	flag.StringVar(&blueprintId, "imageid", "", "镜像ID，可以用过lhbin image list 查询可以使用的镜像")

	return batchOperatorInstances("重置镜像", true, func(region string, insids string) error {
		checkArg(&blueprintId, "镜像ID不能为空，可以用过lhbin image list 查询可以使用的镜像")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.ResetInstances(region, []string{insid}, blueprintId)
	})

}

func ResetInstancesPassword() error {
	var username string
	var password string

	flag.StringVar(&username, "username", "", "用户名")
	flag.StringVar(&password, "password", "", "密码")

	return batchOperatorInstances("重置密码", true, func(region string, insids string) error {
		checkArg(&username, "用户名不能为空")
		checkArg(&password, "密码不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.ResetPassword(region, []string{insid}, username, password)
	})
}
