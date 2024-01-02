package main

import (
	"context"
	_ "github.com/apache/skywalking-go"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	zlog "github.com/vearne/zaplog"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	//"io"
	"net/http"
)

var rdb *redis.Client

func main() {
	zlog.InitLogger("/tmp/aa.log", "debug")

	// 添加Prometheus的相关监控
	// /metrics
	go func() {
		r := gin.Default()
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		r.Run(":9090")
	}()

	prometheus.NewRegistry()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "xxeQl*@nFE", // 密码
		DB:       0,            // 数据库
		PoolSize: 20,           // 连接池大小
	})

	_, err := rdb.Set(context.Background(), "abc", "def", 0).Result()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		val, err := rdb.Incr(c, "helloCounter").Result()
		zlog.Info("test hello", zap.Int64("val", val), zap.Error(err))
		setRes, err := rdb.Set(c, "abc", "def", 0).Result()
		zlog.Info("test hello", zap.String("setRes", setRes), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})
	r.GET("/hello2", func(c *gin.Context) {
		val, err := rdb.Incr(c, "helloCounter").Result()
		zlog.Info("test hello", zap.Int64("val", val), zap.Error(err))
		//req, err := http.NewRequest("GET", "http://localhost:18001/sayHello", nil)
		//resp, err := http.DefaultClient.Do(req)
		//dealErr(err)

		client := resty.New()
		resp, err := client.R().
			Get("http://localhost:18001/sayHello")

		c.JSON(http.StatusOK, gin.H{
			"message": resp.String(),
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		g, _ := errgroup.WithContext(context.Background())

		g.Go(func() error {
			val, err := rdb.Incr(context.Background(), "helloCounter2").Result()
			zlog.Info("ping", zap.Int64("val", val), zap.Error(err))
			return nil
		})
		g.Go(func() error {
			hsetRes, err := rdb.HSet(context.Background(), "xyz", "def", 0).Result()
			zlog.Info("ping", zap.Int64("setRes", hsetRes), zap.Error(err))
			return nil
		})
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8000")
}

func dealErr(err error) {
	if err != nil {
		panic(err)
	}
}
