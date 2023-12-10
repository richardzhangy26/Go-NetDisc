package db

import (
	"database/sql"
	"fmt"
	mydb "go-netdisc/db/mysql"
)

// 文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`)values(?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("filed to prepare statement,err" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println("filed to exec statement,err" + err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("file with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false

}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 数据库查询文件meta
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select `file_sha1`,`file_name`,`file_size`,`file_addr` from tbl_file where `file_sha1`=? and status=1 limit 1",
	)
	if err != nil {
		fmt.Println("filed to prepare statement,err" + err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash,
		&tfile.FileName,
		&tfile.FileSize,
		&tfile.FileAddr,
	)
	if err != nil {
		fmt.Println("filed to exec statement,err" + err.Error())
		return nil, err
	}
	return &tfile, nil
}
