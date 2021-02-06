package skewer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

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

	terminal.Start()

	conf := config.Get()
	bin := conf.Bin

	//TODO goが存在するかの確認

	//シグナル待受
	go func() {
		quit := make(chan os.Signal)
		// 受け取るシグナルを設定
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Println("処理中...")
		process.Kill()

		os.Remove(bin)
		terminal.End()
	}()

	//use goroutine
	err = build.Run(bin, conf.Args)
	if err != nil {
		return xerrors.Errorf("build.Run() error: %w", err)
	}

	//コマンドを実行
	err = process.Run(bin)
	if err != nil {
		return xerrors.Errorf("process.Run() error: %w", err)
	}

	http.HandleFunc("/", proxyHandler)

	proxy := fmt.Sprintf(":%d", conf.Port)
	log.Println(proxy, "Start")
	return http.ListenAndServe(proxy, nil)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {

	//TODO 現状のビルド状況を確認

	//URLを作成
	conf := config.Get()
	url := r.URL.String()

	//TODO クエリは？
	req := fmt.Sprintf("%s://%s:%d%s", conf.Schema, conf.Server, conf.Port, url)

	//相手にリクエスト
	resp, err := http.Get(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	//書き込み
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func cleanup() {
}
