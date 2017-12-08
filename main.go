package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums_web/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Usage(program string) {
	fmt.Printf("\nusage: %s conf/cf.xml\n", program)
	fmt.Printf("\nconf/cf.xml      configure file\n")
}

func main() {
	if len(os.Args) != 2 {
		Usage(os.Args[0])
		os.Exit(-1)
	}
	log.Println("[Main] Starting program")
	defer log.Println("[Main] Exit program successful.")

	conf := Configure(os.Args[1])

	server := fmt.Sprintf(":%v", conf.GinServer.Port)
	//url := fmt.Sprintf("%v", conf.GinServer.Url)
	//配置gin
	gin.SetMode(conf.GinServer.Mode)
	r := gin.New()
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("html/*")

	//制定handler
	//登陆页
	r.GET(conf.GinServer.Url, homeHandler)
	r.GET(conf.GinServer.Url+"/login", loginHandler)
	r.POST(conf.GinServer.Url+"/signup", signupHandler)
	r.POST(conf.GinServer.Url+"/signin", signinHandler)

	//登陆成功后页面
	r.POST(conf.GinServer.Url+"/content", contentHandler)

	//图片内容检测页
	r.GET(conf.GinServer.Url+"/objectdection", objectdectionHandler)
	r.POST(conf.GinServer.Url+"/od", odHandler)
	r.GET(conf.GinServer.Url+"/deeplearning", deeplearningHandler)
	r.POST(conf.GinServer.Url+"/dp", dpHandler)

	//上传页
	r.GET(conf.GinServer.Url+"/upload", uploadHandler)
	r.POST(conf.GinServer.Url+"/ul", ulHandler)

	//搜索页
	r.GET(conf.GinServer.Url+"/search", searchHandler)
	r.POST(conf.GinServer.Url+"/s", sHandler)

	//相册列表
	r.GET(conf.GinServer.Url+"/albumlist", albumlistHandler)
	r.POST(conf.GinServer.Url+"/alist", alistHandler)

	//相册里图片
	r.POST(conf.GinServer.Url+"/download", downloadHandler)

	//测试功能页面
	r.POST(conf.GinServer.Url+"/index", indexHandler)
	r.GET(conf.GinServer.Url+"/index2", index2Handler)
	r.GET(conf.GinServer.Url+"/index2new", index2newHandler)
	r.GET(conf.GinServer.Url+"/b", bHandler)
	r.GET(conf.GinServer.Url+"/b2", b2Handler)
	r.GET(conf.GinServer.Url+"/c", cHandler)

	//起一个http服务器
	s := &http.Server{
		Addr:         server,
		Handler:      r,
		ReadTimeout:  time.Duration(conf.GinServer.Timeout_read_s) * time.Second,
		WriteTimeout: time.Duration(conf.GinServer.Timeout_write_s) * time.Second,
	}
	go func(s *http.Server) {
		log.Printf("[Main] http server start\n")
		err := s.ListenAndServe()
		log.Printf("[Main] http server stop (%+v)\n", err)
	}(s)
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for {
		select {
		case sig := <-signals:
			log.Println("[Main] Catch signal", sig)
			//平滑关闭server
			err := s.Shutdown(context.Background())
			log.Printf("[Main] start gracefully shuts down http serve %+v", err)
			return
		}
	}
}
