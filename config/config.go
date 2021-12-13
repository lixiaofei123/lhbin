package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type DriverName string

const (
	QQCloud DriverName = "qqcloud"
)

var GlobalConfig *Config
var configFilePath string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	configDir := path.Join(homedir, ".lhbin")

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0700)
	}

	if err != nil {
		log.Panic(err)
	}

	GlobalConfig = new(Config)

	configFilePath = path.Join(configDir, "config.yaml")
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {

		data, _ := yaml.Marshal(GlobalConfig)
		err = ioutil.WriteFile(configFilePath, data, 0660)
	}

	if err != nil {
		log.Panic(err)
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Panic(err)
	}

	err = yaml.Unmarshal(data, GlobalConfig)
	if err != nil {
		log.Panic(err)
	}

}

type Config struct {
	Accounts []*AccountConfig `yaml:"accounts"`
}

type AccountConfig struct {
	Driver   DriverName `yaml:"driver"`
	Account  string     `yaml:"account"`
	AKID     string     `yaml:"akid"`
	AKSecret string     `yaml:"aksecret"`
}

func AddAccount(newAccount *AccountConfig) {
	find := false
	for index, account := range GlobalConfig.Accounts {
		if account.Account == newAccount.Account && account.Driver == newAccount.Driver {
			GlobalConfig.Accounts[index] = newAccount
			find = true
			break
		}
	}

	if !find {
		GlobalConfig.Accounts = append(GlobalConfig.Accounts, newAccount)
	}

	data, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(configFilePath, data, 0660)
	if err != nil {
		log.Panic(err)
	}

}

func DeleteAccount(driver DriverName, delAccount string) {
	for index, account := range GlobalConfig.Accounts {
		if account.Account == delAccount && account.Driver == driver {
			GlobalConfig.Accounts = append(GlobalConfig.Accounts[0:index], GlobalConfig.Accounts[index+1:]...)
			break
		}
	}
	data, _ := yaml.Marshal(GlobalConfig)
	err := ioutil.WriteFile(configFilePath, data, 0660)
	if err != nil {
		log.Panic(err)
	}

}
