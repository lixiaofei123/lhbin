package driver

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
}
