package logging

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var config = createConfig()
var logger = createLogger()

const defaultLevel = zapcore.InfoLevel

func createConfig() zap.Config {
	return zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(defaultLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
}

func getLevelFromConfig(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultLevel.String()
}

func SetLevel() {
	atomicLevel := zap.NewAtomicLevel()
	level := getLevelFromConfig("logging.level")

	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		Err(err)
		Warn(fmt.Sprintf("logging level failed to be set to '%s'", level))
	} else {
		config.Level.SetLevel(atomicLevel.Level())
		logger = createLogger()
		Info(fmt.Sprintf("config level has been changed to '%s'", level))
	}
}

func createLogger() *zap.SugaredLogger {
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() //nolint:errcheck
	return logger.Sugar()
}

func GetRouterLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
			)
		}
	}
}

func Debug(message string) {
	logger.Debug(message)
}

func Info(message string) {
	logger.Info(message)
}

func Warn(message string) {
	logger.Warn(message)
}

func Err(err error) {
	logger.Error(err)
}
