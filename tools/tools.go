package tools

import (
	"time"

	"github.com/go-ping/ping"
	"github.com/playmixer/corvid/logger"
)

var (
	log = logger.New("tools")
)

func Ping(ip string) (bool, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return false, err
	}

	pinger.SetPrivileged(true)
	pinger.Timeout = time.Second * 2
	pinger.Count = 1
	err = pinger.Run()
	if err != nil {
		return false, err
	}

	return true, nil
}
