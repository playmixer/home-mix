package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/playmixer/corvid/logger"
	"github.com/playmixer/home-mix/database"
	"github.com/playmixer/home-mix/tools"
)

var (
	logPing      = logger.New("ping")
	chanIpStatus = make(chan ChangeIp, 10)
)

func changeIpStatus(ip, name string, online bool) {
	select {
	case chanIpStatus <- ChangeIp{IP: ip, Name: name, Online: online}:
	default:
	}
}

func handlePing(ctx context.Context) {
	wg.Add(1)
	defer wg.Done()
	t := tools.NewThread()
	eMax := tools.Getenv("PING_THREAD", "2")
	max, err := strconv.Atoi(eMax)
	if err != nil {
		log.ERROR(err.Error())
		max = 2
	}
	t.SetMax(max)

	conn, err := database.Connect()
	if err != nil {
		log.ERROR(err.Error())
		panic(err)
	}

	tick1min := time.NewTicker(time.Minute * 60)
	updateApr()

mainLoop:
	for {
		log.INFO("Start ping pool")
		allDevice := []database.Ping{}
		conn.Find(&allDevice)
		for _, d := range allDevice {
			select {
			case <-ctx.Done():
				log.INFO("Stop ping")
				break mainLoop
			case <-tick1min.C:
				updateApr()
			default:
				t.Wait()
				ip := d.IP
				t.Add()
				go func(addres string) {
					defer t.Done()
					ok, _ := tools.Ping2(addres)
					logPing.INFO(fmt.Sprintf("ping %s %v", addres, ok))

					ping := database.Ping{IP: addres}
					conn.Where(&ping).First(&ping)

					if ping.Online != ok {
						ping.Online = ok
						changeIpStatus(ping.IP, ping.Name, ping.Online)
						tx := conn.Save(&ping)
						if tx.Error != nil {
							log.ERROR("error: cant update row", addres, err.Error())
							return
						}
					}
				}(ip)
			}

		}
		t.WaitAll()
		log.INFO("stop ping pool")
	}

}

func updateApr() {
	conn, err := database.Connect()
	if err != nil {
		log.ERROR(err.Error())
		panic(err)
	}

	log.INFO("upd arp data")
	addrs := tools.ARP()
	for k, v := range addrs {
		mac := database.Ping{Mac: v}
		conn.Where(&mac).First(&mac)
		if mac.ID > 0 {
			if mac.IP == k {
				continue
			} else {
				mac.Mac = ""
				conn.Save(&mac)
			}
		}

		ping := database.Ping{IP: k}
		conn.Where(&ping).First(&ping)
		if ping.ID > 0 {
			ping.Mac = v
			conn.Save(&ping)
		}
	}
}
