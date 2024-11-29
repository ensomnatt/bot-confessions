package handlers

import (
	"bot-cf-simple/internal/db"
	"bot-cf-simple/internal/logger"
	"bot-cf-simple/internal/texts"
	"fmt"
	"strings"

	tg "github.com/OvyFlash/telegram-bot-api"
)

func Start(chatID, usrID int64, bot *tg.BotAPI, userName string) {
  msg := tg.NewMessage(chatID, texts.Start)
  bot.Send(msg)

  logger.Logger.Info("пользователь запустил бота", "chatID", chatID, "usrID", usrID, "userName", userName)
}

func formatTake(msgText string) []string {
  x := strings.ReplaceAll(msgText, "\n", " ")
  take := strings.TrimSpace(x)

  logger.Logger.Debug("форматирован текст тейка", "take", take)
  return strings.Split(take, " ")
}

func TakeTxt(chatID, adminsChatID, usrID int64, msgText, usrName string, bot *tg.BotAPI) {
  take := formatTake(msgText)

  anon := false
  neanon := false
  //find anon or not anon
  for _, w := range take {
    if w == "анон" || w == "Анон" {
      anon = true
      break
    } else if w == "неанон" || w == "Неанон" {
      neanon = true
      break
    }
  }

  //check anon or not anon
  if anon == true {
    //message from user
    msg := tg.NewMessage(adminsChatID, msgText + "\n\n#тейк")
    //message to user which mean that all successful
    msg2 := tg.NewMessage(chatID, "тейк был успешно отправлен")
    
    //send messages
    msgID, err := bot.Send(msg)
    if err != nil {
      logger.Logger.Warn("не удалось отправить сообщение", "error", err, "usrID", usrID, "usrName", usrName, "msgText", msgText)
    }
    //add to db
    db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

    bot.Send(msg2)
    
    logger.Logger.Info("тейк был отправлен", "usrID", usrID, "usrName", usrName)
    return
  } else if neanon == true {
    //message from user
    msg := tg.NewMessage(adminsChatID, msgText + "\n\n#тейк")
    //message to user which mean that all successful
    msg2 := tg.NewMessage(chatID, "тейк был успешно отправлен")
    msg3 := tg.NewMessage(adminsChatID, "тейк от @" + usrName)

    //send messages
    msgID, err :=bot.Send(msg)
    if err != nil {
      logger.Logger.Warn("не удалось отправить сообщение", "error", err, "usrID", usrID, "usrName", usrName, "msgText", msgText)
    }

    //add to db
    db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

    bot.Send(msg2)
    bot.Send(msg3)
    
    logger.Logger.Info("тейк был отправлен", "usrID", usrID, "usrName", usrName)
    return
  } else {
    msg := tg.NewMessage(chatID, texts.Error)
    bot.Send(msg)
  }
}

func Photos(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgPhoto[] tg.PhotoSize) {
  logger.Logger.Debug("обработка фото", "usrID", usrID, "usrName", usrName, "msgPhoto", msgPhoto)
  largest := msgPhoto[len(msgPhoto)-1]
  fileID := largest.FileID

  //photos
  msg := tg.NewPhoto(adminsChatID, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "фото успешно отправлены, теперь напишите свой текст, как обычный тейк")

  //send messages 

  msgID, err := bot.Send(msg)
  if err != nil {
    logger.Logger.Warn("не удалось отправить сообщение", "error", err, "usrID", usrID, "usrName", usrName)
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)

  logger.Logger.Info("фото было отправлено", "usrID", usrID, "usrName", usrName)
}

func Videos(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVideo tg.Video) { 
  logger.Logger.Debug("обработка видео", "usrID", usrID, "usrName", usrName, "msgVideo", msgVideo)
  fileID := msgVideo.FileID

  //photos
  msg := tg.NewVideo(adminsChatID, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "видео успешно отправлены, теперь напишите свой текст, как обычный тейк")

  //send messages 

  msgID, err := bot.Send(msg)
  if err != nil {
    logger.Logger.Warn("не удалось отправить сообщение", "error", err, "usrID", usrID, "usrName", usrName)
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)

  logger.Logger.Info("видео было отправлено", "usrID", usrID, "usrName", usrName)
}

func Voices(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVoice tg.Voice) { 
  logger.Logger.Debug("обработка голосового", "usrID", usrID, "usrName", usrName, "msgVoice", msgVoice)
  fileID := msgVoice.FileID

  //voices 
  msg := tg.NewVoice(adminsChatID, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "голосовое отправлено успешно")
  //send messages 

  msgID, err := bot.Send(msg)
  if err != nil {
    logger.Logger.Warn("не удалось отправить сообщение", "error", err, "usrID", usrID, "usrName", usrName)
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
  logger.Logger.Info("голосовое было отправлено", "usrID", usrID, "usrName", usrName)
}

func VideoNotes(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVideoNote tg.VideoNote) { 
  logger.Logger.Debug("обработка видеосообщения", "usrID", usrID, "usrName", usrName, "msgVideoNote", msgVideoNote)
  fileID := msgVideoNote.FileID
  duration := msgVideoNote.Duration

  //voices 
  msg := tg.NewVideoNote(adminsChatID, duration, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "кружок отправлен успешно")

  //send messages 
  msgID, err := bot.Send(msg)
  if err != nil {
    logger.Logger.Warn("не удалось отправить сообщение", "error", err, "usrID", usrID, "usrName", usrName)
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
  logger.Logger.Info("кружок был отправлен", "usrID", usrID, "usrName", usrName)
}

func Reply(bot *tg.BotAPI, msgText, usrName string, msgID, adminsChatID int64) {
  //get chat id from db
  chatID := db.GetChatIDByMsgID(msgID)

  //send reply 
  msg := tg.NewMessage(chatID, "ответ от админа:\n\n" + msgText)
  //send what all successful
  msg2 := tg.NewMessage(adminsChatID, "ответ успешно отправлен")

  bot.Send(msg)
  bot.Send(msg2)
  logger.Logger.Info("ответ был отправлен", "chatID", chatID, "usrName", usrName)
}

func Ban(msgID int64, userName string) {
  //get chat id from db
  usrID := db.GetUsrIDByMsgID(msgID)

  db.Ban(usrID, userName)
}

func UnBan(msgID int64, userName string) {
  //get chat id from db
  usrID := db.GetUsrIDByMsgID(msgID)

  db.UnBan(usrID, userName)
}

func GetUsers(bot *tg.BotAPI, adminsChatID int64) {
  users := db.GetUsers()

  var usersInfo []string
  for _, v := range users {
    bannedStr := "нет"
    if v.Banned {
      bannedStr = "да" 
    }

    usersInfo = append(usersInfo, fmt.Sprintf("ID: %d,\nимя: @%s,\nзабанен: %s", v.ID, v.UserName, bannedStr))
  }

  msg := tg.NewMessage(adminsChatID, strings.Join(usersInfo, "\n\n"))
  bot.Send(msg)

  logger.Logger.Info("список пользователей был отправлен админам", "list", users)
}

func GetBans(bot *tg.BotAPI, adminsChatID int64) {
  banned_users := db.GetBans()

  var banned_users_info []string
  for _, v := range banned_users {
    banned_users_info = append(banned_users_info, fmt.Sprintf("ID: %d,\nимя: @%s", v.ID, v.UserName))
  }

  msg := tg.NewMessage(adminsChatID, strings.Join(banned_users_info, "   ") + "   \n\n")
  bot.Send(msg)

  logger.Logger.Info("список забаненных пользователей был отправлен админам", "list", banned_users)
}
