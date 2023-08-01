package dbManager

import (
	"bot/botTool"
	. "bot/config"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// _ "github.com/go-sql-driver/mysql"
)

var timerule = regexp.MustCompile(`\d\d:\d\d (A|P)M UTC\+\d`)

type database struct {
	Db *sql.DB
}

func InitMysql(user, token, table string) (db *database) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, token, IPV4, table)
	dbv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return
	}
	err = dbv.Ping()
	if err != nil {
		log.Println(err)
	}
	dbv.SetMaxOpenConns(500)
	dbv.SetMaxIdleConns(50)
	dbv.SetConnMaxLifetime(1 * time.Hour)
	db = &database{dbv}
	return
}

func (db *database) CreateUserTable(userId string) {
	sqlStr := fmt.Sprintf("CREATE TABLE `%s` (userId CHAR(16) UNIQUE,times SMALLINT,name CHAR(80)) CHARSET=utf8mb4", userId)
	result, err := db.Db.Exec(sqlStr)
	if err != nil {
		log.Printf("%s when Exec Database in User", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in User", err)
	}
}

func (db *database) CreateGroupRules(chatId int64) {
	sqlStr := fmt.Sprintf("CREATE TABLE `%d` (regexpStr VARCHAR(2000),replyStr TEXT(20000)) CHARSET=utf8mb4", chatId)
	result, err := db.Db.Exec(sqlStr)
	if err != nil {
		log.Printf("%s when Exec Database in GroupRules", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in GroupRules", err)
	}
}

func (db *database) CreateUserConfig(userId string) {
	sqlStr := fmt.Sprintf("CREATE TABLE `%s` (userId CHAR(16) UNIQUE,username VARCHAR(200),name VARCHAR(2000),time datetime) CHARSET=utf8mb4", userId)
	result, err := db.Db.Exec(sqlStr)
	if err != nil {
		log.Printf("%s when Exec Database in User", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in User", err)
	}
}

func (db *database) CreateChatTable(chatId string) {
	sqlStr := fmt.Sprintf("CREATE TABLE `%s`(groupData CHAR(30) UNIQUE,times SMALLINT) CHARSET=utf8mb4", chatId)
	result, err := db.Db.Exec(sqlStr)
	if err != nil {
		// log.Printf("%s when Exec Database in Chat", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in Chat", err)
	}
}

func (db *database) TableInfo(groups *[]string) {
	sqlStr := `show tables`
	var data string
	rows, err := db.Db.Query(sqlStr)
	if err != nil {
		log.Println("Table info", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&data)
		for i, v := range data {
			if v == 'G' {
				*groups = append(*groups, data[:i])
				break
			}
		}
	}
}

func (db *database) AddMessage(chatId string, message string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("addMessage", err)
		}
	}()
	chatId = chatId + "Group"
	sqlStr := fmt.Sprintf("insert into `%s` (groupData,times) values(?,1) on DUPLICATE key update times=times+1", chatId)
	result, err := db.Db.Exec(sqlStr, message)
	if err != nil {
		driverErr, _ := err.(*mysql.MySQLError)
		if driverErr.Number == 1146 {
			db.CreateChatTable(chatId)
			result, err = db.Db.Exec(sqlStr, message)
			if err != nil {
				log.Println("Addmessage", err)
				return
			}
		} else {
			log.Println("Addmessage", message)
			return
		}
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in Group", err)
	}
}

func (db *database) AddRules(chatId int64, regexpStr string, replyStr string) {
	sqlStr := fmt.Sprintf("insert into `%d` (regexpStr,replyStr) values(?,?)", chatId)
	result, err := db.Db.Exec(sqlStr, regexpStr, replyStr)
	// log.Println(sqlStr)
	if err != nil {
		driverErr, _ := err.(*mysql.MySQLError)
		if driverErr.Number == 1146 {
			db.CreateGroupRules(chatId)
			result, err = db.Db.Exec(sqlStr, regexpStr, replyStr)
			if err != nil {
				log.Println(err, regexpStr, replyStr)
			}
		} else {
			log.Println("regexpStr,replyStr", regexpStr, replyStr)
		}
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in Group", err)
	}
}

func (db *database) AddUser(chatId string, userId string, name string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("AddUser", err)
		}
	}()
	chatId = chatId + "User"
	sqlStr := fmt.Sprintf("insert into `%s` (userId,times,name) values(?,1,?) on DUPLICATE key update times=times+1", chatId)
	result, err := db.Db.Exec(sqlStr, userId, name)
	// log.Println(sqlStr)
	if err != nil {
		driverErr, _ := err.(*mysql.MySQLError)
		if driverErr.Number == 1146 {
			db.CreateUserTable(chatId)
			result, err = db.Db.Exec(sqlStr, userId, name)
			if err != nil {
				log.Println(err, name)
			}
		} else {
			log.Println("name:", len(name))
		}
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in Group", err)
	}
}

func (db *database) AddGroup(update *tgbotapi.Update, message *tgbotapi.Message, chatId string, name string, groupname string, user string, username string, nickname string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("AddGroup", err)
		}
	}()
	var msg string
	sqlStr := "INSERT INTO `user`(`userid`,`username`,`name`,time) VALUES(?,?,?,NOW()) ON DUPLICATE KEY UPDATE `username`= ?,`name`=?,`time`=NOW()"
	result, _ := db.Db.Exec(sqlStr, user, username, nickname, username, nickname)
	_, err := result.RowsAffected()
	// log.Println(sqlStr)
	if err != nil {
		log.Println(err)
		log.Println(user)
		log.Println(username)
		err = nil
	}
	if message.Text == "" {
		sqlStr = "INSERT INTO `config`(`chatId`,`username`, `groupname`) VALUES(?,?,?) ON DUPLICATE KEY UPDATE `username`=?,`groupname`=?"
		result, err = db.Db.Exec(sqlStr, chatId, name, groupname, name, groupname)
		// log.Println(sqlStr)
		if err != nil {
			log.Println(err)
			log.Println(chatId)
			log.Println(len(name))
			log.Println(groupname)
			err = nil
		}
		_, err = result.RowsAffected()

		if err != nil {
			log.Printf("%s when RowsAffected in config", err)
		}
	}
	sqlStr = fmt.Sprintf("select `name`,`username` from `%s` where userid=?", chatId)
	row := db.Db.QueryRow(sqlStr, user)
	var nameDb, usernameDb string
	err = row.Scan(&nameDb, &usernameDb)
	var flag = true
	if err != nil {
		flag = false
		driverErr, _ := err.(*mysql.MySQLError)
		if err == sql.ErrNoRows {
		} else if driverErr.Number == 1146 {
			db.CreateUserConfig(chatId)
		} else {
			log.Println("user:", user)
		}
	}
	if flag {
		if usernameDb != username {
			msg = fmt.Sprintf("change username from @%s to @%s\n", usernameDb, username)
		}
		if nameDb != nickname {
			msg += fmt.Sprintf("change nickname from %s to %s\n", nameDb, nickname)
		}
		if msg == "" {
			sqlStr = fmt.Sprintf("UPDATE `%s` SET time = NOW() WHERE userId = ?", chatId)
			result, _ = db.Db.Exec(sqlStr, user)
			_, err = result.RowsAffected()
			if err != nil {
				log.Printf("%s when RowsAffected in %s", err, chatId)
			}
			return
		}
		msg = fmt.Sprintf("User: [%s](tg://user?id=%s)\n\n%s", user, user, msg)
		if !timerule.MatchString(nickname) {
			botTool.SendMessage(message, msg, false, "Markdown")
		}
	}
	sqlStr = fmt.Sprintf("insert into `%s` (`userId`,`username`,`name`,`time`) values(?,?,?,Now()) ON DUPLICATE KEY UPDATE `username`=?,`name`=?,`time`=Now()", chatId)
	result, _ = db.Db.Exec(sqlStr, user, username, nickname, username, nickname)
	// log.Println(sqlStr)

	_, err = result.RowsAffected()
	if err != nil {
		log.Printf("%s when RowsAffected in Group", err)
	}
}

func (db *database) GetAllWords(chatId string) (result map[string]int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("GetAllWord", err)
		}
	}()
	strSql := fmt.Sprintf("select groupData,times from `%s` order by `times` desc limit 0,150", chatId)
	rows, err := db.Db.Query(strSql)
	if err != nil {
		driverErr, _ := err.(*mysql.MySQLError)
		if driverErr.Number != 1146 {
			log.Println("Getallword", err)
		}
		return
	}
	defer rows.Close()
	result = make(map[string]int)
	for rows.Next() {
		var data string
		var times int
		rows.Scan(&data, &times)
		result[data] = times
	}
	// log.Println(result)
	return
}

func (db *database) GetAllUsers(chatId string) (result [2][]string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("getallUser", err)
		}
	}()
	strSql := fmt.Sprintf("select times,name from `%sUser` order by times desc limit 5", chatId)
	rows, err := db.Db.Query(strSql)
	if err != nil {
		log.Println("GetallUser", err)
		return
	}
	defer rows.Close()
	// result = make([][]string,0)
	for rows.Next() {
		// var data string
		var times string
		var name string
		rows.Scan(&times, &name)
		// log.Println(data, times, name)
		result[0] = append(result[0], times)
		result[1] = append(result[1], name)
	}
	// log.Println(result)
	return
}

func (db *database) GetAllRules(chatId int64) (result [][2]string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("getallUser", err)
		}
	}()
	strSql := fmt.Sprintf("select * from `%d`", chatId)
	rows, err := db.Db.Query(strSql)
	if err != nil {
		driverErr, _ := err.(*mysql.MySQLError)
		if driverErr.Number == 1146 {
			return nil
		} else {
			log.Println("GetallRulers", err)
		}
		return
	}
	defer rows.Close()
	// result = make([][]string,0)
	for rows.Next() {
		// var data string
		var regexpStr string
		var replyStr string
		rows.Scan(&regexpStr, &replyStr)
		// log.Println(data, times, name)
		result = append(result, [2]string{regexpStr, replyStr})
	}
	// log.Println(result)
	return
}

func (db *database) CheckId2User(id string) (result [3]string) {
	sqlStr := "select `name`,`username`,`time` from `user` where userid=?"
	row := db.Db.QueryRow(sqlStr, id)
	var name string
	var username string
	var time string
	row.Scan(&name, &username, &time)
	result = [3]string{username, name, time}
	return
}

func (db *database) CheckUser2Id(username string) (result [3]string) {
	sqlStr := "select `userid`,`name`,`time` from `user` where username=?"
	row := db.Db.QueryRow(sqlStr, username)
	var id string
	var name string
	var time string
	row.Scan(&id, &name, &time)
	result = [3]string{id, name, time}
	return
}

func (db *database) Clear() {
	sqlStr := `show tables`
	var data string
	rows, err := db.Db.Query(sqlStr)
	if err != nil {
		log.Println("clear", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&data)
		strSql := fmt.Sprintf("DROP TABLE `%s`", data)
		_, err := db.Db.Exec(strSql)
		if err != nil {
			log.Println("Clear", err)
		}
	}
}

func (db *database) IsAdmin(userId int64) bool {
	sqlStr := "select * from `admin` where userId=?"
	row := db.Db.QueryRow(sqlStr, userId)
	var id int64
	row.Scan(&id)
	return id != 0
}
