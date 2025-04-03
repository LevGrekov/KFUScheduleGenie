package mybot

import (
	"context"
	"log"

	"github.com/LevGrekov/KFUScheduleGenie/parser"
	"github.com/LevGrekov/KFUScheduleGenie/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	firstName, lastName := utils.TranscribeName(update.Message.Text)
	if firstName == "" || lastName == "" {
		sendSafeMessage(ctx, b, update.Message.Chat.ID, "Пожалуйста, введите Фамилию и Имя через пробел\nПример: Юрий Агачев")
		return
	}

	scheduleHTML, err := parser.ParseSchedule(firstName, lastName)
	if err != nil {
		log.Printf("Parse error: %v", err)
		log.Printf("Parse error: %v %v", firstName, lastName)
		sendSafeMessage(ctx, b, update.Message.Chat.ID, "Ошибка получения расписания")
	}
	sendSafeMessage(ctx, b, update.Message.Chat.ID, scheduleHTML)
}

func sendSafeMessage(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		log.Printf("Send message failed: %v", err)
	}
}
