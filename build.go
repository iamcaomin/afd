package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
)

var BuildFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "service",
		Aliases:  []string{"s"},
		Usage:    "需要打包的服务",
		Required: true,
	},
	&cli.BoolFlag{
		Name:    "dist",
		Aliases: []string{"d"},
		Usage:   "是否需要前端打包",
	},
}

func cmdBuild(ctx *cli.Context) error {
	var cmdString string
	abPath, _ := os.Getwd()
	file, err := os.Open("./docker-compose.yaml")
	if err != nil {
		fmt.Printf("当前目录没有 docker-compose.yaml 文件！\n")
		return err
	}
	defer file.Close()

	var service = ctx.String("service")
	fmt.Printf("准备开始打包服务 %s 镜像\n", service)

	var needDist = ctx.Bool("dist")
	if needDist {
		fmt.Printf("开始打包前端项目\n")
		cmdString = fmt.Sprintf("cd %s && yarn build", service)

		err = stdoutPrint(cmdString)
		if err != nil {
			fmt.Printf("前端打包失败！\n%s\n", err)
			os.Exit(1)
		}

		fmt.Printf("前端打包成功！\n")
		cmdString = fmt.Sprintf("cd %s", abPath)
		exec.Command("cmd", "/C", cmdString)
	}

	fmt.Printf("开始打包 docker 镜像\n")
	cmdString = fmt.Sprintf("docker-compose build --no-cache %s", service)

	err = stdoutPrint(cmdString)
	if err != nil {
		fmt.Printf("镜像打包失败！\n%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("镜像打包成功\n")

	return nil
}
