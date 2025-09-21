//go:build auto

package dytapp

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type DefaultApp struct {
	version   string
	processID string
}

func NewApp(version string) *DefaultApp {
	app := new(DefaultApp)
	app.processID = "development"
	app.version = version
	return app
}

func (dytapp *DefaultApp) Run(debug bool) {
	setPressureTestEnv()
	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c
	fmt.Printf("dayon auto robot close signal : %v", sig)
}
