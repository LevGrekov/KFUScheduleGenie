package mybot

import (
	"context"
	"log/slog"

	"github.com/LevGrekov/KFUScheduleGenie/kfuapi"
	"github.com/LevGrekov/KFUScheduleGenie/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var client = kfuapi.NewClient()

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	executeTeacherSearchFlow(ctx, b, update)
}

func executeTeacherSearchFlow(ctx context.Context, b *bot.Bot, update *models.Update) {

	sendResult := func(text string) {
		sendSafeMessage(ctx, b, update.Message.Chat.ID, text)
	}
	message_text := update.Message.Text

	msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "⌛ Обрабатываю запрос...",
	})

	defer func() {
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: msg.ID,
		})
	}()

	if !utils.IsValidFIO(message_text) {
		sendResult("Вы ввели Некоректные данные")
		return
	}
	teachers, err := client.SearchEmployees(message_text)
	if err != nil {
		slog.Error("Ошибка SearchEmployees", "error", err)
		sendResult("Не получилось обратиться к серверу КФУ")
		return
	}

	switch len(teachers) {
	default:
		sendTeachersKeyboard(ctx, b, update.Message.Chat.ID, teachers)
		return
	case 0:
		sendResult("Преподаватели не найдены")
		return
	case 1:
		teacherID := teachers[0].ID
		sendSchedule(teacherID, func(s string) { sendResult(s) })
		return
	}
}

func sendSchedule(teacherID int, onComplete func(string)) {
	schedule, err := client.GetSchedule(teacherID)
	if err != nil {
		slog.Error("Ошибка Получения Расписания", "error", err)
		onComplete("Ошибка получения расписания")
	}
	onComplete(schedule)
}

func sendSafeMessage(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	const maxChunkSize = 4000

	if len(text) <= maxChunkSize {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: "HTML",
		})
		return
	}

	chunks := utils.SplitMessage(text, maxChunkSize)
	for _, chunk := range chunks {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      chunk,
			ParseMode: "HTML",
		})
	}
}
