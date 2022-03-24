package services

import (
	"fmt"
	"ibuYemekBotu/models"
	"ibuYemekBotu/mongo"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

func TelegramHandler() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic("Telegram Bot Not Found: ", err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	lunchListToday := getLunchList("today")
	lunchListTomorrow := getLunchList("tomorrow")

	c := cron.New()
	c.AddFunc("30 03 * * *", func() {
		log.Printf("Updating lunch list")
		lunchListToday = getLunchList("today")
		lunchListTomorrow = getLunchList("tomorrow")
		log.Printf("Lunch list today: %s", lunchListToday)
		log.Printf("Lunch list tomorrow: %s", lunchListTomorrow)
	})

	c.AddFunc("00 06 * * *", func() {
		sendListSubscribers(lunchListToday, bot)
	})

	c.Start()

	helloMessage := "-İBU Yemek Listesi Botuna Hoş Geldiniz!\n" +
		"-Her sabah 9'da yemek listesini almak için abone olun.\n" +
		"-Abone olmak için /subscribe\n" +
		"-Abonelikten çıkmak için /unsubscribe\n" +
		"-Bugünün Listesini öğrenmek için /today\n" +
		"-Yarının Listesini öğrenmek için /tomorrow\n" +
		"-Kaynak Kod İçin /source\n" +
		"-Yardım almak için /help\n"

	for update := range updates {
		if update.Message != nil { // If we got a message

			if update.Message.Text == "/today" {
				if lunchListToday == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yemek listesi yok")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, lunchListToday)
					bot.Send(msg)
				}
			} else if update.Message.Text == "/tomorrow" {
				if lunchListTomorrow == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yemek listesi yok")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, lunchListTomorrow)
					bot.Send(msg)
				}
			} else if update.Message.Text == "/start" {

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
				bot.Send(msg)
			} else if update.Message.Text == "/help" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMessage)
				bot.Send(msg)
			} else if update.Message.Text == "/source" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://github.com/safderun/ibuYemekBotuGolang")
				bot.Send(msg)
			} else if update.Message.Text == "/subscribe" {
				if mongo.GetUser(update.Message.Chat.ID) == false {
					user := models.User{update.Message.Chat.ID, update.Message.Chat.UserName, "", true}
					fmt.Println(user)
					mongo.Adduser(&user)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Abone oldunuz.")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Zaten abone oldunuz.")
					bot.Send(msg)
				}
			} else if update.Message.Text == "/unsubscribe" {
				if mongo.GetUser(update.Message.Chat.ID) == true {
					mongo.DeleteUser(update.Message.Chat.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Abonelikten çıktınız.")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Zaten abonelikten çıktınız.")
					bot.Send(msg)
				}
			} else if update.Message.Text == "/time" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Saat 11:30-14:00 ve 15:30-18:00 arası.")
				bot.Send(msg)
			} else if update.Message.Text == "/admin" {
				if update.Message.Chat.UserName == "safderun67" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin girişi yapıldı.")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin girişi yapılamadı.")
					bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bilinmeyen komut\nYardım Almak İçin /help")
				bot.Send(msg)
			}
		}
	}
}

func sendListSubscribers(lunchList string, bot *tgbotapi.BotAPI) {
	log.Println("Sending list to subscribers")
	if lunchList == "" {
		log.Println("Lunch list is empty")
	} else {
		userList := mongo.GetAllUsers()
		for _, user := range userList {
			msg := tgbotapi.NewMessage(user.ChatId, lunchList)
			bot.Send(msg)
		}
		log.Println("List sent")
	}
}

func getLunchList(day string) string {
	log.Println("Getting lunch list of " + day)
	lunch := scrapper(day)
	emptyLunch := models.Lunch{"", "", "", "", ""}
	if lunch == emptyLunch {
		return ""
	}
	lunchString := "Çorba: " + lunch.Corba + "\n" +
		"Ana Yemek: " + lunch.AnaYemek + "\n" +
		"İkinci Yemek: " + lunch.YardimciAnaYemek + "\n" +
		"Yan Yemek: " + lunch.YanYemek1 + "\n" +
		"Yan Yemek: " + lunch.YanYemek2

	if day == "today" {
		t := time.Now()
		lunchString = t.Format("02/01/2006") + "\n" + lunchString
		return lunchString
	} else {
		t := time.Now()
		t = t.AddDate(0, 0, 1)
		lunchString = t.Format("02/01/2006") + "\n" + lunchString
		return lunchString
	}
}
