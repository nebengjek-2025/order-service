package mysql

import (
	"fmt"
	"order-service/src/pkg/log"
	"order-service/src/pkg/utils"
	"runtime"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DBConn variable to declare Database Connection
var DBConn *DatabaseConnection

type (
	// DBInterface to provide general Func
	DBInterface interface {
		Connect(string) *DatabaseConnection
		GetDB() (*sqlx.DB, error)
	}

	// Database to provide Database Config
	Database struct {
		Name DatabaseConfig
	}

	// DatabaseConfig currently have master only
	DatabaseConfig struct {
		Master string
	}

	// DatabaseConnection provide struct sqlx connection
	DatabaseConnection struct {
		Connection *sqlx.DB
	}
)

var (
	accessOnce sync.Once
	access     DBInterface
)

func InitConnection(viper *viper.Viper, log log.Log) (DBInterface, error) {
	if access != nil {
		return access, nil
	}

	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		username, password, host, port, database)

	accessOnce.Do(func() {
		dbClient := NewDatabase(dsn)
		conn := dbClient.Connect(database)
		if conn.Connection == nil {
			log.Error("mysql", "Failed to establish MySQL connection", "InitConnection", "")
		}
		access = dbClient
	})

	if access == nil {
		return nil, fmt.Errorf("failed to initialize database connection")
	}

	return access, nil
}

func StatusConnection() (*sqlx.DB, error) {
	if DBConn == nil || DBConn.Connection == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	stats := DBConn.Connection.Stats()
	log.GetLogger().Info("mysql", fmt.Sprintf("DB Status: %+v", stats), "StatusConnection", "")

	return DBConn.Connection, nil
}

// NewDatabase create Database Struct from Config
func NewDatabase(dsn string) *Database {

	dc := DatabaseConfig{
		Master: dsn,
	}

	return &Database{
		Name: dc,
	}
}

// Connect provide sqlx connection
func (db *Database) Connect(dbName string) *DatabaseConnection {

	databaseConn := DatabaseConnection{}

	master := db.Name.Master
	if master != "" {
		db, err := sqlx.Connect("mysql", master)
		if err != nil {
			log.GetLogger().Error("mysql", "Can not connect MySQL", "connect", utils.ConvertString(err))
		}

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
		db.SetConnMaxLifetime(time.Minute * 10)

		databaseConn.Connection = db
	}

	DBConn = &DatabaseConnection{Connection: databaseConn.Connection}
	return DBConn
}

func (db *Database) GetDB() (*sqlx.DB, error) {
	newDB := DBConn.Connection
	if newDB.Stats().OpenConnections > 40 {
		fpcs := make([]uintptr, 1)
		n := runtime.Callers(2, fpcs)
		if n != 0 {
			fun := runtime.FuncForPC(fpcs[0] - 1)
			if fun != nil {
				log.GetLogger().Error("mysql", fmt.Sprintf("Db Conn more than 40, Caller from Func : %s", fun.Name()), "GetDB", string(fun.Name()))
			}
		}
		log.GetLogger().Info("mysql", fmt.Sprintf("DB Conn more than 40, currently : %s", utils.ConvertString(newDB.Stats())), "GetDB", utils.ConvertString(newDB.Stats()))
	}
	return newDB, nil
}
