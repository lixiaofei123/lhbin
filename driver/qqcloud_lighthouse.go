package driver

import (
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
)

type QQCloudLHDriver struct {
	credential *common.Credential
}

func NewQQCloudLHDriver(secretId, secretKey string) Driver {

	return &QQCloudLHDriver{
		credential: common.NewCredential(
			secretId,
			secretKey,
		),
	}
}

func (driver *QQCloudLHDriver) ListRegions() ([]*Region, error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, err := lighthouse.NewClient(driver.credential, "", cpf)
	if err != nil {
		return nil, err
	}

	request := lighthouse.NewDescribeRegionsRequest()
	response, err := client.DescribeRegions(request)
	if err != nil {
		return nil, err
	}

	regions := []*Region{}
	for _, lhregion := range response.Response.RegionSet {
		regions = append(regions, &Region{
			Name:            *lhregion.RegionName,
			Region:          *lhregion.Region,
			State:           RegionAvaliable,
			IsChinaMainland: *lhregion.IsChinaMainland,
		})
	}

	return regions, nil
}

func (driver *QQCloudLHDriver) ListZones(region string) ([]*Zone, error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeZonesRequest()

	response, err := client.DescribeZones(request)
	if err != nil {
		return nil, err
	}

	zones := []*Zone{}

	for _, lhzone := range response.Response.ZoneInfoSet {
		zones = append(zones, &Zone{
			Name: *lhzone.ZoneName,
			Zone: *lhzone.Zone,
		})
	}

	return zones, nil
}

func (driver *QQCloudLHDriver) ListInstances(region string) ([]*InstanceInfo, error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeInstancesRequest()
	request.Limit = common.Int64Ptr(100)

	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}

	instances := []*InstanceInfo{}

	for _, lhinstance := range response.Response.InstanceSet {

		instances = append(instances, lhRespInstaceToInstaceInfo(region, lhinstance))
	}

	return instances, nil

}

func (driver *QQCloudLHDriver) InstanceInfo(region, instanceID string) (*InstanceInfo, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeInstancesRequest()

	request.InstanceIds = common.StringPtrs([]string{instanceID})

	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}

	if len(response.Response.InstanceSet) > 0 {
		return lhRespInstaceToInstaceInfo(region, response.Response.InstanceSet[0]), nil
	}

	return nil, fmt.Errorf("区域[%s]下不存在实例[%s]", region, instanceID)
}

func lhRespInstaceToInstaceInfo(region string, lhinstance *lighthouse.Instance) *InstanceInfo {
	instanceInfo := &InstanceInfo{
		ID:           *lhinstance.InstanceId,
		Name:         *lhinstance.InstanceName,
		Region:       region,
		Zone:         *lhinstance.Zone,
		Cpu:          int(*lhinstance.CPU),
		Memory:       int(*lhinstance.Memory),
		OSName:       *lhinstance.OsName,
		Platform:     *lhinstance.Platform,
		PlatformType: *lhinstance.PlatformType,
		Disk:         int(*lhinstance.SystemDisk.DiskSize),
		PublicIP:     *lhinstance.PublicAddresses[0],
		PeivateIP:    *lhinstance.PrivateAddresses[0],
		Bandwidth:    int(*lhinstance.InternetAccessible.InternetMaxBandwidthOut),
		State:        InstanceState(*lhinstance.InstanceState),
	}

	if len(lhinstance.PublicAddresses) > 0 {
		instanceInfo.PublicIP = *lhinstance.PublicAddresses[0]
	}

	if len(lhinstance.PrivateAddresses) > 0 {
		instanceInfo.PeivateIP = *lhinstance.PrivateAddresses[0]
	}

	if lhinstance.CreatedTime != nil {
		instanceInfo.CreatedTime, _ = time.Parse(time.RFC3339, *lhinstance.CreatedTime)
	}

	if lhinstance.ExpiredTime != nil {
		instanceInfo.ExpiredTime, _ = time.Parse(time.RFC3339, *lhinstance.ExpiredTime)
	}

	return instanceInfo
}

func (driver *QQCloudLHDriver) StopInstances(region string, instanceIDs []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewStopInstancesRequest()

	request.InstanceIds = common.StringPtrs(instanceIDs)

	_, err := client.StopInstances(request)
	return err
}
func (driver *QQCloudLHDriver) StartInstances(region string, instanceIDs []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewStartInstancesRequest()

	request.InstanceIds = common.StringPtrs(instanceIDs)

	_, err := client.StartInstances(request)
	return err
}
func (driver *QQCloudLHDriver) RestartInstances(region string, instanceIDs []string) error {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewRebootInstancesRequest()

	request.InstanceIds = common.StringPtrs(instanceIDs)

	_, err := client.RebootInstances(request)
	return err

}
func (driver *QQCloudLHDriver) TerminateInstances(region string, instanceIDs []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewTerminateInstancesRequest()
	request.InstanceIds = common.StringPtrs(instanceIDs)

	_, err := client.TerminateInstances(request)
	return err

}

func (driver *QQCloudLHDriver) ResetPassword(region string, instanceIDs []string, username, password string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewResetInstancesPasswordRequest()

	request.InstanceIds = common.StringPtrs(instanceIDs)
	request.Password = common.StringPtr(password)
	if username != "" {
		request.UserName = common.StringPtr("username")
	}

	_, err := client.ResetInstancesPassword(request)
	return err
}

func (driver *QQCloudLHDriver) ResetInstances(region string, instanceIDs []string, BlueprintId string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewResetInstanceRequest()

	var err error

	request.BlueprintId = common.StringPtr(BlueprintId)
	for _, instanceID := range instanceIDs {
		request.InstanceId = common.StringPtr(instanceID)
		_, err0 := client.ResetInstance(request)
		if err0 != nil {
			err = err0
		}
	}

	return err
}

func (driver *QQCloudLHDriver) InstancesTrafficPackages(region string, instanceIDs []string) ([]*TrafficPackage, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeInstancesTrafficPackagesRequest()

	request.InstanceIds = common.StringPtrs(instanceIDs)

	response, err := client.DescribeInstancesTrafficPackages(request)
	if err != nil {
		return nil, err
	}

	packages := []*TrafficPackage{}

	for _, lhpackage := range response.Response.InstanceTrafficPackageSet {
		packages = append(packages, &TrafficPackage{
			InstanceId: *lhpackage.InstanceId,
			Used:       *lhpackage.TrafficPackageSet[0].TrafficUsed,
			Total:      *lhpackage.TrafficPackageSet[0].TrafficUsed,
			Remaining:  *lhpackage.TrafficPackageSet[0].TrafficPackageRemaining,
		})
	}

	return packages, nil
}

func (driver *QQCloudLHDriver) ListSnapshots(region, instanceID string) ([]*SnapShot, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeSnapshotsRequest()
	request.Limit = common.Int64Ptr(100)

	request.Filters = []*lighthouse.Filter{
		{
			Name:   common.StringPtr("instance-id"),
			Values: common.StringPtrs([]string{instanceID}),
		},
	}

	response, err := client.DescribeSnapshots(request)
	if err != nil {
		return nil, err
	}

	snapShots := []*SnapShot{}

	for _, lhsnapshot := range response.Response.SnapshotSet {
		snapShots = append(snapShots, lhRespSnapshotToSnapshotInfo(region, lhsnapshot))
	}

	return snapShots, nil

}

func (driver *QQCloudLHDriver) SnapshotInfo(region, snapshotID string) (*SnapShot, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeSnapshotsRequest()

	request.SnapshotIds = common.StringPtrs([]string{"snapshotID"})

	response, err := client.DescribeSnapshots(request)
	if err != nil {
		return nil, err
	}

	if len(response.Response.SnapshotSet) > 0 {
		return lhRespSnapshotToSnapshotInfo(region, response.Response.SnapshotSet[0]), nil
	}

	return nil, fmt.Errorf("区域[%s]下不存在快照[%s]", region, snapshotID)
}

func lhRespSnapshotToSnapshotInfo(region string, lhsnapshot *lighthouse.Snapshot) *SnapShot {
	snapshot := &SnapShot{
		SnapShot: *lhsnapshot.SnapshotId,
		Name:     *lhsnapshot.SnapshotName,
		Percent:  int(*lhsnapshot.Percent),
		State:    SnapShotState(*lhsnapshot.SnapshotState),
	}

	if lhsnapshot.CreatedTime != nil {
		snapshot.CreatedTime, _ = time.Parse(time.RFC3339, *lhsnapshot.CreatedTime)
	}
	return snapshot
}

func (driver *QQCloudLHDriver) DeleteSnapshots(region string, snapshotIDs []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDeleteSnapshotsRequest()

	request.SnapshotIds = common.StringPtrs(snapshotIDs)

	_, err := client.DeleteSnapshots(request)
	return err
}
func (driver *QQCloudLHDriver) CreateSnapshot(region, instanceID, name string) (*SnapShot, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewCreateInstanceSnapshotRequest()

	request.InstanceId = common.StringPtr(instanceID)
	if name != "" {
		request.SnapshotName = common.StringPtr(name)
	}

	response, err := client.CreateInstanceSnapshot(request)
	if err != nil {
		return nil, err
	}
	return &SnapShot{
		SnapShot: *response.Response.SnapshotId,
	}, nil
}

func (driver *QQCloudLHDriver) ApplySnapshot(region, instanceID, snapshotID string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewApplyInstanceSnapshotRequest()

	request.InstanceId = common.StringPtr(instanceID)
	request.SnapshotId = common.StringPtr(snapshotID)

	_, err := client.ApplyInstanceSnapshot(request)
	return err
}

func (driver *QQCloudLHDriver) ListBlueprints(region string, platformType PlatformType, blueprintType BlueprintType) ([]*Blueprint, error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeBlueprintsRequest()
	request.Limit = common.Int64Ptr(100)

	request.Filters = []*lighthouse.Filter{}
	if platformType != AllPlatform {
		request.Filters = append(request.Filters, &lighthouse.Filter{
			Name:   common.StringPtr("platform-type"),
			Values: common.StringPtrs([]string{string(platformType)}),
		})
	}

	if blueprintType != AllBlueprint {
		request.Filters = append(request.Filters, &lighthouse.Filter{
			Name:   common.StringPtr("blueprint-type"),
			Values: common.StringPtrs([]string{string(blueprintType)}),
		})
	}

	response, err := client.DescribeBlueprints(request)
	if err != nil {
		return nil, err
	}

	blueprints := []*Blueprint{}

	for _, lhblueprint := range response.Response.BlueprintSet {
		blueprints = append(blueprints, lhRespBlueprintToBlueprintInfo(lhblueprint))
	}

	return blueprints, nil
}

func (driver *QQCloudLHDriver) BlueprintInfo(region, blueprintID string) (*Blueprint, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeBlueprintsRequest()

	request.BlueprintIds = common.StringPtrs([]string{blueprintID})

	response, err := client.DescribeBlueprints(request)
	if err != nil {
		return nil, err
	}

	if len(response.Response.BlueprintSet) > 0 {
		return lhRespBlueprintToBlueprintInfo(response.Response.BlueprintSet[0]), nil
	}

	return nil, fmt.Errorf("区域[%s]下不存在镜像[%s]", region, blueprintID)
}

func lhRespBlueprintToBlueprintInfo(lhblueprint *lighthouse.Blueprint) *Blueprint {
	return &Blueprint{
		Blueprint:        *lhblueprint.BlueprintId,
		Name:             *lhblueprint.BlueprintName,
		Description:      *lhblueprint.Description,
		OsName:           *lhblueprint.OsName,
		RequiredDiskSize: *lhblueprint.RequiredSystemDiskSize,
		RequiredMemory:   *lhblueprint.RequiredMemorySize,
		State:            BlueprintState(*lhblueprint.BlueprintState),
	}
}

func (driver *QQCloudLHDriver) DeleteBlueprints(region string, blueprintIDs []string) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDeleteBlueprintsRequest()

	request.BlueprintIds = common.StringPtrs(blueprintIDs)

	_, err := client.DeleteBlueprints(request)
	return err
}
func (driver *QQCloudLHDriver) CreateBlueprint(region, instanceId, name, desctiprtion string) (*Blueprint, error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewCreateBlueprintRequest()

	request.BlueprintName = common.StringPtr(name)
	if desctiprtion != "" {
		request.Description = common.StringPtr(desctiprtion)
	}

	request.InstanceId = common.StringPtr(instanceId)

	response, err := client.CreateBlueprint(request)
	if err != nil {
		return nil, err
	}

	return &Blueprint{
		Blueprint: *response.Response.BlueprintId,
	}, nil
}

func (driver *QQCloudLHDriver) ListFirewallRules(region string, instanceID string) ([]*FirewallRule, error) {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDescribeFirewallRulesRequest()

	request.InstanceId = common.StringPtr(instanceID)
	request.Limit = common.Int64Ptr(100)

	response, err := client.DescribeFirewallRules(request)
	if err != nil {
		return nil, err
	}

	rules := []*FirewallRule{}

	for _, lhrole := range response.Response.FirewallRuleSet {
		rules = append(rules, &FirewallRule{
			Protocol:    FirewallRuleProtocol(*lhrole.Protocol),
			Port:        *lhrole.Port,
			CidrBlock:   *lhrole.CidrBlock,
			Action:      FirewallRuleAction(*lhrole.Action),
			Description: *lhrole.FirewallRuleDescription,
		})
	}

	return rules, nil
}
func (driver *QQCloudLHDriver) AddFirewallRules(region string, instanceID string, roles []*FirewallRule) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewCreateFirewallRulesRequest()

	request.InstanceId = common.StringPtr(instanceID)
	request.FirewallRules = []*lighthouse.FirewallRule{}

	for _, role := range roles {
		request.FirewallRules = append(request.FirewallRules, &lighthouse.FirewallRule{
			Protocol:                common.StringPtr(string(role.Protocol)),
			Port:                    common.StringPtr(role.Port),
			CidrBlock:               common.StringPtr(role.CidrBlock),
			Action:                  common.StringPtr(string(role.Action)),
			FirewallRuleDescription: common.StringPtr(role.Description),
		})
	}
	_, err := client.CreateFirewallRules(request)
	return err

}
func (driver *QQCloudLHDriver) UpdateFirewallRules(region string, instanceID string, roles []*FirewallRule) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewModifyFirewallRulesRequest()

	request.InstanceId = common.StringPtr(instanceID)
	request.FirewallRules = []*lighthouse.FirewallRule{}

	for _, role := range roles {
		request.FirewallRules = append(request.FirewallRules, &lighthouse.FirewallRule{
			Protocol:                common.StringPtr(string(role.Protocol)),
			Port:                    common.StringPtr(role.Port),
			CidrBlock:               common.StringPtr(role.CidrBlock),
			Action:                  common.StringPtr(string(role.Action)),
			FirewallRuleDescription: common.StringPtr(role.Description),
		})
	}
	_, err := client.ModifyFirewallRules(request)
	return err
}
func (driver *QQCloudLHDriver) DeleteFirewallRules(region string, instanceID string, roles []*FirewallRule) error {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, _ := lighthouse.NewClient(driver.credential, region, cpf)

	request := lighthouse.NewDeleteFirewallRulesRequest()

	request.InstanceId = common.StringPtr(instanceID)
	request.FirewallRules = []*lighthouse.FirewallRule{}

	for _, role := range roles {
		request.FirewallRules = append(request.FirewallRules, &lighthouse.FirewallRule{
			Protocol:                common.StringPtr(string(role.Protocol)),
			Port:                    common.StringPtr(role.Port),
			CidrBlock:               common.StringPtr(role.CidrBlock),
			Action:                  common.StringPtr(string(role.Action)),
			FirewallRuleDescription: common.StringPtr(role.Description),
		})
	}
	_, err := client.DeleteFirewallRules(request)
	return err
}
