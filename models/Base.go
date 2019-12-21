package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type Payload struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var db *gorm.DB        //database
var db_sas *gorm.DB    //database
var db_master *gorm.DB //database

func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")
	db_name_sas := os.Getenv("DB_SEMAYA_SAS")
	db_name_master := os.Getenv("DB_SEMAYA_MASTER")
	// dbHost := os.Getenv("db_host")

	dbURI := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, db_name)
	dbURISas := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, db_name_sas)
	dbURIMaster := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, db_name_master)

	conn, err := gorm.Open("mysql", dbURI)
	connSas, err := gorm.Open("mysql", dbURISas)
	connMaster, err := gorm.Open("mysql", dbURIMaster)

	if err != nil {
		fmt.Print(err)
	}

	db = conn.LogMode(true)
	db_sas = connSas
	db_master = connMaster

	db.Debug().AutoMigrate(&Attendance{}) //Database migration
}

//returns a handle to the DB object
func GetDB(whatDB string) *gorm.DB {
	if whatDB == "sas" {
		return db_sas
	} else if whatDB == "master" {
		return db_master
	} else if whatDB == "main" {
		return db
	} else {
		return db_master
	}
}
