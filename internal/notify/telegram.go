package notify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

var (
	chatIDs  []int64
	botToken string
	bot      *tgbotapi.BotAPI
)

func InitTelegram(c []int64, b string, botDebug bool) {
	log.Info("Initializing Telegram Notifications Module")

	var err error
	chatIDs = c
	botToken = b

	// Create a new bot instance
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Enable debugging
	bot.Debug = botDebug

	log.Info("Telegram: Authorized on account ", bot.Self.UserName)
	for _, chatID := range chatIDs {
		msg := tgbotapi.NewMessage(chatID, "TestBot now online")
		retMsg, err := bot.Send(msg)

		if err != nil {
			log.Error("Error sending starting Telegram Msg... ", retMsg)
		}
	}

	go updateHandlerTelegram()
}

func updateHandlerTelegram() {
	// Create a new update configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Start receiving updates
	updates := bot.GetUpdatesChan(u)

	// Process updates
	for update := range updates {
		if update.Message != nil { // Check if we've received a message
			clientFound := false
			for _, c := range chatIDs {
				if c == update.Message.Chat.ID {
					log.Info(update.Message.From.UserName, update.Message.Text)
					clientFound = true

					// Reply to the message
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, I am your friendly DockerRight TelegramBot :)\nSadly I am currently not able to do something useful with your messages :(")
					bot.Send(msg)
				}
			}
			if !clientFound {
				log.Warn("Telegram: Received Msg from Unknown Client ", update.Message.Chat.ID, " ", update.Message.From.UserName)
			}
		}
	}
}

func NotifierTelegram(msg string) {
	for _, chatID := range chatIDs {
		sendMsg := tgbotapi.NewMessage(chatID, msg)
		bot.Send(sendMsg)
	}
}
