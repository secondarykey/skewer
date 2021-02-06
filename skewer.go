package skewer

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Listen() error {

	//goが存在するかの確認
	go func() {
		quit := make(chan os.Signal)
		// 受け取るシグナルを設定
		signal.Notify(quit, os.Interrupt)
		<-quit

		log.Println("処理中...")
		time.Sleep(2 * time.Second)
		log.Println("オワタよ")
		os.Exit(0)
	}()

	//use goroutine
	b := build.Run("skewer-bin")

	//シグナル待受
	http.HandleFunc("/", proxyHandler)

	proxy := ":3000"
	log.Println(proxy, "Start")
	return http.ListenAndServe(proxy, nil)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {

	//TODO 現状のビルド状況を確認

	//URLを作成
	url := r.URL.String()
	req := "http://localhost:8080" + url

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
