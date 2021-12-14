package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lixiaofei123/lhbin/driver"
)

const BlueprintCommandName string = "blueprint"

func init() {

	RegisterChildCommand(BlueprintCommandName, "查看管理镜像信息", []string{"image", "bp"})
	RegisterChildCommandOperator(BlueprintCommandName, "list", "列出符合要求的镜像列表", []string{}, SafeOperation(ListBlueprints))
	RegisterChildCommandOperator(BlueprintCommandName, "desc", "查看符合要求的镜像详情", []string{"describe"}, SafeOperation(DescribeBlueprint))
	RegisterChildCommandOperator(BlueprintCommandName, "del", "删除符合要求的镜像", []string{"delete"}, DangerOperation("镜像删除后不能恢复", DeleteBlueprints))
	RegisterChildCommandOperator(BlueprintCommandName, "create", "创建镜像", []string{}, SafeOperation(CreateBlueprint))

}

func ListBlueprints() error {

	var region string
	var platform string
	var imageType string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域")
		flag.StringVar(&platform, "platform", "all", "操作系统，可选值为all、linux、win")
		flag.StringVar(&imageType, "type", "ALL", "镜像类型，可选值为all、app(应用镜像)、system(系统镜像)、private(私有镜像)、shared(共享镜像)")
	}, func() error {
		checkArg(&region, "地域不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
	}

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

	bps, err := cdriver.ListBlueprints(region, platformType, blueprintType)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------")
	fmt.Println("| 镜像名称 | 镜像ID | 系统 | 磁盘 | 内存 | 状态 |")
	fmt.Println("------------------------------------------")

	for _, bp := range bps {
		fmt.Println("|", bp.Name, "|", bp.Blueprint, "|", bp.OsName, "|", bp.RequiredDiskSize, "GB |", bp.RequiredMemory, "GB |", bp.State, "|")
		fmt.Println("----------------------------------------------------")
	}

	fmt.Println("详细信息可以通过 lhbin image desc --region region --imageid lhsnap-xxxxxxxxx 命令进行查看")
	return nil

}

func DescribeBlueprint() error {

	var region string
	var blurprintID string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域")
		flag.StringVar(&blurprintID, "imageid", "", "镜像ID")
	}, func() error {
		checkArg(&region, "地域不能为空")
		checkArg(&blurprintID, "镜像ID不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
	}

	blueprint, err := cdriver.BlueprintInfo(region, blurprintID)
	if err != nil {
		return err
	}

	fmt.Println("-------------------------------")
	fmt.Println("| 镜像名称 | ", blueprint.Name, "|")
	fmt.Println("| 镜像ID | ", blueprint.Blueprint, "|")
	fmt.Println("| 操作系统 | ", blueprint.OsName, "|")
	fmt.Println("| 最小磁盘要求 | ", blueprint.RequiredDiskSize, "GB |")
	fmt.Println("| 最小内存要求 | ", blueprint.RequiredMemory, "GB |")
	fmt.Println("| 描述 | ", blueprint.Description, "|")
	fmt.Println("-------------------------------")

	return nil

}

func DeleteBlueprints() error {

	var region string
	var blurprintID string
	var blurprintIDs string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域")
		flag.StringVar(&blurprintID, "imageid", "", "镜像ID，如果填写此项，则忽略imageids参数")
		flag.StringVar(&blurprintIDs, "imageids", "", "镜像ID列表，用逗号隔开")
	}, func() error {
		if blurprintID != "" {
			blurprintIDs = blurprintID
		}
		checkArg(&region, "地域不能为空")
		checkArg(&blurprintIDs, "镜像ID或者镜像ID列表不能都为空")

		return nil
	}, os.Args[3:])

	if err != nil {
		return err
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

	return nil

}

func CreateBlueprint() error {

	var bpname string
	var desc string
	flag.StringVar(&bpname, "name", "", "镜像名称")
	flag.StringVar(&desc, "desc", "", "镜像描述")

	return batchOperatorInstances("创建镜像", true, func(region string, insids string) error {
		checkArg(&bpname, "镜像名称不能为空")
		return nil
	}, func(cdriver driver.Driver, region, name, insid string, args ...interface{}) error {
		_, err := cdriver.CreateBlueprint(region, insid, bpname, desc)
		return err
	})

}
