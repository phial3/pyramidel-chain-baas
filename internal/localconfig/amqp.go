package localconfig

type (
	AMQP struct {
		Host        string `yaml:"host"`         // 主机地址,127.0.0.1
		Port        int    `yaml:"port"`         // 端口,5672
		User        string `yaml:"user"`         // username,txhy
		Password    string `yaml:"password"`     // pw,txhy2022.com
		Vhost       string `yaml:"vhost"`        // 虚拟主机地址
		Queue       string `yaml:"queue"`        // 队列名称
		ContentType string `yaml:"content_type"` //
	}
)
