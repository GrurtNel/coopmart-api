package init

import (
	"feedback/x/config"
	"feedback/x/db/mongodb"
	"flag"
	"io"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func init() {
	loadConfig()
	initLog()
	initDB()
}

var context *config.Context
var DeviceConfig string

type FeedbackConfig struct {
	HighRate   float32
	CreditRate float32
	MediumRate float32
}

var FConfig FeedbackConfig

func loadConfig() {
	context, _ = config.LoadContext("app.conf", []string{""})
}

func initLog() {
	//config for gin request log
	{
		f, _ := os.Create(path.Join("log", "gin.log"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
	//config for app log use glog
	{
		logDir, _ := context.String("glog.log_dir")
		logStd, _ := context.String("glog.alsologtostderr")
		flag.Set("alsologtostderr", logStd)
		flag.Set("log_dir", logDir)
		flag.Parse()
	}
}

func initDB() {
	DeviceConfig, _ = context.String("device.host")
	// Read configuration.
	mongodb.MaxPool = context.IntDefault("mongo.maxPool", 0)
	mongodb.PATH, _ = context.String("mongo.path")
	mongodb.DBNAME, _ = context.String("mongo.database")
	mongodb.UNAME, _ = context.String("mongo.uname")
	mongodb.PASSWORD, _ = context.String("mongo.password")
	mongodb.CheckAndInitServiceConnection()
}
