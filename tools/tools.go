package tools

import (
	"os"
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
	pinger.Timeout = time.Second * 3
	pinger.Count = 4
	done := make(chan struct{})
	timeout := make(chan struct{})
	go func() {
		tick := time.NewTicker(time.Second*2 + time.Millisecond*900)
		select {
		case <-tick.C:
			close(timeout)
			pinger.Stop()
		case <-done:
		}
	}()
	err = pinger.Run()
	close(done)
	if err != nil {
		return false, err
	}

	select {
	case <-timeout:
		return false, nil
	default:
	}

	return true, nil
}

func Getenv(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		val = def
	}

	return val
}
