package conf

import (
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {

	// log
	log.SetFormatter(&log.JSONFormatter{})
	// log.SetLevel(log.InfoLevel)
	// log.SetReportCaller(true)

	//rl, _ := rotatelogs.New("log/app.%Y%m%d%H%M.log")
	rl, _ := rotatelogs.New("log/app.%Y%m%d.log")
	mw := io.MultiWriter(os.Stdout, rl)
	log.SetOutput(mw)

	// config
	viper.AutomaticEnv()
	viper.AddConfigPath("conf") // path to look for the config file in
	viper.SetEnvPrefix("APP")
	viper.SetConfigType("yml")

	viper.SetConfigName("app")  // name of config file (without extension)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Error("Fatal error config file", err)
	}

	// watch
	viper.WatchConfig()

	// 配置文件读取, prod/dev -> app, 会优先读取prod/dev, 然后是app.yml, 优先环境变量，其次命令行参数
	var profile = viper.GetString("profile")
	log.Info("app profile = ", profile)
	if profile != "" {
		viper.SetConfigName("app." + profile)

		err = viper.MergeInConfig()
		if err != nil { // Handle errors reading the config file
			log.Error("Fatal error config file", err)
		}

		// watch
		viper.WatchConfig()
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", e.Name)
		log.Info("app name ", viper.GetString("name"))
		log.Info("snowflake.machineId = ", viper.GetInt("snowflake.machineId"))
	})

	// default value
	viper.SetDefault("db.maxIdleConn", 10)
	viper.SetDefault("db.maxOpenConn", 100)
}
