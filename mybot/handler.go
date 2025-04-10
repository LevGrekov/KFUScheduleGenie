package mybot

import (
	"context"
	"log/slog"

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

	go executeTeacherSearchFlow(ctx, b, update)

	// message_text := update.Message.Text

	// msg, _ := b.SendMessage(ctx, &bot.SendMessageParams{
	// 	ChatID: update.Message.Chat.ID,
	// 	Text:   "⌛ Обрабатываю запрос...",
	// })

	// defer func() {
	// 	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
	// 		ChatID:    update.Message.Chat.ID,
	// 		MessageID: msg.ID,
	// 	})
	// }()

	// if !utils.IsValidFIO(message_text) {
	// 	sendSafeMessage(ctx, b, update.Message.Chat.ID, "Вы ввели Некоректные данные")
	// 	return
	// }
	// teachers, err := kfuapi.SearchEmployees(message_text)
	// if err != nil {
	// 	slog.Error("Ошибка SearchEmployees: %v", err)
	// 	sendSafeMessage(ctx, b, update.Message.Chat.ID, "Не получилось обратиться к серверу КФУ")
	// 	return
	// }

	// switch len(teachers) {
	// default:
	// 	sendTeachersKeyboard(ctx, b, update.Message.Chat.ID, teachers)
	// 	return
	// case 0:
	// 	sendSafeMessage(ctx, b, update.Message.Chat.ID, "Преподаватели не найдены")
	// 	return
	// case 1:
	// 	teacherID := teachers[0].ID
	// 	sendSchedule(teacherID, func(s string) { sendSafeMessage(ctx, b, update.Message.Chat.ID, s) })
	// 	return
	// }
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
	teachers, err := kfuapi.SearchEmployees(message_text)
	if err != nil {
		slog.Error("Ошибка SearchEmployees: %v", err)
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
	schedule, err := newparser.ParseSchedule(teacherID)
	if err != nil {
		slog.Error("Ошибка ParseSchedule: %v", err)
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
