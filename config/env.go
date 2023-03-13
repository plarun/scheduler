package config

const (
	ENV_CONFIG_FILE = "/root/go/src/github.com/plarun/scheduler/config/env.json"
)

var appCfg *AppConfig

type logs struct {
	Path       string `json:"Path"`
	Prefix     string `json:"Prefix"`
	DateFormat string `json:"DateFormat"`
	Extension  string `json:"Extension"`
}

type AppConfig struct {
	Debug   bool   `json:"Debug"`
	AppName string `json:"AppName"`
	AppRoot string `json:"AppRoot"`
	Service struct {
		Client struct {
			Name string `json:"Name"`
			Logs logs   `json:"Logs"`
		} `json:"Client"`
		EventServer struct {
			Name string `json:"Name"`
			Port int    `json:"Port"`
			Logs logs   `json:"Logs"`
		} `json:"Eventserver"`
		Allocator struct {
			Name string `json:"Name"`
			Port int    `json:"Port"`
			Logs logs   `json:"Logs"`
		} `json:"Allocator"`
		Controller struct {
			Name string `json:"Name"`
			Port int    `json:"Port"`
			Logs logs   `json:"Logs"`
		} `json:"Controller"`
		Validator struct {
			Name string `json:"Name"`
			Port int    `json:"Port"`
			Logs logs   `json:"Logs"`
		} `json:"Validator"`
		Worker struct {
			Name string `json:"Name"`
			Port int    `json:"Port"`
			Logs logs   `json:"Logs"`
		} `json:"Worker"`
	} `json:"Service"`
}

func GetAppConfig() *AppConfig {
	return appCfg
}
