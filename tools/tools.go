package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/playmixer/corvid/logger"
)

var (
	log      = logger.New("tools")
	timeout  time.Duration
	count    int
	timeWite time.Duration
	osystem  string
)

func Init() {
	var err error
	osystem = runtime.GOOS

	eTimeout := Getenv("PING_TIMEOUT", "4000")
	eCount := Getenv("PING_COUNT", "3")
	eTimeWite := Getenv("PING_TIMEWITE", "3800")

	timeout, err = time.ParseDuration(eTimeout + "ms")
	if err != nil {
		log.ERROR(err.Error())
		timeout = time.Second * 4
	}
	count, err = strconv.Atoi(eCount)
	if err != nil {
		log.ERROR(err.Error())
		count = 3
	}
	timeWite, err = time.ParseDuration(eTimeWite + "ms")
	if err != nil {
		log.ERROR(err.Error())
		timeWite = time.Millisecond * 3800
	}
}

func Ping(ip string) (bool, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return false, err
	}

	pinger.SetPrivileged(true)
	pinger.Timeout = timeout
	pinger.Count = count
	done := make(chan struct{})
	timeout := make(chan struct{})
	go func() {
		tick := time.NewTicker(timeWite)

		select {
		case <-tick.C:
			close(timeout)
			tick.Stop()
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

func Ping2(ip string) (bool, error) {
	// var err error
	ctx := context.Background()

	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	key := "-c"
	if osystem == "windows" {
		key = "-n"
	}

	output, err := exec.CommandContext(ctx, "ping", key, fmt.Sprint(count), ip).CombinedOutput()
	if err != nil {
		// log.ERROR(err.Error())
		return false, nil
	}
	if osystem == "windows" && !strings.Contains(string(output), "TTL") {
		return false, nil
	}
	// fmt.Println(string(output))

	return true, nil
}

func Getenv(name, def string) string {
	val := os.Getenv(name)
	if val == "" {
		val = def
	}

	return val
}
