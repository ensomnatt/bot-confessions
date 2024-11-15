package main

import (
	initbot "bot-cf-simple/internal/initBot"
  "bot-cf-simple/internal/handlers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatal(err)
  }

  token := os.Getenv("TOKEN")
  adminsChatIDstr := os.Getenv("ADM_CHAT")

  adminsChatID, updates, bot := initbot.New(token, adminsChatIDstr)

  //start take updates
  for u := range updates {
    if u.Message == nil {
      continue
    }

    chatID, usrID, msgID, msgText, usrName, msgPhoto, msgVideo, msgVoice, msgVideoNote := initbot.CreateVars(u)

    //photos and videos
    if len(msgPhoto) > 0 || msgVideo != nil {
      logFiles(usrName)
      handlers.Files(chatID, adminsChatID, usrID, msgID, bot, usrName)
    }

    //voices
    if msgVoice != nil {
      logVoice(usrName)
      handlers.Voices(chatID, adminsChatID, usrID, bot, usrName, *msgVoice)
    }

    //video notes
    if msgVideoNote != nil {
      logVideoNote(usrName)
      handlers.VideoNotes(chatID, adminsChatID, usrID, bot, usrName, *msgVideoNote)
    }

    //only text
    if msgText != "" {
      switch msgText {
      case "/start":
        logTxt(msgText, usrName)
        handlers.Start(chatID, bot)
      default:
        logTxt(msgText, usrName)
        handlers.TakeTxt(chatID, adminsChatID, usrID, msgText, usrName, bot)
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
