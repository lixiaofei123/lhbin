package cmd

import (
	"fmt"
	"os"
)

const RegionCommandName string = "region"

func init() {

	RegisterChildCommand(RegionCommandName, "查看地域列表", []string{})
	RegisterChildCommandOperator(RegionCommandName, "list", "列出地域列表", []string{}, SafeOperation(ListRegions))
}

func ListRegions() error {

	cdriver, err := parseAndGetDriverWithoutSV(os.Args[3:])
	if err != nil {
		return err
	}

	regions, err := cdriver.ListRegions()
	if err != nil {
		return err
	}

	for _, region := range regions {
		fmt.Println(region.Name, region.Region)
	}

	return nil
}
