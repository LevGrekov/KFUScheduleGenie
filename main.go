package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/LevGrekov/KFUScheduleGenie/mybot"
	"github.com/go-telegram/bot"
)

const (
	BOT_TOKEN = "7503877925:AAHQvB1Rqb3ltBNvaKPaLik68qoUIWOeQ9A"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(mybot.Handler),
	}

	b, err := bot.New(BOT_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
