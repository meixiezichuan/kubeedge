package main

import (
	"os"

	"k8s.io/component-base/logs"

	"github.com/kubeedge/kubeedge/edgedevice/cmd/edgedevice/app"
)

func main() {
	command := app.NewEdgeDeviceCommand()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
