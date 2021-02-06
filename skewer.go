package skewer

import (
	"io"
	"log"
	"net/http"
)

func Listen() error {

	//goが存在するかの確認

	//use goroutine
	b := build.Run("skewer-bin")

	go func() {
		for {
			select {
			case <-b:
			}
		}
	}()

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
