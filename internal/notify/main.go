package notify

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	notifyLevel    int
	enableTelegram bool
)

func Init(notifyLevelStr, telegramBotToken string, telegramChatIDs []int) {
	log.Info("Initializing Notify Module")

	switch strings.ToLower(notifyLevelStr) {
	case "debug":
		notifyLevel = 5
	case "info":
		notifyLevel = 4
	case "warn":
		notifyLevel = 3
	case "error":
		notifyLevel = 2
	case "fatal":
		notifyLevel = 1
	case "panic":
		notifyLevel = 1
	case "none":
		notifyLevel = -1
	default:
		notifyLevel = 4
	}

	log.Info("NotifyLevel set to ", notifyLevel)
	enableTelegram = false

	if telegramBotToken != "" {
		enableTelegram = true
		tgBotDebug := false

		tgChatIDs := make([]int64, len(telegramChatIDs))
		for i, v := range telegramChatIDs {
			tgChatIDs[i] = int64(v)
		}

		InitTelegram(tgChatIDs, telegramBotToken, tgBotDebug)
	}
}

func Notifier(logLevel int, err ...interface{}) {
	// TODO: Add FormatWrapper
	notifyMsg := fmt.Sprintf("%s", err)
	if notifyMsg[:2] == "[[" {
		notifyMsg = notifyMsg[2:]
	}
	if notifyMsg[len(notifyMsg)-2:] == "]]" {
		notifyMsg = notifyMsg[0 : len(notifyMsg)-2]
	}

	if notifyLevel != 0 && notifyLevel >= logLevel {
		if enableTelegram {
			NotifierTelegram(notifyMsg)
		}
	}
}
