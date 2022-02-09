package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

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
				tmp = strings.TrimSpace(tmp)
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
