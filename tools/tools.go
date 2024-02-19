package tools

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
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

func isWindows() bool {
	return osystem == "windows"
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
	if isWindows() {
		key = "-n"
	}

	output, err := exec.CommandContext(ctx, "ping", key, fmt.Sprint(count), ip).CombinedOutput()
	if err != nil {
		// log.ERROR(err.Error())
		return false, nil
	}
	if isWindows() && !strings.Contains(string(output), "TTL") {
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

func ARP() map[string]string {
	res := make(map[string]string)
	if !isWindows() {
		// output, err := exec.Command("arp", "-a").CombinedOutput()
		// if err != nil {
		// 	log.ERROR(err.Error())
		// 	return res
		// }

		f, err := os.Open("/app/data/arp.log")
		if err != nil {
			log.ERROR(err.Error())
			return res
		}

		output, err := io.ReadAll(f)
		if err != nil {
			log.ERROR(err.Error())
			return res
		}

		// output = []byte(`? (192.168.0.21) at <incomplete> on enp3s0
		// ? (192.168.0.4) at <incomplete> on enp3s0
		// ? (192.168.0.78) at <incomplete> on enp3s0
		// ? (192.168.0.81) at 8c:ce:4e:14:96:25 [ether] on enp3s0
		// ? (192.168.0.135) at 04:10:6b:51:ac:ec [ether] on enp3s0
		// ? (192.168.0.11) at <incomplete> on enp3s0
		// ? (192.168.0.13) at <incomplete> on enp3s0
		// ? (172.22.0.2) at 02:42:ac:16:00:02 [ether] on br-43c201b3215e
		// RT (192.168.0.1) at 90:16:ba:74:a9:79 [ether] on enp3s0
		// ? (192.168.0.20) at 2c:f0:5d:2e:73:9e [ether] on enp3s0
		// ? (172.25.0.5) at 02:42:ac:19:00:05 [ether] on br-8491f428c062
		// ? (192.168.0.105) at dc:97:58:2d:64:50 [ether] on enp3s0
		// ? (192.168.0.8) at <incomplete> on enp3s0
		// ? (192.168.0.10) at d8:43:ae:25:52:0d [ether] on enp3s0
		// ? (192.168.0.40) at b8:4d:43:9c:1f:56 [ether] on enp3s0
		// ? (192.168.0.69) at <incomplete> on enp3s0
		// ? (192.168.0.72) at <incomplete> on enp3s0`)

		strOutput := strings.Split(string(output), "\n")
		for _, s := range strOutput {
			re := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*?(\w{2}:\w{2}:\w{2}:\w{2}:\w{2}:\w{2})`)
			matches := re.FindStringSubmatch(s)

			if len(matches) >= 3 {
				ip := matches[1]
				mac := matches[2]
				res[ip] = mac
			}
		}
	}

	return res
}
