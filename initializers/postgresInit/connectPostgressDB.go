package postgresInit

import (
	"database/sql"
	"fmt"
	"github.com/akashc777/csvToPdf/helpers"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	GormDB    *gorm.DB
	SqlDB     *sql.DB
	DBTimeout time.Duration
}

var DBConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute
const dbTimeout = time.Second * 3

func ConnectPostgres(dsn string) error {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(maxOpenDBConn)
	sqlDB.SetMaxIdleConns(maxIdleDbConn)
	sqlDB.SetConnMaxLifetime(maxDbLifetime)

	err = testDB(sqlDB)
	if err != nil {
		return err
	}

	DBConn.GormDB = db
	DBConn.SqlDB = sqlDB
	DBConn.DBTimeout = dbTimeout

	return nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()

	if err != nil {
		fmt.Printf("Error : %+v", err)
		return err
	}

	helpers.MessageLogs.InfoLog.Println("Successfully connected with postgres server !")

	return nil
}
