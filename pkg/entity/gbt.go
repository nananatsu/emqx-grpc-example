package entity

type GBTResult struct {
	Success int `json:"success"`
}

type GBTInstallInfo struct {
	Company          string  `json:"company"`
	Serial           string  `json:"serial"`
	Name             string  `json:"name"`
	CategoryName     string  `json:"categoryName"`
	NetType          string  `json:"netType"`
	Client           string  `json:"client"`
	Address          string  `json:"address"`
	DeviceAddress    string  `json:"deviceAddress"`
	ProductingDate   string  `json:"productingDate"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	ActiveTime       string  `json:"activeTime"`
	OrganizationName string  `json:"organizationName"`
	GbSerial         string  `json:"gbSerial"`
	Responsible      string  `json:"responsible"`
	Mobile           string  `json:"mobile"`
}

type GBTStatusInfo struct {
	Company          string `json:"company"`
	Serail           string `json:"serail"`
	Status           int    `json:"status"`
	SignalStrength   int
	BatteryPercetage int
	Datetime         string `json:"datetime"`
}
