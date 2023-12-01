package meta

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

// 获取文件元信息
func GetFileMeta(filesha1 string) FileMeta {
	fmeta := fileMetas[filesha1]
	return fmeta
}

// RemoveFileMeta 删除文件元信息(最好加上锁保证线程同步)
func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}
