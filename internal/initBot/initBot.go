package initbot

import (
	tg "github.com/OvyFlash/telegram-bot-api"

	"log"
	"strconv"
)

func New(token string, adminsChatIDstr string) (int64, tg.UpdatesChannel, *tg.BotAPI) {
  adminsChatIDi, err := strconv.Atoi(adminsChatIDstr)
  if err != nil {
    log.Fatal("[FATAL]: cannot convert adminsChatIDstr to int")
  }
  adminsChatID := int64(adminsChatIDi)

  bot, err := tg.NewBotAPI(token)
  if err != nil {
    log.Fatal("[FATAL]: cannot create bot")
  }

  updateConfig := tg.NewUpdate(0)
  updateConfig.Timeout = 60

  updates := bot.GetUpdatesChan(updateConfig)
  
  return adminsChatID, updates, bot
}

func CreateVars(u tg.Update) (chatID, msgID int64, msgText, usrName string, msgPhoto []tg.PhotoSize, msgVideo *tg.Video, msgVoice *tg.Voice, msgVideoNote *tg.VideoNote, replyToMsgID, usrID int64) {
  chatID = u.Message.Chat.ID
  msgText = u.Message.Text
  usrName = u.SentFrom().UserName
  msgPhoto = u.Message.Photo
  msgID = int64(u.Message.MessageID)
  msgVideo = u.Message.Video
  msgVoice = u.Message.Voice
  msgVideoNote = u.Message.VideoNote
  if u.Message.ReplyToMessage != nil {
    replyToMsgID = int64(u.Message.ReplyToMessage.MessageID)
  } 
  usrID = u.Message.From.ID

  return chatID, msgID, msgText, usrName, msgPhoto, msgVideo, msgVoice, msgVideoNote, replyToMsgID, usrID
}
