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

type ApolloConfig struct {
	AppID           string   `mapstructure:"appId" json:"appId"`
	Cluster         string   `mapstructure:"cluster" json:"cluster"`
	NameSpaceNames  []string `mapstructure:"nameSpaceNames" json:"nameSpaceNames"`
	MetaAddr        string   `mapstructure:"metaAddr" json:"metaAddr"`
	AccesskeySecret string   `mapstructure:"accesskeySecret" json:"accesskeySecret"`
}

type ServerConfig struct {
	Host       string       `mapstructure:"host" json:"host"`
	Name       string       `mapstructure:"name" json:"name"` //注册到consul中的服务名称
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}
