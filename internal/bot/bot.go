package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"tlgbs/internal/gateway"
	"tlgbs/internal/repository"
)

const (
	// Offset — смещение обновлений (апдейтов) Telegram.
	// Указывает, с какого ID начинать получение сообщений.
	// Значение 0 означает, что бот будет получать все новые обновления с текущего момента.
	Offset = 0

	// TimeOut — время ожидания (в секундах).
	// Telegram-сервер будет удерживать соединение открытым до указанного времени,
	// чтобы бот не делал слишком частые запросы.
	TimeOut = 30
)

type Bot struct {
	api    *tgbotapi.BotAPI
	repo   *repository.UserRepository
	gw     *gateway.Client
	logger *zap.Logger
}

func NewBot(token string, repo *repository.UserRepository, gw *gateway.Client, logger *zap.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create BotAPI %w", err)
	}

	return &Bot{
		api:    api,
		repo:   repo,
		gw:     gw,
		logger: logger,
	}, nil
}

func (b *Bot) BotRun(ctx context.Context) {
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

			chatID := upd.Message.Chat.ID
			user := upd.Message.From
			text := upd.Message.Text

			switch text {
			case "/start":
				if err := b.handleStart(chatID); err != nil {
					b.logger.Error("failed to handle /start", zap.Error(err))
					_ = b.SendAnswerToTelegram(chatID, "Ошибка при получении случайного фильма.")
				}

			case "/random_movie":
				if err := b.handleRandomMovie(ctx, chatID, strconv.FormatInt(user.ID, 10)); err != nil {
					b.logger.Error("failed to handle /random_movie", zap.Error(err))
					_ = b.SendAnswerToTelegram(chatID, "Ошибка при получении случайного фильма.")
				}

			case "/recommendations":
				if err := b.handleRecommendations(ctx, chatID, strconv.FormatInt(user.ID, 10)); err != nil {
					b.logger.Error("failed to handle /recommendations", zap.Error(err))
					_ = b.SendAnswerToTelegram(chatID, "Ошибка при получении рекомендаций.")
				}

			default:
				_ = b.SendAnswerToTelegram(chatID, "Неизвестная команда. Доступные команды: /start, /random_movie, /recommendations.")
			}
		}
	}
}

func (b *Bot) SendAnswerToTelegram(chatID int64, text string) error {
	_, err := b.api.Send(tgbotapi.NewMessage(chatID, text))
	return err
}
