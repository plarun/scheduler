package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	MYSQL_CONFIG_FILE = "/root/go/src/github.com/plarun/scheduler/config/allocator/sql/mysql/query.json"
)

var qry *mySQLQuery

type mySQLQuery struct {
	CheckJobExist       string `json:"job_exists"`
	QueryJobId          string `json:"get_job_id"`
	QueryJob            string `json:"get_job"`
	QueryJobStatus      string `json:"get_job_status"`
	InsertJob           string `json:"insert_job"`
	DeleteJob           string `json:"delete_job"`
	InsertJobStarttimes string `json:"insert_start_times"`
	DeleteJobStarttimes string `json:"delete_start_times"`
	InsertJobStartmins  string `json:"insert_start_mins"`
	DeleteJobStartmins  string `json:"delete_start_mins"`
	InsertJobRelation   string `json:"insert_job_relation"`
	DeleteJobRelation   string `json:"delete_job_relation"`
	UpdateJobRunwindow  string `json:"update_run_window"`
	QueryRunflag        string `json:"get_run_flag"`
	UpdateRunflag       string `json:"update_run_flag"`
}

func GetQuery() *mySQLQuery {
	return qry
}

func LoadConfig() error {
	viper.SetConfigType("json")
	viper.SetConfigFile(MYSQL_CONFIG_FILE)

	query := &mySQLQuery{}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	if err := viper.Unmarshal(&query); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}
