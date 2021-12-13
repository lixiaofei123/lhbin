package driver

import (
	"errors"

	"github.com/lixiaofei123/lhbin/config"
)

type Driver interface {
	ListRegions() ([]*Region, error)
	ListZones(region string) ([]*Zone, error)

	ListInstances(region string) ([]*InstanceInfo, error)
	InstanceInfo(region, instanceID string) (*InstanceInfo, error)
	StopInstances(region string, instanceIDs []string) error
	StartInstances(region string, instanceIDs []string) error
	RestartInstances(region string, instanceIDs []string) error
	TerminateInstances(region string, instanceIDs []string) error
	ResetInstances(region string, instanceIDs []string, BlueprintId string) error
	ResetPassword(region string, instanceIDs []string, username, password string) error

	InstancesTrafficPackages(region string, instanceIDs []string) ([]*TrafficPackage, error)

	ListSnapshots(region, instanceID string) ([]*SnapShot, error)
	SnapshotInfo(region, snapshotID string) (*SnapShot, error)
	DeleteSnapshots(region string, snapshotIDs []string) error
	CreateSnapshot(region, instanceID, name string) (*SnapShot, error)
	ApplySnapshot(region, instanceID, snapshotID string) error

	ListBlueprints(region string, platformType PlatformType, blueprintType BlueprintType) ([]*Blueprint, error)
	BlueprintInfo(region, blueprintID string) (*Blueprint, error)
	DeleteBlueprints(region string, blueprintIDs []string) error
	CreateBlueprint(region, instanceId, name, desctiprtion string) (*Blueprint, error)

	ListFirewallRules(region string, instanceID string) ([]*FirewallRule, error)
	AddFirewallRules(region string, instanceID string, roles []*FirewallRule) error
	UpdateFirewallRules(region string, instanceID string, roles []*FirewallRule) error
	DeleteFirewallRules(region string, instanceID string, roles []*FirewallRule) error

	ListKeyPair(region string) ([]*KeyPair, error)
	CreateKeyPair(region string, name string) (*KeyPair, error)
	ImportKeyPair(region string, name string, publicKey string) (*KeyPair, error)
	DeleteKeyPair(region string, keyids []string) error
	BindKeyPairs(region string, keyids []string, instanceIDs []string) error
	UnBindKeyPairs(region string, keyids []string, instanceIDs []string) error
}

func GetDriver(account *config.AccountConfig) (Driver, error) {
	if account.Driver == config.QQCloud {
		return NewQQCloudLHDriver(account.AKID, account.AKSecret), nil
	}
	return nil, errors.New("没有合适的驱动")
}

func stringer(str []*string) []string {
	var strs []string
	for _, s := range str {
		if s == nil {
			strs = append(strs, "")
			continue
		}
		strs = append(strs, *s)
	}

	return strs
}
