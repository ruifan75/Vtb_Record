package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config *MainConfig
var Changed bool

type UsersConfig struct {
	TargetId     string
	Name         string
	DownloadDir  string
	NeedDownload bool
	TransBiliId  string
	ExtraConfig  map[string]interface{}
}
type ModuleConfig struct {
	//EnableProxy     bool
	//Proxy           string
	Name             string
	Enable           bool
	Users            []UsersConfig
	DownloadProvider string
	ExtraConfig      map[string]interface{}
}
type MainConfig struct {
	CriticalCheckSec int
	NormalCheckSec   int
	LogFile          string
	LogFileSize      int
	LogLevel         string
	RLogLevel        string
	DownloadQuality  string
	DownloadDir      string
	UploadDir        string
	Module           []ModuleConfig
	ExpressPort      string
	EnableTS2MP4     bool
	YtdlpCookies     string
	ExtraConfig      map[string]interface{}
}

var v *viper.Viper

func InitConfig() {
	log.Print("Init config!")
	initConfig()
	log.Print("Load config!")
	_, _ = ReloadConfig()
	//fmt.Println(Config)
}

func initConfig() {
	/*viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")
	viper.SetConfigType("json")*/
	v = viper.NewWithOptions(viper.KeyDelimiter("::::"))
	v.SetConfigFile(viper.ConfigFileUsed())
	v.WatchConfig()
	err := v.ReadInConfig()
	if err != nil {
		fmt.Printf("config file error: %s\n", err)
		os.Exit(1)
	}

	Changed = true
	v.OnConfigChange(func(in fsnotify.Event) {
		Changed = true
	})
}

func ReloadConfig() (bool, error) {
	if !Changed {
		return false, nil
	}
	Changed = false
	err := v.ReadInConfig()
	if err != nil {
		return true, err
	}
	config := &MainConfig{}
	mdMap := make(map[string]*mapstructure.Metadata, 10)
	mdMap[""] = &mapstructure.Metadata{}
	err = v.Unmarshal(config, func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			func(inType reflect.Type, outType reflect.Type, input interface{}) (interface{}, error) {
				if inType.Kind() == reflect.Map && outType.Kind() == reflect.Struct { // we'll decoding a struct
					fieldsMap := make(map[string]reflect.StructField, 10)
					for i := 0; i < outType.NumField(); i++ {
						fieldsMap[strings.ToLower(outType.Field(i).Name)] = outType.Field(i)
					}
					inputMap := input.(map[string]interface{})
					extraConfig := make(map[string]interface{}, 5)
					inputMap["ExtraConfig"] = extraConfig
					for key := range inputMap {
						_, ok := fieldsMap[strings.ToLower(key)]
						if !ok {
							extraConfig[key] = inputMap[key]
						}
					}
				}
				return input, nil
			},
			c.DecodeHook)
	})
	if err != nil {
		fmt.Printf("Struct config error: %s", err)
	}
	Config = config
	UpdateLogLevel()
	return true, nil
}

func LevelStrParse(levelStr string) (level logrus.Level) {
	level = logrus.InfoLevel
	if levelStr == "debug" {
		level = logrus.DebugLevel
	} else if levelStr == "info" {
		level = logrus.InfoLevel
	} else if levelStr == "warn" {
		level = logrus.WarnLevel
	} else if levelStr == "error" {
		level = logrus.ErrorLevel
	}
	return level
}

func UpdateLogLevel() {
	if ConsoleHook != nil {
		level := LevelStrParse(Config.LogLevel)
		ConsoleHook.LogLevel = level
		logrus.Printf("Set logrus console level to %s", level)
	}
}

func PrepareConfig() {
	confPath := flag.String("config", "config.json", "config.json location")
	flag.Parse()
	viper.SetConfigFile(*confPath)
	InitConfig()
}
