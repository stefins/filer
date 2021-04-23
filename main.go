package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	filer "github.com/stefins/filer/src"
)

func main() {
	db := filer.InitDB()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("T_API"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("Authorized on ", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for update := range updates {
		if update.Message.Document != nil {
			file := &filer.File{
				Name:    update.Message.Document.FileName,
				File_ID: update.Message.Document.FileID,
			}
			if file.Insert(db) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "File Already Exists! Thank you :)")
				bot.Send(msg)
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thank you for adding the file :)")
			bot.Send(msg)
			continue
		}
		if update.Message != nil {
			files := filer.Search(update.Message.Text, db)
			for _, file := range files {
				msg := tgbotapi.NewDocumentShare(update.Message.Chat.ID, file.File_ID)
				bot.Send(msg)
			}
			continue
		}
	}
}
