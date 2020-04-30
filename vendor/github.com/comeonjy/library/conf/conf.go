package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName   string     `yaml:"app_name"`
	HttpPort  int16      `yaml:"http_port"`
	PprofPort int16      `yaml:"pprof_port"`
	Redis     *RedisConf `yaml:"redis"`
	Mysql     *DbConfig  `yaml:"mysql"`
	Mail      *MailConf  `yaml:"mail"`
	AesKey    string     `yaml:"aes_key"`
}

type MailConf struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type RedisConf struct {
	Address     string `yaml:"address"`
	Database    int16  `yaml:"database"`
	Password    string `yaml:"password"`
	DialTimeout int64  `yaml:"dial_timeout"`
}

type DbConfig struct {
	Debug       bool   `yaml:"debug"`
	Port        int16  `yaml:"port"`
	MaxIdleConn int16  `yaml:"max_idle_conn"`
	MaxOpenConn int16  `yaml:"max_open_conn"`
	UserName    string `yaml:"username"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	DataBase    string `yaml:"database"`
}

var cfg *Config

// 获取配置对象
func New() *Config {
	return cfg
}

// 加载配置文件
func LoadConfig(confPath string) {
	readConfigFile(confPath)
}

// 读取配置文件
func readConfigFile(file string) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	temp := &Config{}
	if err := yaml.Unmarshal(body, temp); err != nil {
		log.Fatal(err)
	}
	cfg = temp
}

func (c *Config) String() string {
	b, err := json.Marshal(*c)
	if err != nil {
		return fmt.Sprintf("%+v", *c)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *c)
	}
	return out.String()
}
