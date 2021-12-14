package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lixiaofei123/lhbin/driver"
)

const SnapshotCommandName string = "snapshot"

func init() {

	RegisterChildCommand(SnapshotCommandName, "管理快照信息", []string{"ss"})
	RegisterChildCommandOperator(SnapshotCommandName, "list", "列出符合要求的快照列表", []string{}, SafeOperation(ListSnapshots))
	RegisterChildCommandOperator(SnapshotCommandName, "desc", "查看符合要求的快照详情", []string{"describe"}, SafeOperation(DescribeSnapshots))
	RegisterChildCommandOperator(SnapshotCommandName, "del", "删除符合要求的快照", []string{"delete"}, RiskOperation("", DeleteSnapshots))
	RegisterChildCommandOperator(SnapshotCommandName, "create", "创建快照信息", []string{}, SafeOperation(CreateSnapshot))
	RegisterChildCommandOperator(SnapshotCommandName, "apply", "恢复快照", []string{}, DangerOperation("恢复快照会丢失创建快照以后的数据", ApplySnapshot))

}

func ListSnapshots() error {
	fmt.Println("------------------------------------------")
	fmt.Println("| 地域 | 实例名称 | 实例ID | 快照名称 | 快照ID | 创建时间 |状态 |")
	fmt.Println("------------------------------------------")

	err := baseBatchOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
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
	if err != nil {
		return err
	}
	fmt.Println("详细信息可以通过 lhbin ss desc --region region --ssid lhsnap-xxxxxxxxx 命令进行查看")
	return nil

}

func DescribeSnapshots() error {

	var region string
	var snapshotID string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域，必须填写")
		flag.StringVar(&snapshotID, "ssid", "", "快照ID")
	}, func() error {
		checkArg(&region, "地域不能为空")
		checkArg(&snapshotID, "快照ID不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
	}

	snapshot, err := cdriver.SnapshotInfo(region, snapshotID)
	if err != nil {
		return err
	}

	fmt.Println("-------------------------------")
	fmt.Println("| 地域 | ", region, "|")
	fmt.Println("| 快照ID | ", snapshot.SnapShot, "|")
	fmt.Println("| 快照名称 | ", snapshot.Name, "|")
	fmt.Println("| 状态 | ", snapshot.State, "|")
	fmt.Println("| 进度 | ", snapshot.Percent, "|")
	fmt.Println("| 创建时间 | ", snapshot.CreatedTime.Format("2006-01-02 15:04:05"), "|")
	fmt.Println("-------------------------------")

	return nil

}

func DeleteSnapshots() error {

	var region string
	var snapshotID string
	var snapshotIDs string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域，必须填写")
		flag.StringVar(&snapshotID, "ssid", "", "快照ID,如果填写此项，则忽略ssids参数的值")
		flag.StringVar(&snapshotIDs, "ssids", "", "快照ID列表，用逗号隔开")
	}, func() error {
		if snapshotID != "" {
			snapshotIDs = snapshotID
		}
		checkArg(&region, "地域不能为空")
		checkArg(&snapshotIDs, "快照ID不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
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

	return nil
}

func CreateSnapshot() error {

	var ssname string
	flag.StringVar(&ssname, "name", "", "快照名称")

	return batchOperatorInstances("创建快站", true, func(region string, insids string) error {
		checkArg(&ssname, "快照名称不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		_, err := cdriver.CreateSnapshot(region, insid, ssname)
		return err
	})

}

func ApplySnapshot() error {

	var region string
	var insid string
	var snapshotID string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域，必须填写")
		flag.StringVar(&snapshotID, "ssid", "", "快照ID")
		flag.StringVar(&insid, "insid", "", "实例ID")
	}, func() error {
		checkArg(&region, "地域不能为空")
		checkArg(&snapshotID, "快照ID不能为空")
		checkArg(&insid, "实例ID不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
	}

	err = cdriver.ApplySnapshot(region, insid, snapshotID)
	if err != nil {
		fmt.Printf("%s地域的实例%s恢复快照%s失败，原因:%s\n", region, insid, snapshotID, err.Error())
	} else {
		fmt.Printf("%s地域的实例%s恢复快照%s成功\n", region, insid, snapshotID)
	}

	return nil

}
