package main

import (
	"time"

	"github.com/dtmistry/swarm-tool/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	formatter := &log.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = time.RFC3339Nano
	log.SetFormatter(formatter)
	cmd.Execute()
}
