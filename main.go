package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ludyyy-lu/goBlogService/global"
	"github.com/ludyyy-lu/goBlogService/internal/model"
	"github.com/ludyyy-lu/goBlogService/internal/routers"
	"github.com/ludyyy-lu/goBlogService/pkg/logger"
	"github.com/ludyyy-lu/goBlogService/pkg/setting"
	"gopkg.in/natefinch/lumberjack.v2"
)

// go语言的程序执行顺序：全局变量初始化 -> init方法 -> main方法
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	//实际上并没有error
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	err = setupDBEngin()
	if err != nil {
		log.Fatalf("init.setupDBEngin err: %v", err)
	}
}
func main() {
	global.Logger.Infof("%s: goHttpWeb-practice/%s", "ludy-lu", "blog-service")

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	//设置已经映射好的配置和gin的运行模式
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort, //8000
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	return nil
}

func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}
func setupDBEngin() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
