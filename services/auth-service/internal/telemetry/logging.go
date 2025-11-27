package telemetry

import (
	"io"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type LogstashFormatter struct {
	ServiceName string
}

func (f *LogstashFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(map[string]interface{})
	data["@timestamp"] = entry.Time.Format("2006-01-02T15:04:05.000Z")
	data["level"] = entry.Level.String()
	data["message"] = entry.Message
	data["service"] = f.ServiceName
	data["environment"] = "development"

	// Add all fields from the log entry
	for k, v := range entry.Data {
		data[k] = v
	}

	// Convert to JSON
	serialized, err := logrus.StandardLogger().JSONFormatter.Format(&logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
	})
	if err != nil {
		return nil, err
	}

	return serialized, nil
}

// InitLogger initializes structured logging with Logstash integration
func InitLogger(serviceName, logstashHost string) *Logger {
	logger := logrus.New()

	// Set JSON formatter for structured logs
	logger.SetFormatter(&LogstashFormatter{ServiceName: serviceName})

	// Set log level
	logger.SetLevel(logrus.InfoLevel)

	// Configure multiple outputs
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	// Try to connect to Logstash
	if logstashHost != "" {
		if conn, err := net.Dial("tcp", logstashHost); err == nil {
			conn.Close() // Just test connection
			logger.Info("Logstash connection test successful")
			// Note: For production, you would use a proper Logstash hook
			// For now, we'll just log to stdout in JSON format
		} else {
			logger.WithError(err).Warn("Could not connect to Logstash, logging to stdout only")
		}
	}

	// Set multi-writer
	logger.SetOutput(io.MultiWriter(writers...))

	logger.WithFields(logrus.Fields{
		"service": serviceName,
	}).Info("Logger initialized")

	return &Logger{logger}
}
