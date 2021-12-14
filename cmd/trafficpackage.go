package cmd

import (
	"fmt"

	"github.com/lixiaofei123/lhbin/driver"
)

const TPCommandName string = "trafficpackage"

func init() {

	RegisterChildCommand(TPCommandName, "查看流量包详情", []string{"tp"})
	RegisterChildCommandOperator(TPCommandName, "list", "列出符合条件的实例的流量包详情", []string{}, SafeOperation(ListTrafficPackages))
}

func ListTrafficPackages() error {
	fmt.Println("--------------------------------------")
	fmt.Println("| 地域 | 实例ID | 总流量 | 已用流量 | 剩余流量 |")
	fmt.Println("------------------------------------------")
	return baseBatchOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
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
