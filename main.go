package main

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/runtime"
	log "go-micro.dev/v5/logger"
)

func main() {
	err := runtime.Init()
	if err != nil {
		log.Errorf("failed to initialize runtime(%v)", err)
		return
	}

	err = cmd.Execute()
	if err != nil {
		log.Infof("failed to execute command(%v)", err)
	}
}
