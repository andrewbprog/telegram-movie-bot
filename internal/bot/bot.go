package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

const (
	Offset  = 0
	TimeOut = 30
)

type Bot struct {
	api    *tgbotapi.BotAPI
	logger *zap.Logger
}

func NewBot(token string, logger *zap.Logger) (*Bot, error) {
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		api:    b,
		logger: logger,
	}, nil
}

func (b *Bot) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(Offset)
	u.Timeout = TimeOut
	updates := b.api.GetUpdatesChan(u)

	b.logger.Info("telegram bot started")

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("telegram bot stopping")
			return
		case upd := <-updates:
			if upd.Message == nil {
				continue
			}
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Hello! I received: "+upd.Message.Text)
			if _, err := b.api.Send(msg); err != nil {
				b.logger.Error("failed to send message", zap.Error(err))
			}
		}
	}
}
