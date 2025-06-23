package grpc

// Config ...
type Config struct {
	Mode           string `mapstructure:"mode"`          // Mode = debug and release
	Port           string `mapstructure:"port"`          // Port = :8080
	RequestDump    bool   `mapstructure:"request_dump"`  // true or false
	ResponseDump   bool   `mapstructure:"response_dump"` // true or false
	MaxLogBodySize int    `mapstructure:"max_log_body_size"`
}
