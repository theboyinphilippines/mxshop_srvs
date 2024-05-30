package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host        string `mapstructure:"host" json:"host"`
	Port        uint64 `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId"`
	DataId      string `mapstructure:"dataId" json:"dataId"`
	Group       string `mapstructure:"group" json:"group"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"` //注册到consul中的服务名称
	Host       string       `mapstructure:"host" json:"host"` //服务的host
	Tags       []string     `mapstructure:"tags" json:"tags"` //服务的tags
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}
