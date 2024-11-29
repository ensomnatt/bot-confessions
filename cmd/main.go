package main

import (
	"bot-cf-simple/internal/db"
	"bot-cf-simple/internal/handlers"
	initbot "bot-cf-simple/internal/initBot"
	"bot-cf-simple/internal/logger"
	"bot-cf-simple/internal/texts"
	"os"
	"strings"

	tg "github.com/OvyFlash/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
  //init logger
  logger.Init()
  
  //load envs
  err := godotenv.Load(".env")
  if err != nil {
    logger.Logger.Error("не удалось загрузить .env файл", "error", err)
  }

  logger.Logger.Debug("загрузка .env файла завершена", "file", ".env")

  connStr := os.Getenv("DB_CONN_STR")
  token := os.Getenv("TOKEN")
  adminsChatIDstr := os.Getenv("ADM_CHAT")

  //init bot and db
  adminsChatID, updates, bot := initbot.New(token, adminsChatIDstr) 
  db.Init(connStr)
  logger.Logger.Info("бот запущен", "token", token, "adminsChatID", adminsChatIDstr)
  defer db.Close()

  //start taking updates
  for u := range updates {
    if u.Message == nil {
      continue
    } else if db.CheckBan(u.Message.From.ID, u.Message.From.UserName) {
      msg := tg.NewMessage(u.Message.Chat.ID, texts.Banned)
      bot.Send(msg)
      continue
    }

    chatID, msgText, usrName, msgPhoto, msgVideo, msgVoice, msgVideoNote, replyMsgId, usrID := initbot.CreateVars(u)

    //photos and videos
    if len(msgPhoto) > 0 {
      handlers.Photos(chatID, adminsChatID, usrID, bot, usrName, msgPhoto)
    }

    if msgVideo != nil {
      handlers.Videos(chatID, adminsChatID, usrID, bot, usrName, *msgVideo)
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
      //commands
      if chatID == adminsChatID { 
        var command string
        //if replys
        if u.Message.ReplyToMessage != nil && u.Message.ReplyToMessage.From.ID == bot.Self.ID {
          if strings.Contains(msgText, "/ban") {
            command = "/ban"
          } else if strings.Contains(msgText, "/unban") {
            command = "/unban"             
          } else {
            command = "/reply"
          }
        } else {
          if strings.Contains(msgText, "/getusers") {
            command = "/getusers"
          } else if strings.Contains(msgText, "/getbans") {
            command = "/getbans"
          }
        } 

        //handle commands
        switch command {
        case "/ban":
          handlers.Ban(replyMsgId, usrName)
        case "/unban":
          handlers.UnBan(replyMsgId, usrName)
        case "/reply":
          handlers.Reply(bot, msgText, usrName, replyMsgId, adminsChatID)
        case "/getusers":
          handlers.GetUsers(bot, adminsChatID)
        case "/getbans":
          handlers.GetBans(bot, adminsChatID)
        }
      } else {
        switch msgText {
        case "/start":
          handlers.Start(chatID, usrID, bot, usrName)
        default:
          if chatID != adminsChatID {
            handlers.TakeTxt(chatID, adminsChatID, usrID, msgText, usrName, bot)
          }
        }
      } 
    }
  }
}
