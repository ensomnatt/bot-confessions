package main

import (
	"bot-cf-simple/internal/db"
	"bot-cf-simple/internal/handlers"
	"bot-cf-simple/internal/initBot"
	"bot-cf-simple/internal/texts"
	"log"
	"os"
	"strings"

	tg "github.com/OvyFlash/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatal(err)
  }

  log.Println("[INFO]: load .env")

  connStr := os.Getenv("DB_CONN_STR")
  token := os.Getenv("TOKEN")
  adminsChatIDstr := os.Getenv("ADM_CHAT")

  //init bot and db
  adminsChatID, updates, bot := initbot.New(token, adminsChatIDstr) 
  db.Init(connStr)
  log.Println("[INFO]: init bot and db")
  defer db.Close()

  //start taking updates
  for u := range updates {
    if u.Message == nil {
      continue
    } else if db.CheckBan(u.Message.From.ID) {
      msg := tg.NewMessage(u.Message.Chat.ID, texts.Banned)
      bot.Send(msg)
      continue
    }

    chatID, msgText, usrName, msgPhoto, msgVideo, msgVoice, msgVideoNote, replyMsgId, usrID := initbot.CreateVars(u)

    //photos and videos
    if len(msgPhoto) > 0 || msgVideo != nil {
      handlers.Photos(chatID, adminsChatID, usrID, bot, usrName, msgPhoto)
    }

    //voices
    if msgVoice != nil {
      handlers.Voices(chatID, adminsChatID, usrID, bot, usrName, *msgVoice)
    }

    //video notes
    if msgVideoNote != nil {
      handlers.VideoNotes(chatID, adminsChatID, usrID, bot, usrName, *msgVideoNote)
    }

    //only text
    if msgText != "" {
      if chatID == adminsChatID && u.Message.ReplyToMessage != nil && u.Message.ReplyToMessage.From.ID == bot.Self.ID {
        var command string = "/reply"
        //ban  
        if strings.Contains(msgText, "/ban") {
          command = "/ban"
        }
        //unban 
        if strings.Contains(msgText, "/unban") {
          command = "/unban"
        }
        switch command {
        case "/ban":
          handlers.Ban(replyMsgId)
        case "/unban":
          handlers.UnBan(replyMsgId)
        case "/reply":
          handlers.Reply(bot, msgText, usrName, replyMsgId, adminsChatID)
        }
      } else {
        switch msgText {
        case "/start":
          handlers.Start(chatID, bot)
        default:
          if chatID != adminsChatID {
            handlers.TakeTxt(chatID, adminsChatID, usrID, msgText, usrName, bot)
          }
        }
      } 
    }
  }
}
