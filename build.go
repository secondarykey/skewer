package skewer

import (
	"fmt"
	"log"
	"os/exec"

	"golang.org/x/xerrors"
)

func checkGo() error {
	path, err := exec.LookPath("go")
	if err != nil {
		return xerrors.Errorf("go is not avaliable", err)
	}
	log.Println("go is available at", path)
	return nil
}

func build(name string, files []string) error {

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

	err := setCommandPipe(cmd)
	if err != nil {
		return xerrors.Errorf("setCommandPipe() error: %w", err)
	}

	printVerbose("Build Start.")
	printVerbose(cmd)

	err = cmd.Start()
	if err != nil {
		return xerrors.Errorf("go build error: %w", err)
	}

	printVerbose("Build wait...")

	err = cmd.Wait()
	if err != nil {
		return xerrors.Errorf("command wait error: %w", err)
	}

	printVerbose("Build Complate.")

	return nil
}
