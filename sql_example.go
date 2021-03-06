package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"

	PASSWORD = "123456"

	NETWORK = "tcp"

	SERVER = "localhost"

	PORT = 3306

	DATABASE = "test"
)

//user表结构体定义

type User struct {
	Id int `json:"id" form:"id"`

	Username string `json:"username" form:"username"`

	Password string `json:"password" form:"password"`

	Status int `json:"status" form:"status"` // 0 正常状态， 1删除

	Createtime int64 `json:"createtime" form:"createtime"`
}

func main() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
}

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	// dsn := "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	// db, err = sql.Open("mysql", dsn)
	// db, err := sql.Open("mysql", "root:root@(127.0.0.1)/zan_ku")
	dsn := "root:123456@tcp(127.0.0.1:3306)/runoob"
	// conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	} else {
		fmt.Println("连接数据库成功")
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超时的连接就close
	db.SetMaxOpenConns(100)                  //设置最大连接数

	CreateTable(db)
	InsertData(db)
	// QueryOne(db)
	QueryMulti(db)
	UpdateData(db)
	DeleteData(db)

	return nil
}

//创建表
func CreateTable(DB *sql.DB) {
	sql := `CREATE TABLE IF NOT EXISTS users(id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,username VARCHAR(64),password VARCHAR(64),status INT(4),createtime INT(10)); `
	if _, err := DB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)

		return

	}

	fmt.Println("create table successd")

}

//添加数据
func InsertData(DB *sql.DB) {
	result, err := DB.Exec("insert INTO users(username,password,status,createtime) values(?,?,?,?)", "test", "123456", 0, time.Now().Second())

	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)

		return

	}

	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID

	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)

		return

	}

	fmt.Println("Insert data id:", lastInsertID)

	rowsaffected, err := result.RowsAffected() //通过RowsAffected获取受影响的行数

	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)

		return

	}

	fmt.Println("Affected rows:", rowsaffected)

}

//查询单行
func QueryOne(DB *sql.DB) {
	user := new(User) //用new()函数初始化一个结构体对象

	row := DB.QueryRow("select id,username,password from users where id=?", 2)

	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错

	if err := row.Scan(&user.Id, &user.Username, &user.Password); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)

		return

	}

	fmt.Println("Single row data:", *user)

}

//查询多行
func QueryMulti(DB *sql.DB) {
	user := new(User)

	rows, err := DB.Query("select id,username,password,status,createtime from users where id = ?", 6)

	defer func() {
		if rows != nil {
			rows.Close() //关闭掉未scan的sql连接

		}

	}()

	if err != nil {
		fmt.Printf("Query failed,err:%v\n", err)

		return

	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Status, &user.Createtime) //不scan会导致连接不释放

		if err != nil {
			fmt.Printf("Scan failed,err:%v\n", err)

			return

		}

		fmt.Println("scan successd:", *user)

	}

}

//更新数据
func UpdateData(DB *sql.DB) {
	result, err := DB.Exec("UPDATE users set password=? where id=?", "111111", 3)

	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)

		return

	}

	fmt.Println("update data successd:", result)

	rowsaffected, err := result.RowsAffected()

	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v\n", err)

		return

	}

	fmt.Println("Affected rows:", rowsaffected)

}

//删除数据
func DeleteData(DB *sql.DB) {

	result, err := DB.Exec("delete from users where id=?", 2)

	if err != nil {
		fmt.Printf("Insert failed,err:%v\n", err)

		return

	}

	fmt.Println("delete data successd:", result)

	rowsaffected, err := result.RowsAffected()

	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v\n", err)

		return

	}

	fmt.Println("Affected rows:", rowsaffected)

}
