package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/lixiaofei123/lhbin/driver"
)

const KPCommandName string = "keypair"

func init() {

	RegisterChildCommand(KPCommandName, "查看管理密钥对信息", []string{"kp"})
	RegisterChildCommandOperator(KPCommandName, "list", "列出符合条件的密钥对", []string{}, SafeOperation(ListKeyPairs))
	RegisterChildCommandOperator(KPCommandName, "import", "从已经存在的公钥创建密钥对", []string{}, SafeOperation(ImportKeyPair))
	RegisterChildCommandOperator(KPCommandName, "create", "创建新的密钥对", []string{}, SafeOperation(CreateKeyPair))
	RegisterChildCommandOperator(KPCommandName, "del", "删除符合条件的密钥对", []string{"delete"}, RiskOperation("", DeleteKeyPairs))
	RegisterChildCommandOperator(KPCommandName, "bind", "将密钥对绑定到指定的实例上", []string{}, RiskOperation("绑定过程会重启服务器，请注意保存好应用数据", BindKeyPairs))
	RegisterChildCommandOperator(KPCommandName, "unbind", "将密钥对从指定的实例上解绑", []string{}, RiskOperation("解绑过程会重启服务器，请注意保存好应用数据", UnBindKeyPairs))
}

func CreateKeyPair() error {

	var region string
	var keyName string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域")
		flag.StringVar(&keyName, "keyname", "", "密钥对名称")
	}, func() error {
		checkArg(&region, "地域不能为空")
		checkArg(&keyName, "密钥对名称不能为空")
		return nil
	}, os.Args[3:])
	if err != nil {
		return nil
	}

	keypair, err := cdriver.CreateKeyPair(region, keyName)
	if err != nil {
		return err
	}

	fmt.Printf("密钥对%s(%s)创建成功，私钥后续无法查询，请注意保存私钥 \n", keypair.KeyName, keypair.KeyId)
	fmt.Println("公钥为")
	fmt.Println(keypair.PublicKey)
	fmt.Println("私钥为")
	fmt.Println(keypair.PrivateKey)

	return nil

}

func ImportKeyPair() error {

	var region string
	var pubKeyPath string
	var keyName string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域")
		flag.StringVar(&keyName, "keyname", "", "密钥对名称")
		flag.StringVar(&pubKeyPath, "pubKeyPath", "", "公钥文件路径")
	}, func() error {
		checkArg(&region, "地域不能为空")
		checkArg(&keyName, "密钥对名称不能为空")
		checkArg(&pubKeyPath, "公钥文件路径不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return nil
	}

	pubKeyData, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	keypair, err := cdriver.ImportKeyPair(region, keyName, string(pubKeyData))
	if err != nil {
		return err
	}

	fmt.Printf("密钥对%s(%s)导入成功 \n", keypair.KeyName, keypair.KeyId)

	return nil

}

func ListKeyPairs() error {

	var region string

	cdriver, err := parseAndGetDriverWithS(func() {
		flag.StringVar(&region, "region", "", "地域，不填则默认为所有可用区")
	}, os.Args[3:])

	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------")
	fmt.Println("| 地域 | 密钥名称 | 密钥ID | 绑定实例 | 创建时间 |")
	fmt.Println("------------------------------------------")

	if region == "" {
		// 查询所有的
		regions, err := cdriver.ListRegions()
		if err != nil {
			return err
		}

		for _, region := range regions {
			kps, err := cdriver.ListKeyPair(region.Region)
			if err != nil {
				return err
			}
			for _, kp := range kps {
				fmt.Println("|", region.Region, "|", kp.KeyName, "|", kp.KeyId, "|", strings.Join(kp.AssociatedInstanceIds, ","), "|", kp.CreatedTime.Format("2006-01-02 15:04:05"), "|")
				fmt.Println("------------------------------------------")
			}
		}

	} else {
		kps, err := cdriver.ListKeyPair(region)
		if err != nil {
			return err
		}
		for _, kp := range kps {
			fmt.Println("|", region, "|", kp.KeyName, "|", kp.KeyId, "|", strings.Join(kp.AssociatedInstanceIds, ","), "|", kp.CreatedTime.Format("2006-01-02 15:04:05"), "|")
			fmt.Println("------------------------------------------")
		}
	}

	return nil

}

func DeleteKeyPairs() error {

	var region string
	var keyId string
	var keyIds string

	cdriver, err := parseAndGetDriver(func() {
		flag.StringVar(&region, "region", "", "地域")
		flag.StringVar(&keyId, "keyid", "", "密钥对ID，如果设置此值会忽略keyIds参数，如果和keyids都为空则会删除该区域下全部密钥对")
		flag.StringVar(&keyIds, "keyids", "", "密钥对ID列表，用逗号隔开，如果和keyid都为空则会删除该区域下全部密钥对")
	}, func() error {
		if keyId != "" {
			keyIds = keyId
		}

		checkArg(&region, "地域不能为空")
		return nil
	}, os.Args[3:])

	if err != nil {
		return err
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
			return err
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

	return nil

}

func BindKeyPairs() error {

	var keyId string

	flag.StringVar(&keyId, "keyid", "", "密钥对ID")

	return baseBatchOperatorInstances(true, func(region string, insids string) error {
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

func UnBindKeyPairs() error {

	var keyId string

	flag.StringVar(&keyId, "keyid", "", "密钥对ID")

	return baseBatchOperatorInstances(true, func(region string, insids string) error {
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
