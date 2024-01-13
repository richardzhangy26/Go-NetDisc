package handler

import (
	"fmt"
	rPool "go-netdisc/cache/redis"
	dblayer "go-netdisc/db"
	"math"
	"os"
	"path"
	"strings"
	"time"

	// "go-netdisc/redis"
	"go-netdisc/util"
	"net/http"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// 分块上传结构体
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadId   string
	Chunksize  int
	ChunkCount int
}

// 初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1.解析用户请求
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "参数错误", nil).JSONBytes())
		return
	}
	// 2.获得redis连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// 3.初始化用户分块上传信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadId:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		Chunksize:  5 * 1024 * 1024, //5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 4.将初始化信息传到redis缓存
	rConn.Do("Hset", "MP_"+upInfo.UploadId, "filehash", upInfo.FileHash)
	rConn.Do("Hset", "MP_"+upInfo.UploadId, "filesize", upInfo.FileSize)
	rConn.Do("Hset", "MP_"+upInfo.UploadId, "chunkcount", upInfo.ChunkCount)
	// 5.响应请求到客户端
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())

}

// 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// 1.解析用户请求
	r.ParseForm()
	username := r.Form.Get("username")
	uploadId := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")
	// 2.获得redis连接池中的一个连接
	fmt.Println("username", username)
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	// 3.获得文件句柄，用于存储分块内容
	fapth := "/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(fapth), 0744)
	fd, err := os.Create(fapth)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "文件创建失败", nil).JSONBytes())
		return
	}
	defer fd.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// 4.更新redis缓存状态
	rConn.Do("Hset", "MP_"+uploadId, "chkidx_"+chunkIndex, 1)
	// 5.返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// 通知上传合并接口
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1.解析用户请求
	r.ParseForm()
	upid := r.Form.Get("uploadId")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	filename := r.Form.Get("filename")
	// 2.获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	// 3.通过uploadId查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("Hgetall", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete uploadfailed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount += 1
		}
		if totalCount != chunkCount {
			w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
			return
		}

	}
	// 4.合并分块
	// 5.更新唯一文件表及用户文件表
	dblayer.OnFileUploadFinished(filehash, filename, int64(filesize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
