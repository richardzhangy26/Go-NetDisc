package db

import (
	"fmt"
	mydb "go-netdisc/db/mysql"
	"time"
)

// UserFile:用户列表结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

func OnUserFileUploadFinished(username string, filehash string, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file(`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) " +
			"values(?,?,?,?,?)")
	if err != nil {
		return false
	}
	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now().Format("2006-01-02 15:04:05"))
	return err == nil
}

// 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from" +
			" tbl_user_file where user_name=? limit ?")

	if err != nil {
		fmt.Println("prepare error:", err.Error())
		return nil, err
	}
	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println("query error:", err.Error())
		return nil, err
	}
	var userfiles []UserFile
	for rows.Next() {
		var userfile UserFile
		err = rows.Scan(&userfile.FileHash, &userfile.FileName, &userfile.FileSize, &userfile.UploadAt, &userfile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		userfiles = append(userfiles, userfile)
	}
	return userfiles, nil
}
