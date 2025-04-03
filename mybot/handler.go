package mybot

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/LevGrekov/KFUScheduleGenie/kfuapi"
	"github.com/LevGrekov/KFUScheduleGenie/newparser"
	"github.com/LevGrekov/KFUScheduleGenie/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
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
	case 0:
		{
			sendSafeMessage(ctx, b, update.Message.Chat.ID, "Сотрудники не найдены")
			return
		}
	default:
		{
			sendTeachersKeyboard(ctx, b, update.Message.Chat.ID, teachers)
		}
	}

}

func sendTeachersKeyboard(ctx context.Context, b *bot.Bot, chatID int64, teachers []kfuapi.Employee) {
	kb := inline.New(b)

	for _, teacher := range teachers {
		idBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(idBytes, uint64(teacher.ID))
		kb.Row().Button(
			teacher.GetFullName(),
			idBytes,
			onInlineKeyboardSelect,
		)
	}

	kb.Row().Button("Отмена", []byte("cancel"), onCancel)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите преподавателя:",
		ReplyMarkup: kb,
	})

	if err != nil {
		log.Printf("Ошибка при отправке клавиатуры: %v", err)
	}
}

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	teacherID := int(binary.LittleEndian.Uint64(data))
	schedule, err := newparser.ParseSchedule(teacherID)
	if err != nil {
		sendSafeMessage(ctx, b, mes.Message.Chat.ID, "Ошибка получения расписания")
		return
	}
	sendSafeMessage(ctx, b, mes.Message.Chat.ID, schedule)
}

func onCancel(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	// Удаляем сообщение с клавиатурой
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    mes.Message.Chat.ID,
		MessageID: mes.Message.ID,
	})
}

func sendSafeMessage(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("Send message failed: %v", err)
	}
}
