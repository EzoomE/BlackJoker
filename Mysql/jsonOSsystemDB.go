package jsonOSsystemDB

import (
	VARBLACK "BlackJoker/VarData"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", VARBLACK.Config.MysqlUserName, VARBLACK.Config.MysqlPassword, VARBLACK.Config.MysqlIp, VARBLACK.Config.MysqlPort, VARBLACK.Config.MysqlDbName)
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = dbConn.Ping()
	if err != nil {
		return err
	}

	DB = dbConn
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	return nil
}
func InitMysql() error {
	if DB == nil {
		return fmt.Errorf("数据库未连接")
	}
	sqlite := "CREATE TABLE IF NOT EXISTS botinit (Cookie VARCHAR(255), SystemName VARCHAR(255), SysPath VARCHAR(255), Cpu VARCHAR(255), Cpuarchitecture VARCHAR(255), IP VARCHAR(255),AdminName VARCHAR(255));"
	_, err := DB.Exec(sqlite)
	if err != nil {
		return err
	}
	return nil
}
func Infinite() bool {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"
	var count int

	err := DB.QueryRow(query, VARBLACK.Config.MysqlDbName, "botinit").Scan(&count)
	if err != nil {
		log.Println("查询表是否存在时出错:", err)
		return false
	}
	return count > 0
}

func CheckCookieExists(cookieValue string) bool {
	err := InitDB()
	if err != nil {
		fmt.Println("连接数据库错误")
		return false
	}
	db := DB

	if !Infinite() {
		log.Print("正在为您创建表")
		err := InitMysql()
		if err != nil {
			log.Println(err)
			return false
		}
	}

	query := "SELECT COUNT(*) FROM jsonossystem.botinit WHERE Cookie = ?"
	var count int
	err = db.QueryRow(query, cookieValue).Scan(&count)
	if err != nil {
		log.Println("查询出错:", err)
		return false
	}

	return count > 0
}

func MysqlMain(jsonDB string) error {
	err := InitDB()
	if err != nil {
		log.Println("初始化数据库连接出错:", err)
		return err
	}

	func() {
		values := strings.Split(jsonDB, "|")
		if len(values) != 7 {
			log.Println("数据格式错误")
			return
		}

		sql := "INSERT INTO botinit (Cookie,SystemName,SysPath,Cpu,Cpuarchitecture,IP,AdminName) VALUES (?,?,?,?,?,?,?);"
		checkQuery := "SELECT COUNT(*) FROM botinit WHERE Cookie = ? AND SystemName = ? AND SysPath = ? AND Cpu = ? AND Cpuarchitecture = ? AND IP = ? AND AdminName = ?;"
		var count int
		err := DB.QueryRow(checkQuery, values[0], values[1], values[2], values[3], values[4], values[5], values[6]).Scan(&count)
		if err != nil {
			log.Println("查询数据时出错:", err)
			return
		}

		if count > 0 {
			return
		}

		_, err = DB.Exec(sql, values[0], values[1], values[2], values[3], values[4], values[5], values[6])
		if err != nil {
			log.Println("插入数据时出错:", err)
		}
	}()
	return nil
}

func GetClientSystem(hostname string) (string, error) {
	err := InitDB()
	if err != nil {
		log.Println("连接数据库错误")
	}
	db := DB
	var clientOs string
	query := "SELECT SystemName FROM jsonossystem.botinit WHERE Cookie = ?"
	err = db.QueryRow(query, hostname).Scan(&clientOs)
	if err != nil {
		if err == sql.ErrNoRows {
			return "未找到匹配的记录", nil
		}
		return "查询数据库错误", err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("关闭数据库错误")
		}
	}(db)
	return clientOs, nil
}
