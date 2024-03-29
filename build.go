package skewer

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/secondarykey/skewer/config"
	"golang.org/x/xerrors"
)

func testCommand(a []string, pkgs []string) error {

	if len(pkgs) == 0 {
		return fmt.Errorf("test packages required.")
	}

	args := make([]string, 0, 1+len(pkgs)+len(a))
	args = append(args, "test")
	args = append(args, a...)
	args = append(args, pkgs...)

	//指定されたファイルをnameでビルド
	cmd := exec.Command("go", args...)

	err := runCommand("Test", cmd, true)
	if err != nil {
		return xerrors.Errorf("command run error: %w", err)
	}
	return nil
}

func runCommand(title string, cmd *exec.Cmd, test bool) error {

	err := setCommandPipe(cmd, test)
	if err != nil {
		return xerrors.Errorf("setCommandPipe() error: %w", err)
	}

	printVerbose(title, "Start.")
	printVerbose(cmd)

	err = cmd.Start()
	if err != nil {
		return xerrors.Errorf("Go %s command error: %w", title, err)
	}

	printVerbose(title, "wait...")

	err = cmd.Wait()
	if err != nil {
		return xerrors.Errorf("Go %s command wait error: %w", title, err)
	}

	printVerbose(title, "Complate.")

	return nil
}

func buildCommand(name string, files []string) error {

	if len(files) == 0 {
		return fmt.Errorf("build file required.")
	}

	args := make([]string, 0, 3+len(files))
	args = append(args, "build")
	args = append(args, "-o")
	args = append(args, name)
	args = append(args, files...)

	//指定されたファイルをnameでビルド
	cmd := exec.Command("go", args...)

	err := runCommand("Build", cmd, false)
	if err != nil {
		return xerrors.Errorf("command run error: %w", err)
	}

	return nil
}

func rebuildMonitor(s int, ch chan error) {

	conf := config.Get()
	bin := conf.Bin
	mode := conf.Mode
	d := time.Duration(s) * time.Second

	switch mode {
	case config.HTTPMode:
		fmt.Println("* HTTPMode")
		fmt.Println("MonitorPort:", conf.Port)
		fmt.Printf("Command    : go run %s %s\n", conf.Files, conf.Args)
	}

	// TODO TestMode

	for {
		status = getStatus()
		if status.canBuild() {

			time.Sleep(buildDuration())

			switch mode {
			case config.HTTPMode:
				cleanup(bin)
				go startServer(bin, conf.Port, conf.Files, conf.Args, ch)
			case config.TestMode:
				go startTest(conf.Files, ch)
			}
		} else if status == FatalStatus {
			return
		}
		time.Sleep(d)
	}
	return
}

func buildDuration() time.Duration {
	conf := config.Get()
	d := conf.Duration * 1000
	return time.Duration(d) * time.Millisecond
}
