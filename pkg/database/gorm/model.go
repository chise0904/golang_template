package db

import (
	"github.com/chise0904/golang_template/pkg/time"
	"gorm.io/gorm"
)

type DatabaseType string

const (
	// MySQL ...
	MySQL DatabaseType = "mysql"
	// Postgres ...
	Postgres DatabaseType = "postgres"
)

type Config struct {
	Read    Database `mapstructure:"read"`
	Write   Database `mapstructure:"write"`
	Secrets string   `mapstructure:"secrets"`
}

type Connection struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

type Database struct {
	Debug          bool         `mapstructure:"debug"`
	Host           string       `mapstructure:"host"`
	User           string       `mapstructure:"user"`
	Port           int          `mapstructure:"port"`
	Password       string       `mapstructure:"password"`
	Name           string       `mapstructure:"name"`
	Type           DatabaseType `mapstructure:"type"`
	MaxIdleConn    int          `mapstructure:"max_idle_conn"`
	MaxOpenConn    int          `mapstructure:"max_open_conn"`
	MaxLifetimeSec int          `mapstructure:"max_lifetime"`
	ReadTimeout    string       `mapstructure:"read_timeout"`
	WriteTimeout   string       `mapstructure:"write_timeout"`
	SSLEnable      bool         `mapstructure:"ssl_enable"`
	SearchPath     string       `mapstructure:"search_path"`
}

type CommonEmbedding struct {
	Creator   string `gorm:"column:creator"`
	CreatedAt int64  `gorm:"column:created_at"`
	Updater   string `gorm:"column:updater" `
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (c *CommonEmbedding) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.NowMS()
	c.CreatedAt = now
	c.UpdatedAt = now
	return nil
}

func (c *CommonEmbedding) BeforeSave(tx *gorm.DB) (err error) {
	now := time.NowMS()
	c.UpdatedAt = now
	return nil
}
