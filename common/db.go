package common

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DbConnections *DbConnection

type DbConnection struct {
	db map[string]*gorm.DB
}

func NewDbConnection() *DbConnection {
	return &DbConnection{db: make(map[string]*gorm.DB)}
}
func (db *DbConnection) Get(database string) *gorm.DB {
	if _, ok := db.db[database]; !ok {
		if err := connections(database); err != nil {
			panic(err)
		}
	}
	return db.db[database]
}
func connections(name string) error {
	username := GlobalConfig.GetString("db.username")
	password := GlobalConfig.GetString("db.password")
	host := GlobalConfig.GetString("db.host")
	database := GlobalConfig.GetString("db.database")
	charset := GlobalConfig.GetString("db.charset")
	parseTime := GlobalConfig.GetString("db.parseTime")
	loc := GlobalConfig.GetString("db.loc")
	//这里的loc需要http_parse
	loc = url.QueryEscape(loc)
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%v&loc=%s", username, password, host, database, charset, parseTime, loc)
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		return err
	}
	db.LogMode(true)
	DbConnections.db[name] = db
	return nil
}
