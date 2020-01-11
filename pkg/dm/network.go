package dm

//NetworkInfo contains the machine active session ssid&bssid, the gateway info and the network cards
type NetworkInfo struct {
	SSID         string
	BSSID        string
	GatewayIP    string
	GatewayMAC   string
	NetworkPass  string
	Channel      string
	NetworkCards []NetworkCardInfo
}

//NetworkCardInfo represents a phisical network card in the local machine
type NetworkCardInfo struct {
	Name string
	Mac  string
}
