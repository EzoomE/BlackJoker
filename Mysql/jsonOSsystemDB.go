package jsonOSsystemDB

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	userName = "root"
	password = "LUrui15296092828_"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "jsonossystem"
)

var DB *sql.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, password, ip, port, dbName)
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
func InitMysql(jsonDB string) error {
	if DB == nil {
		return fmt.Errorf("数据库未连接")
	}
	sqlite := "CREATE TABLE IF NOT EXISTS Botinit (Cookie VARCHAR(255), SystemName VARCHAR(255), SysPath VARCHAR(255), Cpu VARCHAR(255), Cpuarchitecture VARCHAR(255), IP VARCHAR(255),AdminName VARCHAR(255));"
	_, err := DB.Exec(sqlite)
	if err != nil {
		return err
	}
	return nil
}
func Infinite() bool {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"
	var count int

	err := DB.QueryRow(query, dbName, "Botinit").Scan(&count)
	if err != nil {
		log.Println(err)
		return true
	}
	return count == 0
}

func CheckCookieExists(cookieValue string) bool {
	err := InitDB()
	if err != nil {
		fmt.Println("连接数据库错误")
	}
	db := DB
	query := "SELECT COUNT(*) FROM jsonossystem.botinit WHERE Cookie = ?"
	result, err := db.Exec(query, cookieValue)
	if err != nil {
		log.Println("查询出错:", err)
		return false
	}
	affectedRows, _ := result.RowsAffected()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("关闭数据库错误")
		}
	}(db)
	return affectedRows > 0
}

func MysqlMain(jsonDB string) error {
	err := InitDB()
	if err != nil {
		log.Println(err)
		return err
	}
	if Infinite() {
		err := InitMysql(jsonDB)
		if err != nil {
			log.Println(err)
		}
	}
	func() {
		values := strings.Split(jsonDB, "|")
		if len(values) != 7 {
			log.Println("数据格式错误")
			return
		}
		sql := "INSERT INTO Botinit (Cookie,SystemName,SysPath,Cpu,Cpuarchitecture,IP,AdminName) VALUES (?,?,?,?,?,?,?);"
		checkQuery := "SELECT COUNT(*) FROM Botinit WHERE Cookie = ? AND SystemName = ? AND SysPath = ? AND Cpu = ? AND Cpuarchitecture = ? AND IP = ? AND AdminName = ?;"
		var count int
		err := DB.QueryRow(checkQuery, values[0], values[1], values[2], values[3], values[4], values[5], values[6]).Scan(&count)
		if err != nil {
			log.Println(err)
			return
		}
		if count > 0 {
			return
		}
		_, err = DB.Exec(sql, values[0], values[1], values[2], values[3], values[4], values[5], values[6])
		if err != nil {
			log.Println(err)
		}
	}()
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			fmt.Println("关闭数据库错误")
		}
	}(DB)
	return nil
}
