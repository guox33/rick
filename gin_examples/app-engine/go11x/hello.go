package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/ping", Ping)
	go func() {
		_ = http.ListenAndServe("127.0.0.1:8080", nil)
	}()
	time.Sleep(time.Minute)

	r := gin.Default()
	r.Run()
}

func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Ping(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseForm()
	params := req.Form
	bs, _ := json.Marshal(params)
	_, _ = w.Write(bs)
	w.WriteHeader(http.StatusOK)
}
