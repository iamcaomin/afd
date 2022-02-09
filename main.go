package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

type Config struct {
	Address   string
	BasicPath string
	ImagePath string
	Version   string
}

var config = Config{
	Address:   "******",
	BasicPath: "/docker",
	ImagePath: "/frontend",
	Version:   "1.0.0",
}

var commands = []*cli.Command{
	{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "镜像打包",
		Action:  cmdBuild,
		Flags:   BuildFlags,
	},
	{
		Name:    "upload",
		Aliases: []string{"u"},
		Usage:   "上传镜像",
		Action:  cmdUpload,
		Flags:   UploadFlags,
	},
}

//func login() error {
//	fmt.Printf("准备登录私有仓库 %s\n", config.Address)
//	cmdString := fmt.Sprintf("docker login %s", config.Address)
//	err := stdoutPrint(cmdString)
//	if err != nil {
//		fmt.Printf("登录失败！\n")
//		os.Exit(1)
//	}
//
//	fmt.Printf("登录成功!\n")
//	return nil
//}

func run() int {
	app := cli.NewApp()
	app.Name = "afd"
	app.Usage = "暂时支持 docker-compose 打包，建议首次登录docker私有仓库的用户，先命令行手动登录成功一次！"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "server",
			Usage: "私有仓库服务器地址",
		},
	}
	app.Commands = commands

	return msg(app.Run(os.Args))
}

func msg(err error) int {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
