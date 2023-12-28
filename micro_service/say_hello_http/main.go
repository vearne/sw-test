package main

import (
	"context"
	"fmt"
	_ "github.com/apache/skywalking-go"
	"github.com/redis/go-redis/v9"
	zlog "github.com/vearne/zaplog"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	zlog.InitLogger("/tmp/sayHello.log", "debug")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "xxeQl*@nFE", // 密码
		DB:       0,            // 数据库
		PoolSize: 20,           // 连接池大小
	})

	http.HandleFunc("/sayHello", func(w http.ResponseWriter, r *http.Request) {
		// print Headers
		for key, val := range r.Header {
			fmt.Printf("%v:%v\n", key, val)
		}
		val, err := rdb.Incr(context.Background(), "sayHello").Result()
		zlog.Info("test hello", zap.Int64("val", val), zap.Error(err))

		fmt.Fprintf(w, "Hello, sw-go")
	})

	log.Println("say_hello_http starting...")
	log.Fatal(http.ListenAndServe(":18001", nil))
}
