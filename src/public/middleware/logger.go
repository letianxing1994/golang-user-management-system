package middleware

import (
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const (
	LogPath = "/Users/tianxingle/go/src/Entry_Task/src/log"
)

func AddHook(filename string) (*rotatelogs.RotateLogs, error) {
	LogWriter, err := rotatelogs.New(
		filename+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	return LogWriter, err
}

func LoggerMiddleware() gin.HandlerFunc {
	//get logrus instance
	logClient := log.New()

	//build file
	filename := path.Join(LogPath, "gin_api.log")
	warnFilename := path.Join(LogPath, "warn.log")
	errorFilename := path.Join(LogPath, "error.log")
	debugFilename := path.Join(LogPath, "debug.log")
	panicFilename := path.Join(LogPath, "panic.log")
	fatalFilename := path.Join(LogPath, "fatal.log")

	//build warning files
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 777)
	if err != nil {
		os.Exit(1)
	}

	//set log output path
	logClient.Out = f
	logClient.SetFormatter(&log.TextFormatter{DisableColors: false})
	logClient.SetLevel(log.DebugLevel)

	//add hook
	warnLogWriter, err := AddHook(warnFilename)
	errorLogWriter, err := AddHook(errorFilename)
	debugLogWriter, err := AddHook(debugFilename)
	panicLogWriter, err := AddHook(panicFilename)
	fatalLogWriter, err := AddHook(fatalFilename)

	writeMap := lfshook.WriterMap{
		log.DebugLevel: debugLogWriter,
		log.WarnLevel:  warnLogWriter,
		log.FatalLevel: fatalLogWriter,
		log.ErrorLevel: errorLogWriter,
		log.PanicLevel: panicLogWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{})
	logClient.AddHook(lfHook)

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		//elapsed time
		latency := end.Sub(start)
		//request url
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if statusCode > 499 {
			logClient.Errorf("| %3d | %13v | %15s | %s %s | %s |",
				statusCode,
				latency,
				clientIP,
				method,
				path)
		} else if statusCode > 399 {
			logClient.Warnf("| %3d | %13v | %15s | %s %s | %s |",
				statusCode,
				latency,
				clientIP,
				method,
				path)
		}
	}
}
