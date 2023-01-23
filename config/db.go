package config

const (
	DB_CONFIG_FILE = "/root/go/src/github.com/plarun/scheduler/config/db.json"
)

var dbCfg *DBConfig

type DBConfig struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Schema   string `json:"Schema"`
}

func GetDBConfig() *DBConfig {
	return dbCfg
}
