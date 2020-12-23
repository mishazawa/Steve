package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	local      string
	remote     string
	uploadOnly bool
)

func init() {
	flag.StringVar(&local, "local", ".", "Local folder")
	flag.StringVar(&remote, "remote", "", "Remote host")
	flag.BoolVar(&uploadOnly, "upload", false, "Only upload")
	flag.Parse()

	if len(remote) == 0 {
		panic("No -remote flag")
	}
}

func main() {
	if !uploadOnly {
		fmt.Println("Sync [pull]")
		if err := Download(local, remote); err != nil {
			fmt.Println("Can't sync [pull]", err.Error())
			return
		}

		fmt.Println("Running server")

		current, _ := filepath.Abs(".")

		if err := os.Chdir(local); err != nil {
			fmt.Println("Can't change directory", current, err.Error())
			return
		}

		if err := executeServer([]string{"java", "-jar", "server.jar"}); err != nil {
			fmt.Println("Can't run server", err.Error())
			return
		}

		if err := os.Chdir(current); err != nil {
			fmt.Println("Can't change directory", current, err.Error())
			return
		}
	}

	fmt.Println("Sync [push]")
	if err := Upload(local, fmt.Sprintf("%s/api/upload", remote)); err != nil {
		fmt.Println("Can't sync [push]", err.Error())
		return
	}

	fmt.Println("Done.")
}

func executeServer(executable []string) error {
	if _, err := exec.LookPath(executable[0]); err != nil {
		return err
	}

	cmd := exec.Command(executable[0], executable[1:]...)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
