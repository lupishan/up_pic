package main

import (
	"bytes"
	"fmt"
	"github.com/upyun/go-sdk/upyun"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/up", up)

	log.Println("Starting v1 server ...")
	log.Fatal(http.ListenAndServe(":1210", nil))
}

func hello(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("template/go.html")
	if (err != nil) {
		log.Println(err)
	}
	t.Execute(w, nil)
}

// 上传图片路由
func up(w http.ResponseWriter, r *http.Request){
	up := upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   "static-css",
		Operator: "Operator",
		Password: "Password",
	})

	headerByte := make([]byte, 8)

	b, _ :=  ioutil.ReadAll(r.Body)
	reader1 := bytes.NewReader(b)
	reader2 := bytes.NewReader(b)
	reader2.Read(headerByte)

	// 判断是什么类型的图片
	xStr := fmt.Sprintf("%x", headerByte)
	ext := ""
	switch {
	case xStr == "89504e470d0a1a0a":
		ext = ".png"
	case xStr[:12] == "474946383961" || xStr[:12] == "474946383761":
		ext = ".gif"
	case xStr[:4] == "ffd8":
		ext = ".jpg"
	}

	// 生成文件名（不加后缀）
	h := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(9999))

	picName := h + ext

	// 上传文件
	fmt.Println(up.Put(&upyun.PutObjectConfig{
		Path:		picName,
		Reader:		reader1,
	}))

	w.Write([]byte(picName))
}