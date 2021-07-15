package types

type MerakiDevice = struct {
	Lat         float64  `json:"lat"`
	Lng         float64  `json:"lng"`
	Address     string   `json:"address"`
	Serial      string   `json:"serial"`
	Mac         string   `json:"mac"`
	LanIP       string   `json:"lanIp"`
	Tags        []string `json:"tags"`
	URL         string   `json:"url"`
	NetworkID   string   `json:"networkId"`
	Name        string   `json:"name"`
	Model       string   `json:"model"`
	Firmware    string   `json:"firmware"`
	FloorPlanID string   `json:"floorPlanId"`
}
