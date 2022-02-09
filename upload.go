package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"time"
)

var UploadFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "image",
		Aliases:  []string{"i"},
		Usage:    "需要上传的镜像",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "name",
		Aliases:  []string{"n"},
		Usage:    "镜像名称",
		Required: true,
	},
	&cli.StringFlag{
		Name:        "path",
		Aliases:     []string{"p"},
		Usage:       "上传的相对路径，相对于私有仓库的根目录下的/docker-pbg-local",
		DefaultText: "/cems-frontend",
	},
	&cli.StringFlag{
		Name:        "version",
		Aliases:     []string{"v"},
		Usage:       "需要上传的镜像版本号",
		DefaultText: "1.0.0",
	},
}

func cmdUpload(ctx *cli.Context) error {
	//err := login()
	//if err != nil {
	//	return err
	//}

	server := ctx.String("server")
	if server == "" {
		server = config.Address
	}
	imagePath := ctx.String("path")
	if imagePath == "" {
		imagePath = config.ImagePath
	}
	image := ctx.String("image")
	version := ctx.String("version")
	if version == "" {
		version = config.Version
	}
	name := ctx.String("name")
	dateNow := time.Now().Format("20060102150405")
	tag := fmt.Sprintf("%s%s%s/%s:%s-%s", server, config.BasicPath, imagePath, name, version, dateNow)

	cmdString := fmt.Sprintf("docker tag %s:latest %s", image, tag)
	err := stdoutPrint(cmdString)
	if err != nil {
		fmt.Printf("tag失败！\n%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("tag成功！\n")

	cmdString = fmt.Sprintf("docker push %s", tag)
	err = stdoutPrint(cmdString)
	if err != nil {
		fmt.Printf("镜像上传失败！\n%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("docker上传镜像完成！%s镜像地址为%s\n", image, tag)

	//cmdString = fmt.Sprintf("docker logout %s", server)
	//err = stdoutPrint(cmdString)
	//if err != nil {
	//	fmt.Printf("退出失败！\n%s\n", err)
	//	os.Exit(1)
	//}
	//fmt.Printf("docker 退出登录！\n", image, tag)

	return nil
}
