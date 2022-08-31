package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot is alive :)")
			bot.Send(msg)
			continue
		}
		if update.Message.Document != nil {
			file := &filer.File{
				Name:    update.Message.Document.FileName,
				File_ID: update.Message.Document.FileID,
			}
			if file.Insert(db) == 0 {
				log.Printf("file exists name: %s id: %s\n", file.Name, file.File_ID)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "File already exists! Thank you :)")
				bot.Send(msg)
				continue
			}
			log.Printf("Added new file name: %s id: %s\n", file.Name, file.File_ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thank you for adding the file :)")
			bot.Send(msg)
			continue
		}
		if update.Message != nil {
			if update.Message.Video != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Oops! We don't support video file yet!")
				bot.Send(msg)
				continue
			}
			files := filer.Search(update.Message.Text, db)
			log.Printf("Searching for %s in database\n", update.Message.Text)
			if len(files) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Oops! Unable to find a file with that name!")
				bot.Send(msg)
				continue
			}
			for _, file := range files {
				fle := tgbotapi.FileID(file.File_ID)
				msg := tgbotapi.NewDocument(update.Message.Chat.ID, fle)
				bot.Send(msg)
			}
			continue
		}
	}
}
