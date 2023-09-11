package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
	"github.com/toolkits/pkg/runner"
	"github.com/urfave/cli/v2"

	"github.com/didi/nightingale/v5/src/server"
	"github.com/didi/nightingale/v5/src/webapi"
	"github.com/wumansgy/goEncrypt"
)

var privateKey = []byte(`-----BEGIN  WUMAN RSA PRIVATE KEY -----
MIIEpgIBAAKCAQEA29fttcEDUvhGJEhQXEIH8blptZRF5itec4GtEGtkSr4Wjmsf
2o2XKOr6YEbTOeDA/DdnDSbVzK2ZUscqyBxbKGwI/Bpv9l5K/sh9+Oj2Y8YH53+X
kqRSGvmhHqolhb+gcfH+FKG5IflGuiOREs4h02TVmPAFPTmZjYBeVexJgmPodGPO
e36QVnMeOG8tHOFxItkMvJUpilzs85xdHqTTjWCtk/SjHrp5NGSkHSmionOtrFik
sS/gTX0EzrptmAGHTjZV0NX7Nu8Ma45rVdMRwXrDPbk0yR0iFdBEZ1ceGsNg2Vjr
Z3LCZi3zO+ieA7sBjHARHai5MuFlh9KJ8+YkwwIDAQABAoIBAQCyBCtMXbqfWMMT
ZisMSbu9FPJwQlxHgR6+UWceQJe5nisNr9jfVH/udje/9hncaA5dLU+Y6rV9Q6U/
zl7qI2v9U14DJjU7PidkIF1BTQMWz6he4IaQC9cgWLsK5aP0pbL6EYY4lqwewodu
+pXisF/bmW8MpG7ZoOaiGixJT0hG97aS3YD506RqdnK4a9yG5ycoVUZTzFTIM+aq
MCTOWJLbHJBnY52v7Rqaor1jZ0o+C/Cykbts25VHWZ9ygxBfI/S75jTg/zibpqbo
TnICGBqzzfAsThHYiji3ZgEF2bSadxQZ956Rvm5Dlhk0A6ylwKl1gJTAxyZNRjx2
zogLHNFhAoGBAP3TmQSpyYtjr8xdD+AsinMi3q3p/6+FPBh5wrV9Dvd5aSrKCYnU
j/9LOYmCP6JfKc90K3L0PKTfS1LDZXHoeYpzxB0uD2iTIkYQKraUp7rzqLBMkdTB
nTuOByqnTB45WlWlfy6m/8CO+0r5DkrO1fU8gWkg3+bUQoUTPOhuFkqdAoGBAN25
1oVmFyDvxw62WZ95SPktqzmErFH/X+7wpGInUb0vxnGJ0SP2bVgMLyL9Ecz8p5hC
ZEZVbo/WFs17hqz8Z6Xx4Gq7aoredpsZswpApdODs5sQf72LXj4Cuf1iwMGIyQxa
8hElV9uRIeKfEG2E/HBLLPO8qhGBCsWIT1Jms97fAoGBAPQFIf2usTkVXCPvb9zH
VU79PfEanhoCz9SD8mGCWgomqaleVK8yMEFx8120XzLdpBdyCndYQJkMpqBpgzRw
F7C4PNkEuAGEOhX7YuTmox4DM7BR3H0aqetgTpl9/pqr7qGaGlwiZoubqhDYwRnA
IUfDpHIKDdcfRtgit5KIi1utAoGBAKqoAK8IFsEpDGMMgwq1hS8UsXdB4If0MNht
q3hInycoAGsfEjPF1f8w0Y7yjaLiy/PrFdb0pnZa544ch1nZo8Ub2AkOW0CrXUqf
iyhW/ctA0RqGpmszO8QqwRB/07CiIWw7C5maznaWzCfrGe/RraKYme63xYZXdfz3
n2Xi2oqtAoGBAPB2SL1UafzZiC70e+NeDeaCLTCorCvIrN73RngYXU1OKUynrcjL
xOs8yVFs5cK0sNwEkiozcldfOfU2j70tyrGz+txi5+Db6ex5VXmEKSqdZXRDtqqc
KKRFfpuebqdaR50SDa4Lq6JbqYtwg0CLZeju4Mq41i9F2p4myqVDUsZs
-----END  WUMAN RSA PRIVATE KEY -----
`)

// VERSION go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "not specified"

func hello() {
	fmt.Println("hello")
}
func main() {
	ciphertext, err := base64.StdEncoding.DecodeString(`mFWBdT4Y70ZNEQ7PVIFKwbkefufu52WGXYLrW0Vk1XuajrrEE54dqj4VK2yuGIeMq5bHKAdkDnACB2ABzHLQuobTDpkS0Nj5AlJvwbRDV3pOCB1x0q3aqEooTppeMs8P/WG3YCRDTQPWgZISPsFBQVT1tk77BiImcY4SZM9IL0B4TFUKS9sShnjAebxmJkj8jfYYh7gNzUY0YMvOV6HuiT5C0RsbTe1jwMyN87QEwvpvuPelkeQ8LX1AG+qsn2q4TvOYEKCNfNnePjMIQ/5MlesledwiqUpc/YtY3qj4Qx+8b5luaQ6kyu+zyOXV/A0XjjxIxqLWKU8eAl7eA3o72Q==`)
	if err != nil {
		return
	}

	plaintext, err := goEncrypt.RsaDecrypt(ciphertext, privateKey)
	if err != nil {
		return
	}
	c := cron.New()
	c.AddFunc("@every 1s", hello)
	c.Start()
	defer c.Stop()
	var ch = make(chan int)
	<-ch

	fmt.Println("明文：", string(plaintext)) // test
	app := cli.NewApp()
	app.Name = "n9e"
	app.Version = VERSION
	app.Usage = "Nightingale, enterprise prometheus management"
	app.Commands = []*cli.Command{
		newWebapiCmd(),
		newServerCmd(),
	}
	app.Run(os.Args)
}

func newWebapiCmd() *cli.Command {
	return &cli.Command{
		Name:  "webapi",
		Usage: "Run webapi",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "specify configuration file(.json,.yaml,.toml)",
			},
		},
		Action: func(c *cli.Context) error {
			printEnv()

			var opts []webapi.WebapiOption
			if c.String("conf") != "" {
				opts = append(opts, webapi.SetConfigFile(c.String("conf")))
			}
			opts = append(opts, webapi.SetVersion(VERSION))

			webapi.Run(opts...)
			return nil
		},
	}
}

func newServerCmd() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Run server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "specify configuration file(.json,.yaml,.toml)",
			},
		},
		Action: func(c *cli.Context) error {
			printEnv()

			var opts []server.ServerOption
			if c.String("conf") != "" {
				opts = append(opts, server.SetConfigFile(c.String("conf")))
			}
			opts = append(opts, server.SetVersion(VERSION))

			server.Run(opts...)
			return nil
		},
	}
}

func printEnv() {
	runner.Init()
	fmt.Println("runner.cwd:", runner.Cwd)
	fmt.Println("runner.hostname:", runner.Hostname)
	fmt.Println("runner.fd_limits:", runner.FdLimits())
	fmt.Println("runner.vm_limits:", runner.VMLimits())
}
