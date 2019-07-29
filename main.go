package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	Logger *logrus.Logger
	App    *cli.App
	Cli    IOssCli
)

func init() {
	Logger = newLogger()
}

func main() {
	// 初始化 Cli
	Cli = &AliOssCli{}

	// 注册 App
	App := makeApp()
	registerCommands(App)
	registerFlags(App)

	err := App.Run(os.Args)
	if err != nil {
		Logger.Fatal(err)
	}
}

// 生成 cli
func makeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "oss-cli"
	app.Usage = "OSS Cli 管理工具"
	app.Version = env("APP_VERSION", "0.0.0")
	return app
}

// 注册 命令
func registerCommands(app *cli.App) {
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "test",
			Usage:  "测试配置是否正确",
			Action: Cli.Test,
		},
		{
			Name:   "add",
			Usage:  "上传文件到 OSS 中",
			Action: Cli.Add,
			Flags:  getAddFlags(),
		},
		{
			Name:   "list",
			Usage:  "查看文件列表",
			Action: Cli.ListFiles,
			Flags:  getListFlags(),
		},
	}
}

// 注册 Flags
func registerFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "key, k",
			Usage:  "账户 API `access_key`",
			EnvVar: "ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret, s",
			Usage:  "账户 API `access_key_secret`",
			EnvVar: "ACCESS_KEY_SECRET",
		},
		cli.StringFlag{
			Name:   "bucket, b",
			Usage:  "存储空间 Bucket Name `bucket_name`",
			EnvVar: "BUCKET_NAME",
		},
		cli.StringFlag{
			Name:   "region, r",
			Usage:  "地域 Region Name `region name`",
			EnvVar: "REGION_NAME",
		},
		cli.BoolFlag{
			Name:   "internal, i",
			Usage:  "使用内网进行连接",
			EnvVar: "USE_INTERNAL",
		},
	}
}

func getListFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "prefix, p",
			Usage: "匹配文件前缀",
		},
		cli.StringFlag{
			Name:  "sort",
			Usage: "排序方式 asc, desc",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "最大显示条数, 0 不限制",
			Value: 0,
		},
		cli.BoolFlag{
			Name:  "dir",
			Usage: "只显示目录",
		},
		cli.StringFlag{
			Name:  "search, s",
			Usage: "通过文件名搜索文件",
		},
	}
}

func getAddFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "breakpoint, b",
			Usage: "开启断点续传",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "文件路径",
		},
		cli.StringFlag{
			Name:  "save, s",
			Usage: "OSS 中保存的目录",
		},
		cli.IntFlag{
			Name:  "random",
			Usage: "随机生成新的文件名长度，0 使用原文件名称",
		},
	}
}
