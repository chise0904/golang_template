package db

import (
	"context"
	"fmt"
	"time"

	"github.com/chise0904/golang_template/pkg/trace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
	Teal        = "\u001B[1;36m"
)

type gormLogger struct {
	LogLevel                            logger.LogLevel
	Config                              logger.Config
	SlowThreshold                       time.Duration
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewLogger(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s "
		warnStr      = "%s "
		errStr       = "%s "
		traceStr     = "%s "
		traceWarnStr = "%s "
		traceErrStr  = "%s "
	)

	l := &gormLogger{
		LogLevel:      config.LogLevel,
		SlowThreshold: config.SlowThreshold,
		Config:        config,
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}

	return l
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	return g
}

func (g *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	currentLogger := log.With().Str("request_id", trace.XRequestIDFromContext(ctx)).Logger()
	if g.LogLevel >= logger.Info {
		currentLogger.Info().Msgf(msg, data...)
	}
}

func (g *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {

	currentLogger := log.With().Str("request_id", trace.XRequestIDFromContext(ctx)).Logger()

	if g.LogLevel >= logger.Warn {
		currentLogger.Warn().Msgf(msg, data...)
	}
}

func (g *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	var (
		currentLogger zerolog.Logger
		errs          []error
	)
	currentLogger = log.With().Str("request_id", trace.XRequestIDFromContext(ctx)).Logger()

	if g.LogLevel >= logger.Error {
		for i := 0; i < len(data); i++ {
			if err, ok := data[i].(error); ok && err != nil {
				errs = append(errs, err)
			}
		}
		currentLogger.Error().Errs("errs", errs).Msgf(msg, data...)
	}
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	currentLogger := log.With().Str("request_id", trace.XRequestIDFromContext(ctx)).Logger()

	if g.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && g.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				currentLogger.
					Error().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Str("sql_row", "-").
					Err(err).
					Str("file", utils.FileWithLineNum()).
					Msg(sql)
			} else {
				currentLogger.
					Error().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", rows).
					Err(err).
					Str("file", utils.FileWithLineNum()).
					Msg(sql)
			}
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
			if rows == -1 {
				currentLogger.
					Warn().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Str("sql_row", "-").
					Str("sql_slow_log", slowLog).
					Str("file", utils.FileWithLineNum()).
					Msg(sql)
			} else {
				currentLogger.
					Warn().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", rows).
					Str("sql_slow_log", slowLog).
					Str("file", utils.FileWithLineNum()).
					Msg(sql)
			}
		case g.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				currentLogger.
					Info().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Str("sql_row", "-").
					Str("file", utils.FileWithLineNum()).
					Msg(sql)
			} else {
				currentLogger.
					Info().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", rows).
					Str("file", utils.FileWithLineNum()).
					Msg(sql)
			}
		}
	}
}
