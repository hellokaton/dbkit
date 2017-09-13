package main

import (
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"github.com/urfave/cli"
	"github.com/go-yaml/yaml"
	"github.com/robfig/cron"
)

var (
	// config file root path
	RootPath     = "."
	TarCmd       = "tar"
	MysqlDumpCmd = "mysqldump"
)

// dbkit config
type DbKitConfig struct {
	Db struct {
		Type     string
		Host     string
		Port     int
		Username string
		Password string
		Database string
		Options  []string
	}

	CronExpression string

	Storage struct {
		Type     string
		Account  string
		Location string
	}
}

type ExportResult struct {
	// Path to exported file
	Path string
	// MIME type of the exported file (e.g. application/x-tar)
	MIME string
	// Any error that occured during `Export()`
	Error *Error
}

func main() {

	app := cli.NewApp()
	app.Name = "dbkit"
	app.Usage = "A simple database tools"
	app.Author = "https://github.com/biezhi"
	app.Email = "biezhi.me@gmail.com"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: ".",
			Usage: "config file path",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		RootPath = ctx.String("config")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "start dbkit tools",
			Action: Start,
		},
		{
			Name:  "stop",
			Usage: "stop dbkit tools",
			Action: func(ctx *cli.Context) {
				pid := findPid()
				if pid > 0 {
					proc, err := os.FindProcess(pid)
					if err != nil {
						log.Fatalln(err)
					}
					log.Printf("kill pid: %d\n", pid)
					proc.Kill()
				}
			},
		},
		{
			Name:  "status",
			Usage: "",
			Action: func() {
				pid := findPid()
				if pid > 0 {
					log.Println("dbkit is running.")
				} else {
					log.Println("dbkit is stop.")
				}
			},
		},
	}
	app.Run(os.Args)
	os.Exit(0)
}

// dbkit service start
func Start(_ *cli.Context) {
	conf := ParseConfig()

	if conf.CronExpression == "" {
		log.Fatalln("Cron表达式不能为空")
		return
	}

	c := cron.New()
	spec := conf.CronExpression
	c.AddFunc(spec, func() {
		log.Println("开始备份数据库...")

		var exp *ExportResult
		if conf.Db.Type == "mysql" {
			exp = ExportMySql(*conf)
		}

		if exp.Error != nil {
			log.Printf("遇到error")
			log.Printf(exp.Error.CmdOutput)
			log.Fatal(exp.Error.err)
		} else {
			log.Printf("备份成功到: %s\n", exp.Path)
		}

	})
	c.Start()
	select {}
}

// parse config.yml
func ParseConfig() *DbKitConfig {
	var config *DbKitConfig
	configPath := filepath.Join(RootPath, "config.yml")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err.Error())
	}
	return config
}
