package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	errs "github.com/chloyka/gorig/utils/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	file           *os.File
	bufferedWriter *zapcore.BufferedWriteSyncer
}

func newLogger(cfg *configTypes.LoggerConfig) (*Logger, error) {
	if err := os.MkdirAll(cfg.LogsDir, 0755); err != nil {
		return nil, errs.Wrap(errs.ErrLoggerCreateDir, err)
	}

	fileNum := getNextFileNumber(cfg.LogsDir, cfg.MaxLogFiles)
	filename := fmt.Sprintf("%d-%s.log", fileNum, time.Now().Format(time.DateTime))
	filePath := filepath.Join(cfg.LogsDir, filename)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errs.Wrap(errs.ErrLoggerOpenFile, err)
	}

	bufferedWriter := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(file),
		Size:          cfg.BufferSize,
		FlushInterval: cfg.FlushInterval,
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		bufferedWriter,
		parseLogLevel(cfg.Level),
	)

	logger := zap.New(core)

	return &Logger{
		Logger:         logger,
		file:           file,
		bufferedWriter: bufferedWriter,
	}, nil
}

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func getNextFileNumber(logsDir string, maxLogFiles int) int {
	entries, err := os.ReadDir(logsDir)
	if err != nil || len(entries) == 0 {
		return 1
	}

	var nums []int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".log") {
			continue
		}

		parts := strings.SplitN(name, "-", 2)
		if len(parts) == 0 {
			continue
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		nums = append(nums, n)
	}

	if len(nums) == 0 {
		return 1
	}

	sort.Ints(nums)
	lastNum := nums[len(nums)-1]

	nextNum := lastNum + 1
	if nextNum > maxLogFiles {
		nextNum = 1
	}

	return nextNum
}

func (l *Logger) Close() error {

	if l.bufferedWriter != nil {
		_ = l.bufferedWriter.Sync()
		_ = l.bufferedWriter.Stop()
	}
	_ = l.Logger.Sync()
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
