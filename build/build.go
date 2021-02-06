package build

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/secondarykey/skewer/terminal"

	"golang.org/x/xerrors"
)

func Check() error {
	path, err := exec.LookPath("go")
	if err != nil {
		return xerrors.Errorf("go is not avaliable", err)
	}
	log.Println("go is available at", path)
	return nil
}

func Run(name string, files []string) error {

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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	terminal.Verbose("build start")
	err := cmd.Start()
	if err != nil {
		return xerrors.Errorf("go build error: %w", err)
	}

	terminal.Verbose("build wait...")
	err = cmd.Wait()
	if err != nil {
		return xerrors.Errorf("command wait error: %w", err)
	}

	terminal.Verbose("build complate")

	return nil
}
