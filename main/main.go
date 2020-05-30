package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"github.com/yiyedata/restapi"
	"github.com/yiyedata/restapi/main/dao/test"
	"github.com/yiyedata/restapi/utils"
	"runtime"
	"syscall"
	"time"
)

type Message struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	// AppEnv := os.Getenv("AppEnv")
	test.Init("root", "root", "tcp(:3306)/demo")

	app := restapi.New()
	app.RootGroup("", func(c *restapi.Context) {
		// c.String(200, []byte("ok"))

		// log.Printf("1")
		c.Next()
		// log.Printf("2")
	})
	group := app.Group("/test", func(c *restapi.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	})
	{
		group.GET("/:id", func(c *restapi.Context) {
			var msg = &Message{Code: 0, Message: "ok"}
			var t, err = test.SelectByID(utils.ParseInt64(c.GetPara("id")))
			if err != nil {
				msg.Code = 1
				msg.Message = err.Error()
			} else {
				msg.Data = t
			}
			c.JSON(msg)
		})
		group.GET("/index/:key", func(c *restapi.Context) {
			var msg = &Message{Code: 0, Message: "ok"}
			var t, err = test.SelectByKey(c.GetPara("key"))
			if err != nil {
				msg.Code = 1
				msg.Message = err.Error()
			} else {
				msg.Data = t
			}
			c.JSON(msg)
		})
		group.POST("/update", func(c *restapi.Context) {
			var msg = &Message{Code: 0, Message: "ok"}
			json := c.Body()
			id := json.GetString("id")
			key := json.GetString("key")
			value := json.GetString("value")
			if value != "" {
				_, err := test.UpdateValue(utils.ParseInt64(id), value)
				if err != nil {
					msg.Code = 1
					msg.Message = err.Error()
				}
			} else {
				_, err := test.UpdateKey(utils.ParseInt64(id), key)
				if err != nil {
					msg.Code = 1
					msg.Message = err.Error()
				}
			}
			c.JSON(msg)
		})
		group.POST("/insert", func(c *restapi.Context) {
			json := c.Body()
			key := json.GetString("key")
			value := json.GetString("value")
			obj, err := test.Insert(key, value)
			var msg = &Message{Code: 0, Message: "ok"}
			if err != nil {
				msg.Code = 1
				msg.Message = err.Error()
			} else {
				msg.Data = obj
			}
			c.JSON(msg)
		})
	}
	app.GET("/health", func(c *restapi.Context) {
		var msg = &Message{Code: 0, Message: "health"}
		c.JSON(msg)
	})
	go func() {
		err := app.Run(":8087")
		fmt.Println(err)
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
}
