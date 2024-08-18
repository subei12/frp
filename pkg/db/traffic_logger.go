package db

import (
	"database/sql"
	"fmt"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	db  *sql.DB
	cfg v1.MySQLServerConfig
)

func InitDB(mysqlCfg v1.MySQLServerConfig) {
	// 将配置保存到包级别变量
	cfg = mysqlCfg

	// 构建数据源名称
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Errorf("Error opening database: %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		log.Errorf("Error connecting to the database: %v\n", err)
	}
	log.Infof("Database connection established successfully")
}

func LogDailyTraffic(name string, proxyType string, direction string, trafficBytes int64) {
	if cfg.Port <= 0 {
		return
	}

	date := time.Now().Format(time.DateOnly)
	var query string

	if direction == "in" {
		query = `
			INSERT INTO daily_traffic (name, proxy_type, date, traffic_in)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE traffic_in = traffic_in + ?`
	} else if direction == "out" {
		query = `
			INSERT INTO daily_traffic (name, proxy_type, date, traffic_out)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE traffic_out = traffic_out + ?`
	}

	// 执行 SQL 语句
	_, err := db.Exec(query, name, proxyType, date, trafficBytes, trafficBytes)
	if err != nil {
		log.Errorf("Error updating daily traffic log: %v\n", err)
	}
	log.Infof("Daily traffic log updated: name=%s, proxy_type=%s, date=%s, traffic_%s=%d", name, proxyType, date, direction, trafficBytes)
}
