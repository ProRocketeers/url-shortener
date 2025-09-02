package infrastructure

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	gormlogger "gorm.io/gorm/logger"
)

type ZerologGormLogger struct {
	log                zerolog.Logger
	level              gormlogger.LogLevel
	slowQueryThreshold time.Duration
}

func NewZerologGormLogger(log zerolog.Logger, level gormlogger.LogLevel, slowQueryThreshold time.Duration) *ZerologGormLogger {
	return &ZerologGormLogger{
		log:                log,
		level:              level,
		slowQueryThreshold: slowQueryThreshold,
	}
}

// implements `gorm/logger.Interface`
func (l *ZerologGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *ZerologGormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Info {
		l.log.Info().Msgf(msg, args...)
	}
}

func (l *ZerologGormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Warn {
		l.log.Warn().Msgf(msg, args...)
	}
}

func (l *ZerologGormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Error {
		l.log.Error().Msgf(msg, args...)
	}
}

func (l *ZerologGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	event := l.log.With().
		Str("sql", sql).
		Dur("elapsedMs", elapsed).
		Int64("affected", rows).
		Logger()

	switch {
	case err != nil && l.level >= gormlogger.Error:
		event.Error().Err(err).Msg("GORM error")
	case elapsed > l.slowQueryThreshold && l.level > gormlogger.Warn:
		event.Warn().Msg("Slow SQL query")
	case l.level >= gormlogger.Info:
		event.Debug().Msg("SQL executed")
	}
}
