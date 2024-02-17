package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/playmixer/corvid/logger"
	"github.com/playmixer/home-mix/database"
	"github.com/playmixer/home-mix/tools"
)

type ChangeIp struct {
	IP     string
	Name   string
	Online bool
}

var (
	log          = logger.New("log")
	logPing      = logger.New("ping")
	wg           = sync.WaitGroup{}
	chanIpStatus = make(chan ChangeIp, 10)
	bot          *tgbotapi.BotAPI
)

func changeIpStatus(ip, name string, online bool) {
	select {
	case chanIpStatus <- ChangeIp{IP: ip, Name: name, Online: online}:
	default:
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.ERROR("Error loading .env file")
	}

	database.Init()

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.ERROR(err.Error())
		return
	}
}

func main() {
	log.INFO("Starting ...")
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go handlePing(ctx)
	go proccess(ctx)

	done := make(chan os.Signal, 10)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done
	log.INFO("Stopping ...")
	cancel()
	time.Sleep(time.Second)
	wg.Wait()
	log.INFO("Stop app")
}

func proccess(ctx context.Context) {
	_chatId := os.Getenv("TELEGRAM_CHAT_ID")
	chatId, err := strconv.Atoi(_chatId)
	if err != nil {
		log.ERROR("error: not valid telegram chatId")
		panic(err)
	}
proccessLoop:
	for {
		select {
		case <-ctx.Done():
			break proccessLoop
		case ipStatus := <-chanIpStatus:
			message := fmt.Sprintf("%s %s offline", ipStatus.Name, ipStatus.IP)
			if ipStatus.Online {
				message = fmt.Sprintf("%s %s online", ipStatus.Name, ipStatus.IP)
			}

			_, err := bot.Send(tgbotapi.NewMessage(int64(chatId), fmt.Sprintf("status ip: %s", message)))
			if err != nil {
				log.ERROR(err.Error())
			}
		}
	}
}

func handlePing(ctx context.Context) {
	defer wg.Done()
	t := tools.NewThread()
	t.SetMax(5)

	conn, err := database.Connect()
	if err != nil {
		log.ERROR(err.Error())
		panic(err)
	}

mainLoop:
	for {
		for i := 1; i <= 254; i++ {
			select {
			case <-ctx.Done():
				log.INFO("Stop ping")
				break mainLoop
			default:
				t.Wait()
				ip := fmt.Sprintf("192.168.0.%v", i)
				t.Add()
				go func(addres string) {
					defer t.Done()
					ok, _ := tools.Ping(addres)
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
					time.Sleep(time.Second)
				}(ip)
			}

		}
		t.WaitAll()
	}

}
