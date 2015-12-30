package locator

type Request struct {
	HomeMobileCountryCode uint16             `json:"homeMobileCountryCode,omitempty"` // The mobile country code stored on the SIM card (100-999).
	HomeMobileNetworkCode uint16             `json:"homeMobileNetworkCode,omitempty"` // The mobile network code stored on the SIM card (0-32767).
	RadioType             string             `json:"radioType,omitempty"`             // The mobile radio type. Supported values are lte, gsm, cdma, and wcdma.
	Carrier               string             `json:"carrier,omitempty"`               // The clear text name of the cell carrier / operator.
	ConsiderIp            bool               `json:"considerIp"`                      // Should the clients IP address be used to locate it, defaults to true.
	CellTowers            []*CellTower       `json:"cellTowers,omitempty"`            // Array of cell towers
	WifiAccessPoints      []*WifiAccessPoint `json:"wifiAccessPoints,omitempty"`      // Array of wifi access points
	IPAddress             string             `json:"ipaddress,omitempty"`             // Client IP Address
	Fallbacks             *Fallbacks         `json:"fallbacks,omitempty"`             // The fallback section is a custom addition to the GLS API.
}

type Fallbacks struct {
	LAC bool `json:"lacf"` // If no exact cell match can be found, fall back from exact cell position estimates to more coarse grained cell location area estimates, rather than going directly to an even worse GeoIP based estimate.
	IP  bool `json:"ipf"`  // If no position can be estimated based on any of the provided data points, fall back to an estimate based on a GeoIP database based on the senders IP address at the time of the query.
}

type CellTower struct {
	MobileCountryCode uint16 `json:"mobileCountryCode"`        // The mobile country code.
	MobileNetworkCode uint16 `json:"mobileNetworkCode"`        // The mobile network code.
	LocationAreaCode  uint16 `json:"locationAreaCode"`         // The location area code for GSM and WCDMA networks. The tracking area code for LTE networks.
	CellId            uint32 `json:"cellId"`                   // The cell id or cell identity.
	SignalStrength    int16  `json:"signalStrength,omitempty"` // The signal strength for this cell network, either the RSSI or RSCP.
	Age               uint32 `json:"age,omitempty"`            // The number of milliseconds since this networks was last detected.
	TimingAdvance     uint8  `json:"timingAdvance,omitempty"`  // The timing advance value for this cell network.
}

type WifiAccessPoint struct {
	MacAddress         string `json:"macAddress"`                   // The BSSID of the WiFi network.
	SignalStrength     int16  `json:"signalStrength,omitempty"`     // The received signal strength (RSSI) in dBm.
	Age                uint32 `json:"age,omitempty"`                // The number of milliseconds since this network was last detected.
	Channel            uint8  `json:"channel,omitempty"`            // The WiFi channel, often 1 - 13 for networks in the 2.4GHz range.
	SignalToNoiseRatio uint16 `json:"signalToNoiseRatio,omitempty"` // The current signal to noise ratio measured in dB.
}

type Response struct {
	Location Point   `json:"location"` // The user’s estimated latitude and longitude, in degrees.
	Accuracy float64 `json:"accuracy"` // The accuracy of the estimated location, in meters.
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
