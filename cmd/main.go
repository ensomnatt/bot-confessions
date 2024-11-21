package main

import (
	db "bot-cf-simple/internal/db"
	"bot-cf-simple/internal/handlers"
	initbot "bot-cf-simple/internal/initBot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatal(err)
  }

  //connStr := os.Getenv("DB_CONN_STR")
  //log.Println(connStr)

  token := os.Getenv("TOKEN")
  adminsChatIDstr := os.Getenv("ADM_CHAT")

  //init bot and db
  adminsChatID, updates, bot := initbot.New(token, adminsChatIDstr)
  db.Init()
  defer db.Close()

  //start take updates
  for u := range updates {
    if u.Message == nil {
      continue
    }

    chatID, msgID, msgText, usrName, msgPhoto, msgVideo, msgVoice, msgVideoNote, replyMsgId := initbot.CreateVars(u)

    //photos and videos
    if len(msgPhoto) > 0 || msgVideo != nil {
      logFiles(usrName)
      handlers.Files(chatID, adminsChatID, msgID, bot, usrName)
    }

    //voices
    if msgVoice != nil {
      logVoice(usrName)
      handlers.Voices(chatID, adminsChatID, bot, usrName, *msgVoice)
    }

    //video notes
    if msgVideoNote != nil {
      logVideoNote(usrName)
      handlers.VideoNotes(chatID, adminsChatID, bot, usrName, *msgVideoNote)
    }

    //only text
    if msgText != "" {
      if chatID == adminsChatID && u.Message.ReplyToMessage != nil && u.Message.ReplyToMessage.From.ID == bot.Self.ID {
        logReply(usrName, replyMsgId)
        handlers.Reply(bot, msgText, replyMsgId)
      } else {
        switch msgText {
        case "/start":
          logTxt(msgText, usrName)
          handlers.Start(chatID, bot)
        default:
          if chatID != adminsChatID {
            logTxt(msgText, usrName)
            handlers.TakeTxt(chatID, adminsChatID,msgText, usrName, bot)
          }
        }
      } 
    }
  }
}

func logTxt(text string, usrName string) {
  log.Printf("пользователь @%s написал: %s", usrName, text)
}

func logFiles(usrName string) {
  log.Printf("пользователь @%s отправил файлы", usrName)
}

func logVoice(usrName string) {
  log.Printf("пользователь @%s отправил голосовое", usrName)
}

func logVideoNote(usrName string) {
  log.Printf("пользователь @%s отправил кружок", usrName)
}

func logReply(usrName string, msgID int64) {
  log.Printf("админ @%s ответил на сообщение %v", usrName, msgID)
}
