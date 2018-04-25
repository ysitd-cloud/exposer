package manager

type Connection struct {
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
	Port     int    `json:"port"`
}
