package skewer

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/secondarykey/skewer/config"
)

func startProxyServer(port int) error {

	http.HandleFunc("/", proxyHandler)

	proxy := fmt.Sprintf(":%d", port)
	log.Println(proxy, "Start")

	return http.ListenAndServe(proxy, nil)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {

	status := getStatus()
	//TODO 現状のビルド状況を確認
	if status == WaitingForRebootStatus {
		conf := config.Get()
		bin := conf.Bin
		cleanup(bin)
		go startServer(bin, conf.AppPort, conf.Args)
		return
	} else if status != OKStatus {
		w.Write([]byte("<h1>skewer error</h1>"))
		return
	}

	//URLを作成
	conf := config.Get()
	url := r.URL.String()

	//TODO クエリは？
	req := fmt.Sprintf("%s://%s:%d%s", conf.Schema, conf.Server, conf.AppPort, url)

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
