package cmd

import (
	"west.garden/template/common/aws_s3"
	"west.garden/template/common/config"
	"west.garden/template/common/log"
	"west.garden/template/common/redis"
	"west.garden/template/common/wrapper"
	"west.garden/template/model"
	pkgRedis "west.garden/template/common/redis"
	"west.garden/template/router"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	svc "github.com/judwhite/go-svc"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

type Application struct {
	wrapper    wrapper.Wrapper
	ginEngine  *gin.Engine
	httpServer *http.Server
	cron       *cron.Cron
}

var cfgFile *string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the api",
	Long: `usage example:
	server(.exe) start -c config.json
	start the api`,
	Run: func(cmd *cobra.Command, args []string) {
		app := &Application{}
		if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	cfgFile = startCmd.Flags().StringP("config", "c", "", "api config file (required)")
	startCmd.MarkFlagRequired("config")
}

func (app *Application) Init(_ svc.Environment) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	log.Init(&cfg.Logger)

	if err := model.Init(&cfg.Mysql); err != nil {
		log.Log.Error("Init Mysql Err:", err.Error())
		return err
	}
	
	if err = redis.Init(&cfg.Redis); err != nil {
		log.Log.Error("Init Redis Err:", err.Error())
	}
	if err = pkgRedis.Init(&cfg.Redis); err != nil {
		log.Log.Error("Init pkg Redis Err:", err.Error())
	}

	if err = aws_s3.Init(&cfg.S3); err != nil {
		log.Log.Error("Init S3 Err:", err.Error())
		return err
	}
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	app.ginEngine = router.InitApiRouter(&cfg)

	buf, _ := json.Marshal(cfg)
	fmt.Println("config:", string(buf))

	return nil
}

func (app *Application) Start() error {
	cfg := config.GetConfig().Http
	app.wrapper.Wrap(func() {
		app.httpServer = &http.Server{
			Handler:        app.ginEngine,
			Addr:           cfg.ListenAddr,
			ReadTimeout:    cfg.ReadTimeout * time.Second,
			WriteTimeout:   cfg.WriteTimeout * time.Second,
			IdleTimeout:    cfg.IdleTimeout * time.Second,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
		}
		if err := app.httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	})
	log.Log.Info("Api Server Started,Listen on:", cfg.ListenAddr)
	return nil
}

func (app *Application) Stop() error {
	if app.httpServer != nil {
		if err := app.httpServer.Shutdown(context.Background()); err != nil {
			fmt.Printf("Api Server shutdown error:%v\n", err)
		}
		fmt.Println("Api Server shutdown")
	}

	model.Close()
	app.wrapper.Wait()
	fmt.Println("Shutdown end")
	return nil
}
