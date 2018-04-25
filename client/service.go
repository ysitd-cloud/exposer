package client

type Service struct {
	ServiceName string `json:"service"`
	Port        int    `json:"port"`
}

type connection struct {
	Service  `json:",inline"`
	Hostname string `json:"hostname"`
}
