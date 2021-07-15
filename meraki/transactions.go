package meraki

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"stellar.af/meraki-reboot/types"
	"stellar.af/meraki-reboot/util"
)

func findCorresponding(result gjson.Result, aKey, aValue, bKey string) *string {
	var bValue string = ""
	result.ForEach(func(key, value gjson.Result) bool {
		a := value.Get(aKey).String()
		if a == aValue {
			bValue = value.Get(bKey).String()
			return false
		}
		return true
	})
	if bValue != "" {
		return &bValue
	}
	return nil
}

func GetOrganizationID(orgName string) (orgID string, err error) {
	allOrgs, err := MerakiRequest("GET", "/api/v1/organizations", emptyQuery)
	util.Check("Error fetching organizations from Meraki dashboard", err)
	matching := findCorresponding(allOrgs, "name", orgName, "id")
	if matching != nil {
		return *matching, nil
	}
	return "", fmt.Errorf("Unable to find matching Meraki organization for '%s'\n", orgName)

}

func GetNetworkID(orgID string, networkName string) (networkID string, err error) {
	allNets, err := MerakiRequest("GET", fmt.Sprintf("/api/v1/organizations/%s/networks", orgID), emptyQuery)
	util.Check("Error getting networks for organization ID '%s'", err, orgID)
	matching := findCorresponding(allNets, "name", networkName, "id")
	if matching != nil {
		return *matching, nil
	}
	return "", fmt.Errorf("Unable to find network matching '%s' in organization '%s'", networkName, orgID)
}

func GetNetworkDevices(networkID string, exclusions []string) (devices []*types.MerakiDevice, err error) {
	allDevices, err := MerakiRequest("GET", fmt.Sprintf("/api/v1/networks/%s/devices", networkID), emptyQuery)
	util.Check("Error getting devices for network ID '%s'", err, networkID)
	allDevices.ForEach(func(key, value gjson.Result) bool {
		var device types.MerakiDevice
		e := json.Unmarshal([]byte(value.Raw), &device)
		if e != nil {
			err = e
			return false
		}
		hasExclusions := util.CompareArrays(exclusions, device.Tags)
		if !hasExclusions {
			devices = append(devices, &device)
		}
		return true
	})
	return devices, err
}

func GetDevice(serial string) *types.MerakiDevice {
	d, err := MerakiRequest("GET", fmt.Sprintf("/api/v1/devices/%s", serial), emptyQuery)
	util.Check("Error getting device with serial number '%s'", err, serial)
	var device types.MerakiDevice
	err = json.Unmarshal([]byte(d.Raw), &device)
	util.Check("Error parsing device with serial number '%s'", err, serial)
	return &device

}

func RebootDevice(serial string) (success bool, err error) {
	res, err := MerakiRequest("POST", fmt.Sprintf("/api/v1/devices/%s/reboot", serial), emptyQuery)
	success = res.Get("success").Bool()
	return success, err
}
