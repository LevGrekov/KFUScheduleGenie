package mybot

import (
	"context"
	"encoding/binary"
	"log"

	"github.com/LevGrekov/KFUScheduleGenie/kfuapi"
	"github.com/LevGrekov/KFUScheduleGenie/newparser"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

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
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    mes.Message.Chat.ID,
		MessageID: mes.Message.ID,
	})
}
