package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
)

type Config struct {
	Address  string
	Username string
	Password string
}

var config = Config{
	Address:  "server",
	Username: "username",
	Password: "password",
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
	},
}

func build() {

}

func upload() {

}

func reload() {

}
func processError(err error) {
	if err != nil {
		fmt.Printf("somthing was wrong: %s", err)
	}
}
func getCurrentPathByRunner() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
		fmt.Printf("当前工作路径：%s\n", abPath)
	}
	return abPath
}
func stdoutPrint(ctx *cli.Context, cmdString string) error {
	cmd := exec.Command("cmd", "/C", cmdString)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != nil {
					fmt.Printf("程序出现错误：%q", ctx.Err())
				} else {
					fmt.Printf("程序被终止！")
				}
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					return
				}
				fmt.Printf(readString)
			}
		}
	}(&wg)
	err = cmd.Start()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}
func cmdBuild(ctx *cli.Context) error {
	abPath := getCurrentPathByRunner()

	//fmt.Printf("准备登录私有仓库 %s\n", config.Address)
	cmdString := fmt.Sprintf("docker login %s -p %s -u %s", config.Address, config.Password, config.Username)
	cmd := exec.Command("cmd", cmdString)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("登录失败\n")
		return err
	}

	fmt.Println("登录成功!\n")
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
		cmdString = fmt.Sprintf("cd ../%s && yarn build", service)
		err = stdoutPrint(ctx, cmdString)
		//_, err = cmd.Output()
		if err != nil {
			fmt.Printf("前端打包失败！\n")
			return err
		}

		fmt.Printf("前端打包成功！\n")
		cmdString = fmt.Sprintf("cd %s", abPath)
		exec.Command("cmd", "/C", cmdString)
	}

	fmt.Printf("开始打包 docker 镜像\n")
	cmdString = fmt.Sprintf("docker-compose build %s", service)
	err = stdoutPrint(ctx, cmdString)
	if err != nil {
		return err
	}

	fmt.Printf("镜像打包成功\n")

	return nil
}

func cmdUpload(c *cli.Context) error {

	return nil
}

func run() int {
	app := cli.NewApp()
	app.Name = "afd"
	app.Usage = "暂时支持 docker-compose 打包"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "server",
			Usage: "私有仓库服务器地址",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "服务器用户名",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "服务器密码",
		},
	}
	app.Commands = commands

	return msg(app.Run(os.Args))
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
