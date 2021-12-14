package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/lixiaofei123/lhbin/config"
	"github.com/lixiaofei123/lhbin/driver"
)

func checkArg(arg *string, errText string) {
	if *arg == "" {
		fmt.Println(errText)
		os.Exit(-1)
	}
}

func parseAndGetDriverWithoutSV(arguments []string) (driver.Driver, error) {
	return parseAndGetDriver(func() {}, func() error { return nil }, arguments)
}

func parseAndGetDriverWithS(setArgsFunc func(), arguments []string) (driver.Driver, error) {
	return parseAndGetDriver(setArgsFunc, func() error { return nil }, arguments)
}

// func parseAndGetDriverWithV(vaildFunc func() error, arguments []string) (driver.Driver, error) {
// 	return parseAndGetDriver(func() {}, vaildFunc, arguments)
// }

func parseAndGetDriver(setArgsFunc func(), vaildFunc func() error, arguments []string) (driver.Driver, error) {

	setArgsFunc()

	var driverName string
	var account string // 账号
	flag.StringVar(&driverName, "driver", "qqcloud", "云厂商类型，目前仅支持腾讯云，默认为腾讯云")
	flag.StringVar(&account, "account", "", "账号名称，区分多用户使用，可随意指定。不指定则为默认Driver的第一个账号")
	flag.CommandLine.Parse(arguments)

	if err := vaildFunc(); err != nil {
		return nil, err
	}

	acc, err := config.FindAcount(config.DriverName(driverName), account)
	if err != nil {
		return nil, err
	}

	return driver.GetDriver(acc)
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

type OperationFunc func(showHelp bool) error

func SafeOperation(callback func() error) OperationFunc {
	return func(showHelp bool) error {
		return callback()
	}
}

func RiskOperation(tips string, callback func() error) OperationFunc {

	return func(showHelp bool) error {
		if !showHelp {
			fmt.Println("警告，下面的操作具有一定的风险性，请谨慎操作:")
			if tips != "" {
				fmt.Println(tips)
			}
			fmt.Print("请输入Y来确认是否进行下一步操作（不区分大小写，输入其他任意字符取消操作）:")
			var confirm string
			fmt.Scan(&confirm)
			fmt.Println("")
			if strings.ToLower(confirm) == "y" {
				return callback()
			} else {
				fmt.Println("操作已经取消")
			}
		} else {
			return callback()
		}
		return nil
	}

}

func DangerOperation(tips string, callback func() error) OperationFunc {

	return func(showHelp bool) error {
		if !showHelp {
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
				return callback()
			} else {
				fmt.Println("输入错误，操作已经取消")
			}
		} else {
			return callback()
		}

		return nil
	}

}

var commands map[string]*childCommand = map[string]*childCommand{}
var commandAliasMap map[string]string = map[string]string{}

type childCommand struct {
	command          string
	operatorAliasMap map[string]string
	operators        map[string]*commandOperator
	tips             string
}

type commandOperator struct {
	operatorFunc OperationFunc
	tips         string
}

func RegisterChildCommand(command string, tips string, alias []string) {

	c := &childCommand{
		command:          command,
		tips:             tips,
		operatorAliasMap: map[string]string{},
		operators:        map[string]*commandOperator{},
	}

	if len(alias) > 0 {
		for _, alia := range alias {
			commandAliasMap[alia] = command
		}
	}

	commandAliasMap[command] = command
	commands[command] = c

}

func RegisterChildCommandOperator(commandName, oprator, tips string, alias []string, operation OperationFunc) {
	if commandName, ok := commandAliasMap[commandName]; ok {
		if command, ok := commands[commandName]; ok {
			command.operators[oprator] = &commandOperator{
				tips:         tips,
				operatorFunc: operation,
			}

			if len(alias) > 0 {
				for _, alia := range alias {
					command.operatorAliasMap[alia] = oprator
				}
			}

			command.operatorAliasMap[oprator] = oprator
		}
	}
}

func printHeader() {
	fmt.Println("----------------------------------------")
	fmt.Println("--------------轻量服务器LHBIN--------------")
	fmt.Println("----------------V 0.0.2 -----------------")
	fmt.Println()
}

func printHelp() {

	printHeader()

	fmt.Print("本工具目前支持以下命令:\n\n")
	for key, command := range commands {
		fmt.Println("      ", key, "     ", command.tips)
	}

	fmt.Println()
	fmt.Println("输入 lhbin [命令名称] --help 查看命令帮助信息")
}

func printChildCommandHelp(childCommand string) {

	printHeader()

	if commandName, ok := commandAliasMap[childCommand]; ok {
		if command, ok := commands[commandName]; ok {
			fmt.Println(command.tips)
			fmt.Println()
			for key, operator := range command.operators {
				fmt.Println("      ", key, "     ", operator.tips)
			}
			fmt.Println()
			fmt.Println("输入 lhbin [命令名称] [操作名称]  --help 查看操作帮助信息")
		}
	}
}

func printChildCommandOperatorHelp(childCommand string, operatorName string) {

	printHeader()

	if commandName, ok := commandAliasMap[childCommand]; ok {
		if command, ok := commands[commandName]; ok {
			if operatorName, ok := command.operatorAliasMap[operatorName]; ok {
				if operator, ok := command.operators[operatorName]; ok {
					fmt.Println(operator.tips)
					fmt.Println()
					operator.operatorFunc(true)
				}
			}

		}
	}
}

func ExecuteCommand() {
	if len(os.Args) == 1 || (os.Args[1] == "--help" || os.Args[1] == "-help") {
		printHelp()
		return
	}
	if len(os.Args) == 2 || (os.Args[2] == "--help" || os.Args[2] == "-help") {
		printChildCommandHelp(os.Args[1])
		return
	}

	if len(os.Args) > 3 && (os.Args[3] == "--help" || os.Args[3] == "-help") {
		printChildCommandOperatorHelp(os.Args[1], os.Args[2])
		return
	}

	if len(os.Args) >= 3 {
		childCommand := os.Args[1]
		operatorName := os.Args[2]
		if commandName, ok := commandAliasMap[childCommand]; ok {
			if command, ok := commands[commandName]; ok {
				if operatorName, ok := command.operatorAliasMap[operatorName]; ok {
					if operator, ok := command.operators[operatorName]; ok {
						err := operator.operatorFunc(false)
						if err != nil {
							fmt.Printf("操作失败，原因是:%s \n", err.Error())
						} else {
							fmt.Println("操作成功.....")
						}
						return
					}
				}

			}
		}
	}

	printHelp()

}
