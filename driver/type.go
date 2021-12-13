package driver

import "time"

type RegionState string

const (
	RegionAvaliable RegionState = "AVAILABLE"
)

type Region struct {
	Name            string
	Region          string
	State           RegionState
	IsChinaMainland bool
}

type Zone struct {
	Name string
	Zone string
}

type InstanceState string

const (
	Pendding     InstanceState = "PENDING"
	LaunchFailed InstanceState = "LAUNCH_FAILED"
	Running      InstanceState = "RUNNING"
	Stoped       InstanceState = "STOPPED"
	Starting     InstanceState = "STARTING"
	Stoping      InstanceState = "STOPPING"
	Rebooting    InstanceState = "REBOOTING"
	Shutdown     InstanceState = "SHUTDOWN"
	Terminating  InstanceState = "TERMINATING"
)

type InstanceInfo struct {
	ID           string
	Name         string
	Region       string
	Zone         string
	Cpu          int
	Memory       int
	OSName       string
	Platform     string
	PlatformType string
	Disk         int
	PublicIP     string
	PrivateIP    string
	Bandwidth    int
	State        InstanceState
	CreatedTime  time.Time
	ExpiredTime  time.Time
}

type TrafficPackage struct {
	InstanceId string
	Used       int64
	Total      int64
	Remaining  int64
}

type SnapShotState string

const (
	SnapShotNormal      SnapShotState = "NORMAL"
	SnapShotCreating    SnapShotState = "CREATING"
	SnapShotRollbacking SnapShotState = "ROLLBACKING"
)

type SnapShot struct {
	SnapShot    string
	Name        string
	State       SnapShotState
	Percent     int
	CreatedTime time.Time
}

type BlueprintState string

type PlatformType string

const (
	AllPlatform   PlatformType = "ALL"
	LinuxPlatform PlatformType = "LINUX_UNIX"
	WinPlatform   PlatformType = "WINDOWS"
)

type BlueprintType string

const (
	AllBlueprint     BlueprintType = "ALL"
	AppBlueprint     BlueprintType = "APP_OS"
	PureBlueprint    BlueprintType = "PURE_OS"
	PrivateBlueprint BlueprintType = "PRIVATE"
	SharedBlueprint  BlueprintType = "SHARED"
)

type Blueprint struct {
	Blueprint        string
	Name             string
	Description      string
	OsName           string
	RequiredDiskSize int64
	RequiredMemory   int64
	State            BlueprintState
}

type FirewallRuleProtocol string

const (
	TcpRuleProtocol  FirewallRuleProtocol = "TCP"
	UdpRuleProtocol  FirewallRuleProtocol = "UDP"
	IcmpRuleProtocol FirewallRuleProtocol = "ICMP"
	AllPRulerotocol  FirewallRuleProtocol = "ALL"
)

type FirewallRuleAction string

const (
	AcceptRuleAction FirewallRuleAction = "ACCEPT"
	DropAcRuletion   FirewallRuleAction = "DROP"
)

type FirewallRule struct {
	Protocol    FirewallRuleProtocol
	Port        string
	CidrBlock   string
	Action      FirewallRuleAction
	Description string
}

type KeyPair struct {
	KeyId                 string
	KeyName               string
	PublicKey             string
	AssociatedInstanceIds []string
	CreatedTime           time.Time
	PrivateKey            string
}
