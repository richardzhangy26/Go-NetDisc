package meta

import (
	mydb "go-netdisc/db"
)

type FileMeta struct {
	Filesha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.Filesha1] = fmeta
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.Filesha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 获取文件元信息
func GetFileMeta(filesha1 string) FileMeta {
	fmeta := fileMetas[filesha1]
	return fmeta
}

// 从mysql获取文件元信息
func GetFileMetaDB(filesha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(filesha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		Filesha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil

}

// RemoveFileMeta 删除文件元信息(最好加上锁保证线程同步)
func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}
