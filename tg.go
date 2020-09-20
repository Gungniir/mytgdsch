package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strings"
)

var bot *tgbotapi.BotAPI
var err error

func initBotApi() {
	bot, err = tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.RemoveWebhook()
	if err != nil {
		log.Fatal(err)
	}

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(botUrl + botPath + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + botPath + bot.Token)
	go startListeningHttp()

	go func() {
		for update := range updates {
			handleUpdate(&update)
		}
	}()
}

func startListeningHttp() {
	log.Printf("Listening %s:%s...", listen, botPort)
	err := http.ListenAndServe(listen+":"+botPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func sendToAll(s string) {
	chats := getEnabledChats()
	log.Printf("Start messsanging to %d user(s)\r\n", len(chats))

	for _, chat := range chats {
		_, err := bot.Send(tgbotapi.NewMessage(chat, s))
		if err != nil {
			log.Printf("[chatID: %d]: %s", chat, err)
			disableChat(chat)
		}
	}

	log.Printf("Success :)\r\n")
}

func handleUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if !update.Message.IsCommand() {
		msg.Text = "К сожалению, я не понимаю, что ты написал :("
		_, err := bot.Send(msg)

		if err != nil {
			log.Fatal(err)
		}
	}

	switch strings.ToLower(update.Message.Command()) {
	case "start":
		msg.Text = "Привеет! Данный бот напоминает, в какой кабинет нужно идти в какой сейчас будет урок :)\r\n" +
			"Для того, чтобы включить напоминания введи /enable\r\n" +
			"Чтобы увидеть весь список комманд, введи /help"
	case "help":
		msg.Text = "Я мало что умею, но вот список моих возможностей:\r\n" +
			"/enable – включить напомнианя\r\n" +
			"/disable – выключить напоминания\r\n" +
			"/getschedule – вывести расписание\r\n" +
			"/help – показать этот список"
	case "enable":
		if isChatEnabled(update.Message.Chat.ID) {
			msg.Text = "Уведомления уже включены :)\r\n" +
				"Чтобы их выключить введи /disable"
			break
		}
		enableChat(update.Message.Chat.ID)
		msg.Text = "Уведомления включены, ожидай :)\r\n" +
			"Чтобы их выключить введи /disable"
	case "disable":
		if !isChatEnabled(update.Message.Chat.ID) {
			msg.Text = "Уведомления и так выключены :)\r\n" +
				"Чтобы их включить введи /enable"
			break
		}
		disableChat(update.Message.Chat.ID)
		msg.Text = "Уведомления выключены :)\r\n" +
			"Чтобы их включить введи /enable"
	case "getschedule":
		msg.Text = GetSchedule().String()
		msg.ParseMode = "Markdown"
	default:
		msg.Text = "Моя твоя не понимать. Если не знаешь, что мне написать, введи /help"
	}
	_, err := bot.Send(msg)

	if err != nil {
		log.Fatal(err)
	}
}
