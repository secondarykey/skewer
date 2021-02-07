package skewer

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/secondarykey/skewer/build"
	"github.com/secondarykey/skewer/config"
	"github.com/secondarykey/skewer/process"
	"github.com/secondarykey/skewer/terminal"
	"golang.org/x/xerrors"
)

func Listen(opts ...config.Option) error {

	err := config.Set(opts)
	if err != nil {
		return xerrors.Errorf("config.Set() error: %w", err)
	}

	conf := config.Get()

	var ch chan error

	terminal.Start(conf.Verbose)

	go func() {
		err = monitoring(conf.Args, conf.IgnoreFiles)
		ch <- err
		if err != nil {
			log.Printf("monitoring error: %+v\n", err)
		}
	}()

	//TODO goが存在するかの確認
	bin := conf.Bin

	//シグナル待受
	go func() {
		quit := make(chan os.Signal)
		// 受け取るシグナルを設定
		signal.Notify(quit, os.Interrupt)
		<-quit

		cleanup(bin)
		terminal.End()
	}()

	setStatus(WaitingForRebootStatus)
	go rebuildMonitor(5, ch)

	//TODO ch を監視
	go func() {
		err := <-ch
		if err != nil {
		}
	}()

	return startProxyServer(conf.Port)
}

func checkConnection(port int) (chan bool, error) {

	ch := make(chan bool)

	//TODO 10秒待ったらエラー
	//TODO エラーが起こったら終了するイメージ

	go func() {
		var conn net.Conn
		var err error
		for range time.Tick(100 * time.Millisecond) {
			conn, err = net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				break
			}
			log.Println(err)
		}
		defer conn.Close()
		ch <- true
		close(ch)
	}()

	return ch, nil
}

func cleanup(bin string) {

	err := process.Kill()
	if err != nil {
		log.Println(err)
	}

	//TODO 付け焼き刃
	time.Sleep(1 * time.Second)

	if _, err := os.Stat(bin); err == nil {
		err = os.Remove(bin)
		if err != nil {
			log.Println(err)
		}
	}
}

func startServer(bin string, port int, args []string, ch chan error) {

	setStatus(BuildStatus)
	//use goroutine
	err := build.Run(bin, args)
	if err != nil {
		log.Printf("build error: ============\n%+v", err)
		setStatus(BuildErrorStatus)
		ch <- err
		return
	}

	setStatus(StartupStatus)
	//コマンドを実行
	err = process.Run(bin, ch)
	if err != nil {
		log.Printf("process error: ============\n%+v", err)
		setStatus(StartupErrorStatus)
		ch <- err
		return
	}

	go func() {
		check, err := checkConnection(port)
		if err != nil {
			log.Printf("connection error: ============\n%+v", err)
			setStatus(StartupErrorStatus)
			return
		}

		v := <-check
		if v {
			log.Println("Complete Build and Launch")
			setStatus(OKStatus)
		} else {
		}
	}()
}
