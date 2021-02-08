package skewer

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/secondarykey/skewer/config"

	"golang.org/x/xerrors"
)

func Listen(opts ...config.Option) error {

	//TODO Goが存在しない場合

	err := config.Set(opts)
	if err != nil {
		return xerrors.Errorf("config.Set() error: %w", err)
	}

	conf := config.Get()

	ch := make(chan error)
	done := make(chan error)

	startTerminal(conf.Verbose)

	go func() {
		for {
			select {
			case err := <-ch:
				if err != nil {
					msg := fmt.Sprintf("%+v", err)
					printVerbose(msg)
					if getStatus() == FatalStatus {
						done <- err
						return
					}
				}
			default:
			}
		}
	}()

	go func() {
		notifyMonitoring(conf.Args, conf.IgnoreFiles, ch)
	}()

	bin := conf.Bin
	//シグナル待受
	go func() {
		quit := make(chan os.Signal)
		// 受け取るシグナルを設定
		signal.Notify(quit, os.Interrupt)
		<-quit

		cleanup(bin)
		endTerminal()
	}()

	setStatus(WaitingForRebootStatus)
	go rebuildMonitor(5, ch)

	return <-done
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
			//TODO おかしい
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

	err := kill()
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

	log.Println("Build and Launch.")

	setStatus(BuildStatus)
	//use goroutine
	err := build(bin, args)
	if err != nil {
		log.Println("Build Error")
		setStatus(BuildErrorStatus)
		ch <- err
		return
	}

	setStatus(StartupStatus)
	//コマンドを実行
	err = run(bin)
	if err != nil {
		log.Println("Process Run Error")
		setStatus(StartupErrorStatus)
		ch <- err
		return
	}

	if port == 0 {
		log.Println("Complete Build and Launch")
		setStatus(OKStatus)
	} else {
		go func() {
			check, err := checkConnection(port)

			if err != nil {
				log.Println("HTTP Connection Error")
				setStatus(StartupErrorStatus)
				//TODO おかしい
				ch <- err
				return
			}

			v := <-check
			if v {
				log.Println("Complete Build and Launch")
				setStatus(OKStatus)
			} else {
				log.Println("HTTP Server Launch Error")
				setStatus(StartupErrorStatus)
			}
		}()
	}
}
