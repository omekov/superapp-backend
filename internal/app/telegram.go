package app

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBot struct {
	token     string
	channelID string
}

func NewTelegramBot(token string) telegramBot {
	return telegramBot{
		token: token,
	}
}

func (tb telegramBot) Init() error {
	bot, err := tgbotapi.NewBotAPI(tb.token)
	if err != nil {
		return err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Команда для открытия Web App
		if update.Message.Text == "/start" {
			userID := update.Message.From.ID

			// Проверяем подписку
			isSubscribed, err := tb.isUserSubscribed(bot, userID)
			if err != nil {
				log.Println(err)
				continue
			}

			if isSubscribed {
				// Открываем Web App
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы подписаны! Можете открыть Web App: <ссылка>")
				bot.Send(msg)
			} else {
				// Просим подписаться
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, подпишитесь на наш канал, чтобы продолжить: https://t.me/your_channel")
				bot.Send(msg)
			}
		}
	}
	return nil
}

func (tb telegramBot) isUserSubscribed(bot *tgbotapi.BotAPI, userID int) (bool, error) {

	// Используем метод getChatMember для проверки статуса пользователя в канале
	member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		SuperGroupUsername: tb.channelID,
		UserID:             userID,
	})
	if err != nil {
		return false, err
	}

	// Проверяем статус пользователя
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		return true, nil
	}

	return false, nil
}
