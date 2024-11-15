package handlers

import (
	"bot-cf-simple/internal/texts"
	"strconv"
	"strings"

	tg "github.com/OvyFlash/telegram-bot-api"
)

func Start(chatID int64, bot *tg.BotAPI) {
  msg := tg.NewMessage(chatID, texts.Start)
  bot.Send(msg)
}

func TakeTxt(chatID, adminsChatID, usrID int64, msgText, usrName string, bot *tg.BotAPI) {
  take := strings.Split(msgText, " ")
  usrIDstr := strconv.Itoa(int(usrID))

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
//message to admins about user id
    msg2 := tg.NewMessage(adminsChatID, "ID: " + "`" + usrIDstr + "`")
    msg2.ParseMode = "MarkdownV2"
    //message to user which mean that all successful
    msg3 := tg.NewMessage(chatID, "тейк был успешно отправлен")
    
    //send messages
    bot.Send(msg)
    bot.Send(msg2)
    bot.Send(msg3)
    return
  } else if neanon == true {
    //message from user
    msg := tg.NewMessage(adminsChatID, msgText + "\n\n#тейк")
    bot.Send(msg)
    //message to admins about user id
    msg2 := tg.NewMessage(adminsChatID, "тейк от: " + "@" + usrName + ", " + "ID: " + "`" + usrIDstr + "`")
    msg2.ParseMode = "MarkdownV2"
    bot.Send(msg2)
    //message to user which mean that all successful
    msg3 := tg.NewMessage(chatID, "тейк был успешно отправлен")
    bot.Send(msg3)
    return
  } else {
    msg := tg.NewMessage(chatID, texts.Error)
    bot.Send(msg)
  }
}

func Files(chatID, adminsChatID, usrID, msgID int64, bot *tg.BotAPI, usrName string) {
  usrIDstr := strconv.Itoa(int(usrID))

  //files
  msg := tg.NewForward(adminsChatID, chatID, int(msgID))
  bot.Send(msg)
  //id
  msg3 := tg.NewMessage(adminsChatID, "ID: " + usrIDstr)
  bot.Send(msg3)
  //user must send a text
  msg2 := tg.NewMessage(chatID, "файлы отправлены, теперь напишите свой текст, как обычный тейк")
  bot.Send(msg2)
}

func Voices(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVoice tg.Voice) {
  usrIDstr := strconv.Itoa(int(usrID))

  fileID := msgVoice.FileID

  //voices 
  msg := tg.NewVoice(adminsChatID, tg.FileID(fileID))
  bot.Send(msg)
  //id 
  msg2 := tg.NewMessage(adminsChatID, "тейк от: " + "@" + usrName + ", " + "ID: " + "`" + usrIDstr + "`")
  msg2.ParseMode = "MarkdownV2"
  bot.Send(msg2)
  //user message
  msg3 := tg.NewMessage(chatID, "голосовое отправлено успешно")
  bot.Send(msg3)
}

func VideoNotes(chatID, adminsChatID, usrID int64, bot *tg.BotAPI, usrName string, msgVideoNote tg.VideoNote) {
  usrIDstr := strconv.Itoa(int(usrID))

  fileID := msgVideoNote.FileID
  duration := msgVideoNote.Duration

  //voices 
  msg := tg.NewVideoNote(adminsChatID, duration, tg.FileID(fileID))
  bot.Send(msg)
  //id 
  msg2 := tg.NewMessage(adminsChatID, "тейк от: " + "@" + usrName + ", " + "ID: " + "`" + usrIDstr + "`")
  msg2.ParseMode = "MarkdownV2"
  bot.Send(msg2)
  //user message
  msg3 := tg.NewMessage(chatID, "кружок отправлен успешно")
  bot.Send(msg3)
}
