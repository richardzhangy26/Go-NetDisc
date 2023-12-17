package db

import (
	"fmt"
	mydb "go-netdisc/db/mysql"
)

// 通过用户名和密码注册用户
func UserSignup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`) values(?,?)")
	if err != nil {
		fmt.Println("failed to insert,err : ", err.Error())
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("failed to insert,err : ", err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}
	return false
}

// 用户登陆校验
func UserSignIn(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=?")
	if err != nil {
		fmt.Println("failed to select,err : ", err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("failed to select,err : ", err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: ", err.Error())
		return false
	}
	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == password {
		return true
	} else {
		return false
	}

}

// 刷新用户登陆的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token(`user_name`,`user_token`)values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// User : 用户表model
type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// GetUserInfo : 查询用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
