package localconfig

type Serve struct {
	Mode string `json:"mode" yaml:"mode" default:"debug"`
	Port string `json:"port" yaml:"port"`
}

func checkMode(mode string) {
	switch mode {
	case "debug":
	case "release":
	case "test":
		return
	default:
		panic("unknown mode")
	}
}
