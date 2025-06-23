package web

import "github.com/chise0904/golang_template/pkg/pagination"

type Config struct {
	Mode           string `mapstructure:"mode"`
	Port           string `mapstructure:"port"`
	RequestDump    bool   `mapstructure:"request_dump"`
	MaxLogBodySize int    `mapstructure:"max_log_body_size"`
}

type ResponsePayLoadMetaData struct {
	*pagination.Pagination
} //`json:"meta"`
