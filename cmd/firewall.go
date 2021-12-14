package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lixiaofei123/lhbin/driver"
)

const FirewallCommandName string = "firewall"

func init() {

	RegisterChildCommand(FirewallCommandName, "管理防火墙信息", []string{"iptable"})
	RegisterChildCommandOperator(FirewallCommandName, "list", "列出符合条件的实例的防火墙信息", []string{}, SafeOperation(ListFirewalls))
	RegisterChildCommandOperator(FirewallCommandName, "del", "列出符合条件的实例的防火墙信息中指定的防火墙规则", []string{"delete"}, SafeOperation(DeleteFirewallRules))
	RegisterChildCommandOperator(FirewallCommandName, "add", "向符合条件的实例的防火墙信息中添加新的防火墙规则", []string{"create"}, SafeOperation(AddFirewallRules))
	RegisterChildCommandOperator(FirewallCommandName, "update", "重置符合条件的实例的防火墙信息中的防火墙规则(会删除原有的全部规则并应用新添加的规则)", []string{"reset"}, RiskOperation("此操作会删除原有的所有的防火墙规则", UpdateFirewallRules))
}

func ListFirewalls() error {
	return baseBatchOperatorInstances(false, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) {
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

func PrintFileWallRuleTips() {
	fmt.Println("防火墙的规则写法如下:")
	fmt.Println("protocol|port|cidr|action|desc")
	fmt.Println("其中，protocal的取值是TCP、UDP、ICMP、ALL")
	fmt.Println("port可以是端口1,端口2这样的形式，也可以是起始端口-结束端口，端口的范围是1-65535，也可以输入ALL代表1-65535")
	fmt.Println("cidr为ip或者ip/port，例如0.0.0.0/0或者0.0.0.0，可以不填，默认为0.0.0.0/0")
	fmt.Println("action的取值为ACCEPT或者DROP，可以不填，默认为ACCEPT")
	fmt.Println("desc可以不填，默认为空。如果是删除防火墙规则，此项留空")

	fmt.Println("示例1 TCP|8080-8090|0.0.0.0/0|ACCEPT|测试")
	fmt.Println("示例2 TCP|8080-8090")
	fmt.Println()
}

func firewallRuleFromStr(ruleStr string) (*driver.FirewallRule, error) {

	attrs := strings.Split(ruleStr, "|")
	if len(attrs) < 2 {
		return nil, errors.New("至少需要输入协议和端口")
	}

	rule := &driver.FirewallRule{
		Protocol:  driver.FirewallRuleProtocol(attrs[0]),
		Port:      attrs[1],
		CidrBlock: "0.0.0.0/0",
		Action:    "ACCEPT",
	}

	if len(attrs) >= 3 {
		rule.CidrBlock = attrs[2]
	}

	if len(attrs) >= 4 {
		rule.Action = driver.FirewallRuleAction(attrs[3])
	}

	if len(attrs) >= 5 {
		rule.Description = attrs[4]
	}

	return rule, nil

}

func DeleteFirewallRules() error {

	PrintFileWallRuleTips()

	fmt.Println("请输入你要删除的防火墙规则:")
	var rule string
	fmt.Scan(&rule)

	deleteRule, err := firewallRuleFromStr(rule)
	if err != nil {
		return err
	}

	return batchOperatorInstances("删除防火墙规则", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.DeleteFirewallRules(region, insid, []*driver.FirewallRule{deleteRule})
	})
}

func AddFirewallRules() error {

	PrintFileWallRuleTips()

	fmt.Println("请输入你要添加的防火墙规则:")
	var rule string
	fmt.Scan(&rule)

	addRule, err := firewallRuleFromStr(rule)
	if err != nil {
		return err
	}

	return batchOperatorInstances("添加防火墙规则", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.AddFirewallRules(region, insid, []*driver.FirewallRule{addRule})
	})
}

func UpdateFirewallRules() error {

	PrintFileWallRuleTips()

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

		addRule, err := firewallRuleFromStr(rule)
		if err != nil {
			fmt.Println("防火墙规则输入错误:")
		} else {
			newrules = append(newrules, addRule)
		}

	}

	return batchOperatorInstances("更新防火墙规则", true, func(region string, insids string) error { return nil }, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		return cdriver.UpdateFirewallRules(region, insid, newrules)
	})

}
