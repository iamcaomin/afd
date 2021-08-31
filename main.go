package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	Address string
}

var config = Config{
	Address: "server address",
}

var commands = []*cli.Command{
	{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "镜像打包",
		Action:  cmdBuild,
		Flags: []cli.Flag{
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
		},
	},
	{
		Name:    "upload",
		Aliases: []string{"u"},
		Usage:   "上传镜像",
		Action:  cmdUpload,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "上传的路径，相对路径，相对于私有仓库的根目录下的/docker-pbg-local",
				Required: true,
			},
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
				Name:     "version",
				Aliases:  []string{"v"},
				Usage:    "需要上传的镜像版本号",
				Required: true,
			},
		},
	},
}

func login() error {
	fmt.Printf("准备登录私有仓库 %s\n", config.Address)
	cmdString := fmt.Sprintf("docker login %s", config.Address)
	err := stdoutPrint(cmdString)
	if err != nil {
		fmt.Printf("登录失败！\n")
		os.Exit(1)
	}

	fmt.Printf("登录成功!\n")
	return nil
}

func printLog(reader io.ReadCloser) error {
	bucket := make([]byte, 1024)
	buffer := make([]byte, 100)
	for {
		num, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				err = nil
			}
			return err
		}

		if num > 0 {
			line := ""
			bucket = append(bucket, buffer[:num]...)
			tmp := string(bucket)
			if strings.Contains(tmp, "\n") {
				ts := strings.Split(tmp, "\n")
				if len(ts) > 1 {
					line = strings.TrimSpace(strings.Join(ts[:len(ts)-1], "\n"))
					bucket = []byte(ts[len(ts)-1])
				} else {
					line = ts[0]
					bucket = bucket[:0]
				}
				fmt.Printf("%s\n", line)
			}
		}
	}
}

func stdoutPrint(cmdString string) error {
	//fmt.Printf("cmdString: %s\n", cmdString)
	fmt.Print("---------- 分割线 ----------\n")
	cmd := exec.Command("cmd", "/C", cmdString)
	closed := make(chan struct{})
	defer close(closed)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error Starting command: %s.....\n", err.Error())
		return err
	}

	go printLog(stdout)
	go printLog(stderr)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error Waiting command: %s.....\n", err.Error())
		return err
	}
	return nil
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
	cmdString = fmt.Sprintf("docker-compose build %s", service)

	err = stdoutPrint(cmdString)
	if err != nil {
		fmt.Printf("镜像打包失败！\n%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("镜像打包成功\n")

	return nil
}

func cmdUpload(ctx *cli.Context) error {
	err := login()
	if err != nil {
		return err
	}

	server := ctx.String("server")
	if server == "" {
		server = config.Address
	}
	imagePath := ctx.String("path")
	image := ctx.String("image")
	version := ctx.String("version")
	name := ctx.String("name")
	dateNow := time.Now().Format("20060102150405")
	tag := fmt.Sprintf("%s/docker-pbg-local%s/%s:%s-%s", server, imagePath, name, version, dateNow)

	cmdString := fmt.Sprintf("docker tag %s:latest %s", image, tag)
	err = stdoutPrint(cmdString)
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

	return nil
}

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
