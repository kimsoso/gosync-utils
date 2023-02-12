package config

import (
	"btsync-utils/libs/consts"
	bnet "btsync-utils/libs/net"
	"btsync-utils/libs/utils"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	WatchDir   string                `yaml:"watchDir"`
	RemoteHost []map[string][]string `yaml:"remoteHost"`
	SubNet     string                `yaml:"subnet"`
	HttpPort   int                   `yaml:"httpPort"`
	Port       uint                  `yaml:"port"`
	Log        bool                  `yaml:"log"`
	Verfile    string                `yaml:"verfile"`
}

type Conf struct {
	path  string
	maps  sync.Map
	viper *viper.Viper
}

func NewConf(path string) *Conf {
	if _, err := os.Stat(path); err != nil {
		panic("配置文件不存在:" + path)
	}

	out := &Conf{
		path:  path,
		viper: viper.New(),
	}
	out.setDefault()
	out.wathFile()

	return out
}

func (c *Conf) wathFile() {
	c.viper.SetConfigFile(c.path)
	c.viper.SetConfigType("yaml")

	if err := c.viper.ReadInConfig(); err != nil {
		log.Panic("配置文件错误,", err)
	} else if err := c.load(); err != nil {
		log.Panic("配置文件格式错误,", err)
	}

	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件被修改", time.Now().Format(consts.TIME_LAYOUT))
		c.load()
	})
}

func (c *Conf) setDefault() {
	c.viper.SetDefault("WatchDir", "files/")
	c.viper.SetDefault("SubNet", "192.168.1.0/24")
	c.viper.SetDefault("HttpPort", 8080)
	c.viper.SetDefault("Log", true)
	c.viper.SetDefault("Port", 9111)
	c.viper.SetDefault("MyIp", "")
}

func (c *Conf) load() error {
	confs := new(Config)

	if err := c.viper.Unmarshal(confs); err != nil {
		return err
	} else {
		if rpath, err := filepath.Abs(confs.WatchDir); err != nil {
			panic("监控目录不存在")
		} else if _, err := os.Stat(confs.WatchDir); err != nil {
			panic("监控目录不存在")
		} else {
			c.maps.Store("WatchDir", rpath)
		}

		c.maps.Store("HttpPort", confs.HttpPort)
		c.maps.Store("SubNet", confs.SubNet)
		c.maps.Store("Log", confs.Log)
		c.maps.Store("Port", confs.Port)

		if confs.RemoteHost != nil && len(confs.RemoteHost) > 0 {
			c.maps.Store("RemoteHost", confs.RemoteHost[0])
		}

		c.maps.Store("Verfile", confs.Verfile)

		if _, ipNet, err := net.ParseCIDR(confs.SubNet); err == nil {
			if ip, err := bnet.LocalIpByIpNet(*ipNet); err == nil {
				// log.Println("config: my local ip is:", ip)
				c.maps.Store("Ip", ip)
			}
		}

		if _, ipsetted := c.maps.Load("Ip"); !ipsetted {
			if ip, err := utils.GetLocalIp(); err == nil {
				c.maps.Store("Ip", ip)
				log.Println("parse cidr error:", err)
			} else {
				log.Fatal("can't get local ip")
			}
		}

		log.Println("below is configs <<")

		c.maps.Range(func(k, v any) bool {
			log.Println(k.(string), ":", v)
			return true
		})
		log.Println(">> above is configs")

		return nil
	}
}

func (c *Conf) MyWatchDir() string {
	if out, ok := c.maps.Load("WatchDir"); ok {
		return out.(string)
	} else {
		return ""
	}
}

func (c *Conf) MyHttpPort() int {
	if out, ok := c.maps.Load("HttpPort"); ok {
		return out.(int)
	} else {
		return 0
	}
}
func (c *Conf) MyLog() bool {
	if out, ok := c.maps.Load("Log"); ok {
		return out.(bool)
	} else {
		return false
	}
}
func (c *Conf) MyPort() uint {
	if out, ok := c.maps.Load("Port"); ok {
		return out.(uint)
	} else {
		return 0
	}
}
func (c *Conf) MyRemoteHost() map[string][]string {
	if out, ok := c.maps.Load("RemoteHost"); ok {
		return out.(map[string][]string)
	} else {
		return nil
	}
}
func (c *Conf) MyVerfile() string {
	if out, ok := c.maps.Load("Verfile"); ok {
		return out.(string)
	} else {
		return ""
	}
}

func (c *Conf) MySubNet() string {
	if out, ok := c.maps.Load("SubNet"); ok {
		return out.(string)
	} else {
		return ""
	}
}

func (c *Conf) MyIp() string {
	if out, ok := c.maps.Load("Ip"); ok {
		return out.(string)
	} else {
		return ""
	}
}

func (c *Conf) MyClientSubnets() (subnets []string) {
	subnets = []string{}
	for subnet := range c.MyRemoteHost() {
		subnets = append(subnets, subnet)
	}
	return
}

func (c *Conf) RealFilepath(refFilepath string) string {
	return filepath.Join(c.MyWatchDir(), refFilepath)
}

func (c *Conf) UploadingFilepath(refFilepath string) string {
	return c.RealFilepath(refFilepath) + ".bsup"
}

func (c *Conf) RelFilepath(realFilepath string) string {
	basePath, _ := filepath.Abs(c.MyWatchDir())
	realpath, _ := filepath.Abs(realFilepath)
	out, _ := filepath.Rel(basePath, realpath)
	return out
}
