package config

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/koding/multiconfig"

	"github.com/didi/nightingale/v5/src/pkg/httpx"
	"github.com/didi/nightingale/v5/src/pkg/logx"
	"github.com/didi/nightingale/v5/src/server/naming"
	"github.com/didi/nightingale/v5/src/server/reader"
	"github.com/didi/nightingale/v5/src/server/writer"
	"github.com/didi/nightingale/v5/src/storage"
)

var (
	C    = new(Config)
	once sync.Once
)

func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			handled := false

			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
				handled = true
			}
			if strings.HasSuffix(fpath, "conf") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
				handled = true
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
				handled = true
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
				handled = true
			}

			if !handled {
				fmt.Println("config file invalid, valid file exts: .conf,.yaml,.toml,.json")
				os.Exit(1)
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)

		if C.Heartbeat.IP == "" {
			// auto detect
			C.Heartbeat.IP = fmt.Sprint(GetOutboundIP())

			if C.Heartbeat.IP == "" {
				fmt.Println("heartbeat ip auto got is blank")
				os.Exit(1)
			}
		}

		C.Heartbeat.Endpoint = fmt.Sprintf("%s:%d", C.Heartbeat.IP, C.HTTP.Port)
		C.Heartbeat.Cluster = C.ClusterName

		C.Alerting.RedisPub.ChannelKey = C.Alerting.RedisPub.ChannelPrefix + C.ClusterName

		fmt.Println("heartbeat.ip:", C.Heartbeat.IP)
		fmt.Printf("heartbeat.interval: %dms\n", C.Heartbeat.Interval)
	})
}

type Config struct {
	RunMode     string
	ClusterName string
	Log         logx.Config
	HTTP        httpx.Config
	BasicAuth   gin.Accounts
	Heartbeat   naming.HeartbeatConfig
	Alerting    Alerting
	NoData      NoData
	Redis       storage.RedisConfig
	Gorm        storage.Gorm
	MySQL       storage.MySQL
	Postgres    storage.Postgres
	WriterOpt   writer.GlobalOpt
	Writers     []writer.Options
	Reader      reader.Options
	Ibex        Ibex
}

type Alerting struct {
	NotifyScriptPath  string
	NotifyConcurrency int
	RedisPub          RedisPub
}

type RedisPub struct {
	Enable        bool
	ChannelPrefix string
	ChannelKey    string
}

type NoData struct {
	Metric   string
	Interval int64
}

type Ibex struct {
	Address       string
	BasicAuthUser string
	BasicAuthPass string
	Timeout       int64
}

func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("auto get outbound ip fail:", err)
		os.Exit(1)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}