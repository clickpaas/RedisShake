package main

import (
	"RedisShake/cmd/commands"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help") {
		commands.ConvertArgs2Toml(true, false)
		os.Exit(0)
	}

	commandLine := strings.Join(os.Args[1:], ",")
	if strings.Contains(commandLine, "reader") || strings.Contains(commandLine, "writer") ||
		strings.Contains(commandLine, "filter") || strings.Contains(commandLine, "advanced") ||
		strings.Contains(commandLine, "module") {
		tomlPath, err := commands.ConvertArgs2Toml(false, strings.Contains(commandLine, "dry_run"))
		if err != nil || tomlPath == "" {
			os.Exit(0)
		}
		// os.Args = []string{os.Args[0], tomlPath}
		goos := runtime.GOOS
		if goos == "windows" {
			err = executeWithCopyOutput("./redis-shake.exe", tomlPath)
		} else {
			err = executeWithCopyOutput("./redis-shake", tomlPath)
		}
		if err != nil {
			log.Fatal(err)
		}
	}

}

func executeWithCopyOutput(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	// 获取输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return err
	}

	// 使用 goroutine 实时复制输出
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	// 等待命令完成
	return cmd.Wait()
}
