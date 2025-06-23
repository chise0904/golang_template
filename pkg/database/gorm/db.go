package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(config *Config) (Connection, error) {
	var (
		conn Connection
	)
	readDB, err := setupDatabase(&config.Read)
	if err != nil {
		return conn, err
	}

	log.Info().Msg("connect to read db success")

	writeDB, err := setupDatabase(&config.Write)
	if err != nil {
		return conn, err
	}

	log.Info().Msg("connect to write db success")

	conn = Connection{
		ReadDB:  readDB,
		WriteDB: writeDB,
	}
	return conn, nil
}

func setupDatabase(database *Database) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(10) * time.Second

	if database.WriteTimeout == "" {
		database.WriteTimeout = "10s"
	}

	if database.ReadTimeout == "" {
		database.ReadTimeout = "10s"
	}

	var dialector gorm.Dialector

	switch database.Type {
	case MySQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true&readTimeout=%s&writeTimeout=%s", database.User, database.Password, database.Host+":"+strconv.Itoa(database.Port), database.Name, database.ReadTimeout, database.WriteTimeout)
		dialector = mysql.Open(dsn)
	case Postgres:
		dsn := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s `, database.User, database.Password, database.Host, database.Port, database.Name)
		if database.SSLEnable {
			dsn += " sslmode=require"
		} else {
			dsn += " sslmode=disable"
		}
		if strings.TrimSpace(database.SearchPath) != "" {
			dsn = fmt.Sprintf("%s search_path=%s", dsn, database.SearchPath)
		}
		dialector = postgres.Open(dsn)
	default:
		return nil, errors.New("Not support driver")
	}

	colorful := false
	logLevel := logger.Silent
	if database.Debug {
		colorful = true
		logLevel = logger.Info
	}

	newLogger := NewLogger(logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logLevel,    // Log level
		Colorful:      colorful,    // Disable color
	})

	var conn *gorm.DB

	// 嘗試重新連線database
	err := backoff.Retry(func() error {
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger:         newLogger,
			TranslateError: true,
		})
		if err != nil {
			return err
		}
		conn = db
		sqlDB, err := conn.DB()
		if err != nil {
			log.Error().Msgf("err %s", err.Error())

			return err
		}

		err = sqlDB.Ping()
		if err != nil {
			log.Error().Msgf("err %s", err.Error())
			return err
		}
		return nil

	}, bo)

	if err != nil {
		return nil, err
	}

	// set default idle conn
	if database.MaxIdleConn == 0 {
		database.MaxIdleConn = 10
	}

	if database.MaxOpenConn == 0 {
		database.MaxOpenConn = 20
	}

	if database.MaxLifetimeSec == 0 {
		database.MaxLifetimeSec = 14400
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(database.MaxIdleConn)
	sqlDB.SetMaxOpenConns(database.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(database.MaxLifetimeSec))
	return conn, nil
}
