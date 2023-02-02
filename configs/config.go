package configs

type (
	MySQLConfig struct {
		Host         string `mapstructure:"host"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		DB           string `mapstructure:"dbname"`
		Port         int    `mapstructure:"port"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
	}

	RedisConfig struct {
		Host         string `mapstructure:"host"`
		Password     string `mapstructure:"password"`
		Port         int    `mapstructure:"port"`
		DB           int    `mapstructure:"db"`
		PoolSize     int    `mapstructure:"pool_size"`
		MinIdleConns int    `mapstructure:"min_idle_conns"`
	}
)
