package handlers

import (
	"bot-cf-simple/internal/db"
	"bot-cf-simple/internal/texts"
	"log"
	"strings"

	tg "github.com/OvyFlash/telegram-bot-api"
)

func logAnon(msgText, usrName string) {
  log.Println("[INFO]: user @", usrName, "send anon message: ", msgText)
}

func logUnAnon(msgText, usrName string) {
  log.Println("[INFO]: user @", usrName, "send not anon message: ", msgText)
}

func logStart(usrName string) {
  log.Println("[INFO]: user @", usrName, "press /start")
}

func logFiles(usrName string) {
  log.Println("[INFO]: user @", usrName, "send files")
}

func logVoice(usrName string) {
  log.Println("[INFO]: user @", usrName, "send voice")
}

func logVideoNote(usrName string) {
  log.Println("[INFO]: user @", usrName, "send video note")
}

func logBan() {
  log.Println("[INFO]: call db.Ban")
}

func logUnBan() {
  log.Println("[INFO]: call db.UnBan")
}

func logReply(usrName string) {
  log.Println("[INFO]: admin ", usrName, "send reply")
}

func Start(chatID int64, bot *tg.BotAPI) {
  msg := tg.NewMessage(chatID, texts.Start)
  bot.Send(msg)
}

func TakeTxt(chatID, adminsChatID, usrID int64, msgText, usrName string, bot *tg.BotAPI) {
  take := strings.Split(msgText, " ")

  anon := false
  neanon := false
  //find anon or not anon
  for _, w := range take {
    if w == "анон" {
      anon = true
      break
    } else if w == "неанон" {
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
      log.Println("[ERROR]: cannot send message")
    }
    //add to db
    db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

    bot.Send(msg2)
    logAnon(msgText, usrName)
    return
  } else if neanon == true {
    //message from user
    msg := tg.NewMessage(adminsChatID, msgText + "\n\n#тейк")
    //message to user which mean that all successful
    msg2 := tg.NewMessage(chatID, "тейк был успешно отправлен")

    //send messages
    msgID, err :=bot.Send(msg)
    if err != nil {
      log.Println("[ERROR]: cannot send message")
    }

    //add to db
    db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

    bot.Send(msg2)
    logUnAnon(msgText, usrName)
    return
  } else {
    msg := tg.NewMessage(chatID, texts.Error)
    bot.Send(msg)
  }
}

func Photos(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgPhoto[] tg.PhotoSize) {
  largest := msgPhoto[len(msgPhoto)-1]
  fileID := largest.FileID

  //photos
  msg := tg.NewPhoto(adminsChatID, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "фото успешно отправлены, теперь напишите свой текст, как обычный тейк")

  //send messages 

  msgID, err := bot.Send(msg)
  if err != nil {
    log.Println("[ERROR]: cannot send message")
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)

  logFiles(usrName)
}

func Videos(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVideo tg.Video) { 
  fileID := msgVideo.FileID

  //photos
  msg := tg.NewVideo(adminsChatID, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "видео успешно отправлены, теперь напишите свой текст, как обычный тейк")

  //send messages 

  msgID, err := bot.Send(msg)
  if err != nil {
    log.Println("[ERROR]: cannot send message")
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)

  logFiles(usrName)
}

func Voices(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVoice tg.Voice) { 
  fileID := msgVoice.FileID

  //voices 
  msg := tg.NewVoice(adminsChatID, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "голосовое отправлено успешно")
  //send messages 

  msgID, err := bot.Send(msg)
  if err != nil {
    log.Println("[ERROR]: cannot send message")
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
  logVoice(usrName)
}

func VideoNotes(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVideoNote tg.VideoNote) { 
  fileID := msgVideoNote.FileID
  duration := msgVideoNote.Duration

  //voices 
  msg := tg.NewVideoNote(adminsChatID, duration, tg.FileID(fileID))
  //user message
  msg2 := tg.NewMessage(chatID, "кружок отправлен успешно")

  //send messages 
  msgID, err := bot.Send(msg)
  if err != nil {
    log.Println("[ERROR]: cannot send message")
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
  logVideoNote(usrName)
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
  logReply(usrName)
}

func Ban(msgID int64) {
  //get chat id from db
  usrID := db.GetUsrIDByMsgID(msgID)

  db.Ban(usrID)
  logBan()
}

func UnBan(msgID int64) {
  //get chat id from db
  usrID := db.GetUsrIDByMsgID(msgID)

  db.UnBan(usrID)
  logUnBan()
}
