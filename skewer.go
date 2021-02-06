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

	//TODO goが存在するかの確認
	terminal.Start()

	conf := config.Get()
	args := conf.Args
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

	go startServer(bin, conf.AppPort, args)

	return startProxyServer(conf.Port)
}

func checkConnection(port int) (chan bool, error) {

	ch := make(chan bool)

	//TODO 10秒待ったらエラー
	//TODO エラーが起こったら終了するイメージ

	go func() {
		var conn net.Conn
		var err error
		for range time.Tick(20 * time.Millisecond) {
			conn, err = net.Dial("tcp", fmt.Sprintf(":%d", port))
			if err == nil {
				log.Println("ok")
				break
			}
			log.Println("err")
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

	err = os.Remove(bin)
	if err != nil {
		log.Println(err)
	}
}

func startServer(bin string, port int, args []string) {

	setStatus(BuildStatus)
	//use goroutine
	err := build.Run(bin, args)
	if err != nil {
		setStatus(BuildErrorStatus)
		return
	}

	setStatus(StartupStatus)
	//コマンドを実行
	err = process.Run(bin)
	if err != nil {
		setStatus(StartupErrorStatus)
		return
	}

	go func() {
		ch, err := checkConnection(port)
		if err != nil {
		}

		v := <-ch
		if !v {
			setStatus(StartupErrorStatus)
		} else {
			setStatus(OKStatus)
		}
	}()
}
