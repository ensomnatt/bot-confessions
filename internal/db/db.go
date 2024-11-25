package db

import (
	"bot-cf-simple/internal/logger"
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init(connStr string) {
  //connect to database
  var err error
  db, err = sql.Open("postgres", connStr)
  if err != nil {
    logger.Logger.Error("не удалось подключиться к базе данных", "error", err, "connStr", connStr)
  }
  
  //create takes table
  query := `CREATE TABLE IF NOT EXISTS takes (
    id SERIAL PRIMARY KEY,
    msg_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL
  )`
  _, err = db.Exec(query)
  if err != nil {
    logger.Logger.Error("не удалось создать таблицу takes", "error", err)
  }

  logger.Logger.Debug("таблица takes создана")

  //create users table 
  query = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    user_name VARCHAR(50) NOT NULL,
    banned BOOLEAN DEFAULT FALSE
  )`
  _, err = db.Exec(query)
  if err != nil {
    logger.Logger.Error("не удалось создать таблицу users", "error", err)
  }

  logger.Logger.Debug("таблица users создана")

  err = db.Ping()
  if err != nil {
    logger.Logger.Error("не удалось произвести связь к базе данных", "error", err)
  } else {
    logger.Logger.Info("связь к базе данных установлена")
  } 
}

func Add(msgID, chatID, usrID int64, usrName string) {
  //add to takes
  query := `INSERT INTO takes (msg_id, chat_id, user_id) VALUES ($1, $2, $3)`
  _, err := db.Exec(query, msgID, chatID, usrID)
  if err != nil {
    logger.Logger.Warn("не удалось добавить в таблицу takes", "error", err, "msgID", msgID, "chatID", chatID, "usrID", usrID)
  }

  logger.Logger.Debug("добавлено в таблицу takes", "msgID", msgID, "chatID", chatID, "usrID", usrID)

  //add to users

  //check if user exists
  if checkIfUserExists(usrID) {
    return
  }

  query = `INSERT INTO users (user_id, user_name) VALUES ($1, $2)`
  _, err = db.Exec(query, usrID, usrName)
  if err != nil {
    logger.Logger.Warn("не удалось добавить в таблицу users", "error", err, "usrID", usrID, "usrName", usrName)
  }

  logger.Logger.Debug("добавлено в таблицу users", "usrID", usrID, "usrName", usrName)
}

func checkIfUserExists(usrID int64) bool {
  var dbUserID int64
  query := `SELECT user_id FROM users WHERE user_id = $1`
  _ = db.QueryRow(query, usrID).Scan(&dbUserID)

  if dbUserID == 0 {
    return false
  } else {
    return true
  }
}

func GetChatIDByMsgID(msgID int64) (int64) {
  var chatID int64
  query := `SELECT chat_id FROM takes WHERE msg_id = $1`
  err := db.QueryRow(query, msgID).Scan(&chatID)
  if err != nil {
    logger.Logger.Warn("не удалось получить chatID из таблицы takes", "error", err, "msgID", msgID, "chatID", chatID) 
  }

  logger.Logger.Debug("получен chatID из таблицы takes", "msgID", msgID, "chatID", chatID)
  return chatID
}

func GetUsrIDByMsgID(msgID int64) (int64) {
  var usrID int64
  query := `SELECT user_id FROM takes WHERE msg_id = $1`
  err := db.QueryRow(query, msgID).Scan(&usrID)
  if err != nil {
    logger.Logger.Warn("не удалось получить usrID из таблицы takes", "error", err, "msgID", msgID, "usrID", usrID)
  }

  logger.Logger.Debug("получен usrID из таблицы takes", "msgID", msgID, "usrID", usrID)
  return usrID
}

func GetUsrNameByUsrID(usrID int64) (string) {
  var usrName string
  query := `SELECT user_name FROM users WHERE user_id = $1`
  err := db.QueryRow(query, usrID).Scan(&usrName)
  if err != nil {
    logger.Logger.Warn("не удалось получить usrName из таблицы users", "error", err, "usrID", usrID)   
  } 

  logger.Logger.Debug("получен usrName из таблицы users", "usrID", usrID, "usrName", usrName)
  return usrName
} 

func Ban(usrID int64, userName string) {
  query := `UPDATE users SET banned = true WHERE user_id = $1`
  _, err := db.Exec(query, usrID)
  if err != nil {
    logger.Logger.Warn("не удалось забанить пользователя", "error", err, "usrID", usrID)  
  }

  logger.Logger.Info("забанен пользователь", "usrID", usrID, "usrName", userName)
}

func UnBan(usrID int64, userName string) {
  query := `UPDATE users SET banned = false WHERE user_id = $1`
  _, err := db.Exec(query, usrID)
  if err != nil {
    logger.Logger.Warn("не удалось разбанить пользователя", "error", err, "usrID", usrID, "usrName", userName)
  }

  logger.Logger.Info("разбанен пользователь", "usrID", usrID, "usrName", userName)
}

func CheckBan(usrID int64, userName string) bool {
  var banned bool
  query := `SELECT banned FROM users WHERE user_id = $1`
  _ = db.QueryRow(query, usrID).Scan(&banned)

  logger.Logger.Debug("получен banned из таблицы users", "usrID", usrID, "banned", banned, "usrName", userName)
  return banned
}

func Close() {
  if db != nil {
    db.Close()
  }

  logger.Logger.Debug("база данных закрыта")
}
