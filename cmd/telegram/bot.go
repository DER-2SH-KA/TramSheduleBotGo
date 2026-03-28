package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DER-2SH-KA/TramSheduleBotGo/cmd/parser"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	Token             string
	LongPollingTiming time.Duration
	Parser            *parser.Parser
}

func (b *Bot) Start() error {
	bot, err := telego.NewBot(
		b.Token,
		telego.WithDefaultDebugLogger(),
	)
	if err != nil {
		return fmt.Errorf("Ошибка при создании бота: %s\n", err)
	}

	ctxLongPolling := context.Background()

	updates, err := bot.UpdatesViaLongPolling(ctxLongPolling, nil)
	if err != nil {
		return fmt.Errorf("Ошибка создания канала обновлений по принципу Long Polling: %s\n", err)
	}

	handler, err := th.NewBotHandler(bot, updates)
	if err != nil {
		return fmt.Errorf("Ошибка создания обработчика обновлений: %s\n", err)
	}

	handler.Handle(handleStartCommand, th.CommandEqual("start"))           // Обработка команды /start.
	handler.Handle(b.handleParseScheduleMessages, th.AnyMessageWithText()) // Обработка текстовых сообщений.
	handler.Handle(handleUnknownMessage, th.AnyMessage())                  // Обработка неизвестных сообщений.

	log.Println("Бот запущен...")

	defer func(handler *th.BotHandler) {
		err := handler.Stop()
		if err != nil {
			log.Println("Ошибка завершения работы обработчика сообщений", err)
			os.Exit(1)
		}
	}(handler)

	err = handler.Start()
	if err != nil {
		return fmt.Errorf("Ошибка запуска обработчика сообщений", err)
	}

	return nil
}

func createKeyboard() *telego.ReplyKeyboardMarkup {
	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("2 ЧС-КМ Б"),
			tu.KeyboardButton("2 ЧС-КМ В"),
			tu.KeyboardButton("8 ЧС-Т Б"),
			tu.KeyboardButton("8 ЧС-Т В"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("2 КМ-ЧС Б"),
			tu.KeyboardButton("2 КМ-ЧС В"),
			tu.KeyboardButton("8 Т-ЧС Б"),
			tu.KeyboardButton("8 Т-ЧС В"),
		),
	).WithResizeKeyboard()

	return keyboard
}

func handleStartCommand(ctx *th.Context, update telego.Update) error {
	message, err := ctx.Bot().SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:      update.Message.Chat.ChatID(),
			Text:        "Добро пожаловать! Выберите маршрут для получения расписания на клавиатуре.",
			ReplyMarkup: createKeyboard(),
		},
	)
	if err != nil {
		return err
	}

	if message != nil {
		log.Println("Message was send:", message)
	}

	return nil
}

func (b *Bot) handleParseScheduleMessages(ctx *th.Context, update telego.Update) error {
	messageText := update.Message.Text
	messageReplyText := ""

	switch messageText {
	case "2 ЧС-КМ Б":
		schedule, err := b.Parser.Parse(parser.ROUTE_2, parser.WORK_DAYS, parser.DIRECTION_FROM_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "2 ЧС-КМ В":
		schedule, err := b.Parser.Parse(parser.ROUTE_2, parser.WEEKENDS, parser.DIRECTION_FROM_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "8 ЧС-Т Б":
		schedule, err := b.Parser.Parse(parser.ROUTE_8, parser.WORK_DAYS, parser.DIRECTION_FROM_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "8 ЧС-Т В":
		schedule, err := b.Parser.Parse(parser.ROUTE_8, parser.WEEKENDS, parser.DIRECTION_FROM_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "2 КМ-ЧС Б":
		schedule, err := b.Parser.Parse(parser.ROUTE_2, parser.WORK_DAYS, parser.DIRECTION_TO_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "2 КМ-ЧС В":
		schedule, err := b.Parser.Parse(parser.ROUTE_2, parser.WEEKENDS, parser.DIRECTION_TO_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "8 Т-ЧС Б":
		schedule, err := b.Parser.Parse(parser.ROUTE_8, parser.WORK_DAYS, parser.DIRECTION_TO_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	case "8 Т-ЧС В":
		schedule, err := b.Parser.Parse(parser.ROUTE_8, parser.WEEKENDS, parser.DIRECTION_TO_CH_SLOBODA)
		if err != nil {
			return err
		}
		messageReplyText = schedule
	default:
		messageReplyText = "Неизвестный маршрут"
	}

	message, err := ctx.Bot().SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:      update.Message.Chat.ChatID(),
			Text:        messageReplyText,
			ReplyMarkup: createKeyboard(),
		},
	)
	if err != nil {
		return err
	}

	if message != nil {
		log.Println("Message was send:", message)
	}

	return nil
}

func handleUnknownMessage(ctx *th.Context, update telego.Update) error {
	message, err := ctx.Bot().SendMessage(
		ctx,
		&telego.SendMessageParams{
			ChatID:      update.Message.Chat.ChatID(),
			Text:        "Извините, я вас не понимаю. Проверьте запрос на возможные ошибки и повторите попытку",
			ReplyMarkup: createKeyboard(),
		},
	)
	if err != nil {
		return err
	}

	if message != nil {
		log.Println("Message was send:", message)
	}

	return nil
}
