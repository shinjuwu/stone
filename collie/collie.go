package collie

import (
	"collie/cluster"
	"collie/conf"
	"collie/console"
	"collie/log"
	"collie/module"
	"os"
	"os/signal"
)

func Run(logLevel int32, mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New("", logLevel)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
	}

	log.Info("dayon starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Info("dayon closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
