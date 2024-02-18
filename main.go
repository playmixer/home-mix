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
	log = logger.New("log")
	wg  = sync.WaitGroup{}
	bot *tgbotapi.BotAPI
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.ERROR("Error loading .env file")
	}

	database.Init()
	tools.Init()

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.ERROR(err.Error())
		return
	}
}

func main() {
	log.INFO("Starting ...")
	ctx, cancel := context.WithCancel(context.Background())
	go handlePing(ctx)
	go proccess(ctx)
	go startHttp()

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
