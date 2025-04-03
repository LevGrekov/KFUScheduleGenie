package mybot

import (
	"context"
	"fmt"

	"github.com/LevGrekov/KFUScheduleGenie/kfuapi"
	"github.com/LevGrekov/KFUScheduleGenie/newparser"
	"github.com/LevGrekov/KFUScheduleGenie/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	message_text := update.Message.Text
	if !utils.IsValidFIO(message_text) {
		sendSafeMessage(ctx, b, update.Message.Chat.ID, "Некоректные данные")
		return
	}
	teachers, err := kfuapi.SearchEmployees(message_text)
	if err != nil {
		sendSafeMessage(ctx, b, update.Message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	switch len(teachers) {
	default:
		sendTeachersKeyboard(ctx, b, update.Message.Chat.ID, teachers)
		return
	case 0:
		sendSafeMessage(ctx, b, update.Message.Chat.ID, "Преподаватели не найдены")
		return
	case 1:
		teacherID := teachers[0].ID
		schedule, err := newparser.ParseSchedule(teacherID)
		if err != nil {
			sendSafeMessage(ctx, b, update.Message.Chat.ID, "Ошибка получения расписания")
		}
		sendSafeMessage(ctx, b, update.Message.Chat.ID, schedule)
		return
	}
}

func sendSafeMessage(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	const (
		maxMessageLength = 4096
		maxChunkSize     = 4000
	)

	if len(text) <= maxChunkSize {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: "HTML",
		})
		return
	}

	// Разбиваем длинное сообщение на части
	chunks := utils.SplitMessage(text, maxChunkSize)
	for _, chunk := range chunks {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      chunk,
			ParseMode: "HTML",
		})
	}
}
