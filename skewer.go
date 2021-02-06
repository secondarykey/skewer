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

	go func() {
		err = monitoring(conf.Args)
		if err != nil {
			log.Printf("monitoring error: %+v\n", err)
		}
	}()

	//TODO goが存在するかの確認
	terminal.Start()

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
	go rebuildMonitor(10)

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

	//TODO 付け焼き刃
	time.Sleep(1 * time.Second)

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
		//TODO build error
		log.Printf("build error: ============\n%+v", err)
		setStatus(BuildErrorStatus)
		return
	}

	setStatus(StartupStatus)
	//コマンドを実行
	err = process.Run(bin)
	if err != nil {
		log.Printf("process error: ============\n%+v", err)
		setStatus(StartupErrorStatus)
		return
	}

	//TODO エラー時の挙動
	go func() {
		ch, err := checkConnection(port)
		if err != nil {
			log.Printf("connection error: ============\n%+v", err)
			setStatus(StartupErrorStatus)
			return
		}

		v := <-ch
		if !v {
			setStatus(StartupErrorStatus)
		} else {
			setStatus(OKStatus)
		}
	}()
}
