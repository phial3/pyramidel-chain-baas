package localconfig

type Mysql struct {
	Host      string `yaml:"host"`      // 主机地址,127.0.0.1
	Port      int    `yaml:"port"`      // 端口,3306
	User      string `yaml:"user"`      // username,root
	Password  string `yaml:"password"`  // pw,123456
	DB        string `yaml:"db"`        // 数据库,baas-go
	Charset   string `yaml:"charset"`   // 字符集,utf8_unicode_ci
	Parsetime bool   `yaml:"parsetime"` // true
	Loc       bool   `yaml:"loc"`       // 是否本地时区？,true
}
