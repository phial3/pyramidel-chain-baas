package localconfig

type Serve struct {
	Mode string `json:"mode" yaml:"mode,omitempty"`
	Port string `json:"port" yaml:"port"`
}
