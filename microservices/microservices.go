package microservices

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/gofiber/fiber/v2"
)

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Stop()
	Cleanup() error
	Log(tag string, message string)

	// HTTP Services
	Use(args ...interface{})
	Group(prefix string, h ...ServiceHandleFunc)

	// HTTP Methods
	GET(path string, h ...ServiceHandleFunc)
	POST(path string, h ...ServiceHandleFunc)
	PUT(path string, h ...ServiceHandleFunc)
	PATCH(path string, h ...ServiceHandleFunc)
	DELETE(path string, h ...ServiceHandleFunc)
}

// Microservice is the centralized service management
type Microservice struct {
	fiber       *fiber.App
	exitChannel chan bool
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx *fiber.Ctx) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{
		fiber: fiber.New(config.FiberConfig),
	}
}

// Start start all registered services
func (ms *Microservice) Start() error {

	httpN := len(ms.fiber.Stack())
	var exitHTTP chan bool
	if httpN > 0 {
		exitHTTP = make(chan bool, 1)
		go func() {
			ms.startHTTP(exitHTTP)
		}()
	}

	// There are 2 ways to exit from Microservices
	// 1. The SigTerm can be send from outside program such as from k8s
	// 2. Send true to ms.exitChannel
	osQuit := make(chan os.Signal, 1)
	ms.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-osQuit:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		case <-ms.exitChannel:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		}
	}

	ms.Cleanup()
	return nil
}

// Stop stop the services
func (ms *Microservice) Stop() {
	if ms.exitChannel == nil {
		return
	}
	ms.exitChannel <- true
}

// Cleanup clean resources up from every registered services before exit
func (ms *Microservice) Cleanup() error {
	ms.Log("Microservices", "Start cleanup")

	if database.DBConn != nil {
		sqlDB, _ := database.DBConn.DB()
		sqlDB.Close()
	}
	if cache.Client != nil {
		cache.Client.Close()
	}

	return nil
}

// Log log message to console
func (ms *Microservice) Log(tag string, message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	log.Println(
		"[",
		time.Now().Format("2006-01-02 15:04:05"),
		"]",
		tag+":",
		fns[len(fns)-1],
		line,
		"|",
		message,
	)
}
