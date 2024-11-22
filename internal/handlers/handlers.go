package handlers

import (
	"bot-cf-simple/internal/db"
	"bot-cf-simple/internal/texts"
	"log"
	"strings"

	tg "github.com/OvyFlash/telegram-bot-api"
)

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
      log.Println("cannot send message", err)
    }
    //add to db
    db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

    bot.Send(msg2)
    return
  } else if neanon == true {
    //message from user
    msg := tg.NewMessage(adminsChatID, msgText + "\n\n#тейк")
    //message to user which mean that all successful
    msg2 := tg.NewMessage(chatID, "тейк был успешно отправлен")

    //send messages
    msgID, err :=bot.Send(msg)
    if err != nil {
      log.Println("cannot send message", err)
    }

    //add to db
    db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

    bot.Send(msg2)
    return
  } else {
    msg := tg.NewMessage(chatID, texts.Error)
    bot.Send(msg)
  }
}

func Files(chatID, adminsChatID, msgID, usrID int64, bot *tg.BotAPI, usrName string) {
  //files
  msg := tg.NewForward(adminsChatID, chatID, int(msgID))
  //user must send a text
  msg2 := tg.NewMessage(chatID, "файлы отправлены, теперь напишите свой текст, как обычный тейк")

  //send messages
  msgIDbot, err := bot.Send(msg)
  if err != nil {
    log.Println("cannot send message", err)
  }

  //add to db
  db.Add(int64(msgIDbot.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
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
    log.Println("cannot send message", err)
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
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
    log.Println("cannot send message", err)
  }

  //add to db
  db.Add(int64(msgID.MessageID), chatID, usrID, usrName)

  bot.Send(msg2)
}

func Reply(bot *tg.BotAPI, msgText string, msgID, adminsChatID int64) {
  //get chat id from db
  chatID, err := db.GetChatIDByMsgID(msgID)
  if err != nil {
    log.Fatal("cannot get from table", err)
  }

  //send reply 
  msg := tg.NewMessage(chatID, "ответ от админа:\n\n" + msgText)
  //send what all successful
  msg2 := tg.NewMessage(adminsChatID, "ответ успешно отправлен")

  bot.Send(msg)
  bot.Send(msg2)
  log.Println("ответ успешно отправлен")
}

func Ban(msgID int64) {
  //get chat id from db
  usrID, err := db.GetUsrIDByMsgID(msgID)
  if err != nil {
    log.Fatal("cannot get from table", err)
  }

  db.Ban(usrID)
}
