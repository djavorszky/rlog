package rlog

import (
	"bytes"
	fmt "fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	mu       sync.RWMutex
	registry = make(map[int32]service)
	counter  = make(chan int32)
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	go increment()
}

func increment() int32 {
	var id int32
	for {
		counter <- id
		id++
	}
}

type service struct {
	app  string
	name string
}

// SetOut sets the output of the logger. It is the responsibility of the
// caller to close the resource. Not calling it will result in rlog printing
// to stderr
func SetOut(out io.Writer) {
	logrus.SetOutput(out)
	logrus.SetFormatter(logger{})
}

// Server implements implements the methods of the rlog service
type Server struct{}

// Register informs the server that a new service has joined
func (s *Server) Register(ctx context.Context, msg *RegisterRequest) (*RegisterResponse, error) {
	id := <-counter

	mu.Lock()
	registry[id] = service{app: msg.App, name: msg.Service}
	mu.Unlock()

	logrus.WithFields(logrus.Fields{
		"app":     msg.App,
		"service": msg.Service,
	}).Info("Registered Service")

	return &RegisterResponse{Id: id, Message: "REGISTER_SUCCESS"}, nil
}

// Debug will log the received message on debug level
func (s *Server) Debug(ctx context.Context, msg *LogMessage) (*LogResponse, error) {
	meta, err := getMeta(msg.Id)
	if err != nil {
		return &LogResponse{}, err
	}

	meta.Debug(msg.Message)

	return &LogResponse{}, nil
}

// Fatal will log the received message on error level - fatal would kill rlog.
func (s *Server) Fatal(ctx context.Context, msg *LogMessage) (*LogResponse, error) {
	meta, err := getMeta(msg.Id)
	if err != nil {
		return &LogResponse{}, err
	}

	// TODO: In case Fatal is called, send a notification somewhere
	meta.WithField("error", "critical").Error(msg.Message)

	return &LogResponse{}, nil
}

// Error will log the received message on error level
func (s *Server) Error(ctx context.Context, msg *LogMessage) (*LogResponse, error) {
	meta, err := getMeta(msg.Id)
	if err != nil {
		return &LogResponse{}, err
	}

	meta.Error(msg.Message)

	return &LogResponse{}, nil
}

// Warn will log the received message on warn level
func (s *Server) Warn(ctx context.Context, msg *LogMessage) (*LogResponse, error) {
	meta, err := getMeta(msg.Id)
	if err != nil {
		return &LogResponse{}, err
	}

	meta.Warn(msg.Message)

	return &LogResponse{}, nil
}

// Info will log the received message on info level
func (s *Server) Info(ctx context.Context, msg *LogMessage) (*LogResponse, error) {
	meta, err := getMeta(msg.Id)
	if err != nil {
		return &LogResponse{}, err
	}

	meta.Info(msg.Message)

	return &LogResponse{}, nil
}

func getMeta(id int32) (*logrus.Entry, error) {
	service, err := getService(id)
	if err != nil {
		return &logrus.Entry{}, err
	}

	return logrus.WithFields(logrus.Fields{
		"app":     service.app,
		"service": service.name,
	}), nil
}

func getService(id int32) (service, error) {
	mu.RLock()
	defer mu.RUnlock()

	s, ok := registry[id]
	if !ok {
		return service{}, fmt.Errorf("ERR_SERVICE_NOT_REGISTERED")
	}

	return s, nil
}

type logger struct{}

func (l logger) Format(entry *logrus.Entry) ([]byte, error) {
	line := fmt.Sprintf("%s - [%s::%s] [%s] %s", entry.Time.Format(time.RFC3339), entry.Data["app"], entry.Data["service"], entry.Level, entry.Message)

	if len(entry.Data) > 2 {
		var buf bytes.Buffer

		for k, v := range entry.Data {
			if k == "app" || k == "service" {
				continue
			}

			val, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("couldn't cast value of %q to string", k)
			}

			buf.WriteString(", ")
			buf.WriteString(k)
			buf.WriteString("=")
			buf.WriteString(val)
			buf.WriteString("")
		}

		line = fmt.Sprintf("%s (%s)", line, strings.TrimLeft(buf.String(), ", "))
	}

	return []byte(line + "\n"), nil
}
