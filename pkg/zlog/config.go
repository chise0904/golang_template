package zlog

type Config struct {
	Env            string `mapstructure:"env"`
	AppID          string `mapstructure:"app_id"`
	Level          int8   `mapstructure:"level"` // 0 Debug 1 INFO 2 Warn 3 Error
	Debug          bool   `mapstructure:"debug"`
	EnableCaller   bool   `mapstructure:"enable_caller"`
	CallerMinLevel int8   `mapstructure:"caller_min_level"`
}
