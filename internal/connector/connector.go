// Package connector defines the APConnector interface for managing AP radio settings.
package connector

// ClientInfo holds AP-reported client stats.
type ClientInfo struct {
	RSSI      *int     `json:"rssi"`
	SNR       *int     `json:"snr"`
	Signal    *int     `json:"signal"`
	Channel   *int     `json:"channel"`
	Band      *string  `json:"band"`
	APName    *string  `json:"apName"`
	PowerSave *bool    `json:"powerSave"`
	RxRate    *int     `json:"rxRate"`
	TxRate    *int     `json:"txRate"`
	WifiMode  *string  `json:"wifiMode"`
	Activity  *int     `json:"activity"`
	LastSeen  *int64   `json:"lastSeen"`
	Raw       map[string]any
}

// RadioConfig describes desired radio settings.
type RadioConfig struct {
	Channel  int
	WidthMHz int
	TxPower  *int // nil = no change
}

// APConnector is the interface for AP management backends.
type APConnector interface {
	Name() string
	GetClient(mac string) (*ClientInfo, error)
	SetRadio(apMAC string, cfg RadioConfig) error
	ReconnectClient(mac string) error
}
