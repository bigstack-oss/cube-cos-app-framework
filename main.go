package main

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd"
	log "go-micro.dev/v5/logger"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Infof("failed to execute command(%v)", err)
	}
}
