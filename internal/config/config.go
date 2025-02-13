package config

import "github.com/spf13/viper"

type Config struct {
	App   AppConfig
	MySql MySqlConfig
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
	JWTExpire int    `mapstructure:"jwt_expire"`
}

type MySqlConfig struct {
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	Database              string `mapstructure:"database"`
	MaxIdleConns          int    `mapstructure:"max_idle_conns"`
	MaxOpenConns          int    `mapstructure:"max_open_conns"`
	ConnectionMaxLifetime int    `mapstructure:"connection_max_lifetime"`
}

var (
	GlobalConfig Config
)

func Init(configPath string) error {
	// 指定配置文件的路径
	viper.SetConfigFile(configPath)
	// 自动绑定环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
