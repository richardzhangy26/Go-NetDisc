package handler

import (
	"encoding/json"
	"fmt"
	"go-netdisc/meta"
	"go-netdisc/util"
	"io"
	"net/http"
	"os"
	"time"
)

// 上传文件处理
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//接受文件流存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("filed to get data err: %v\n", err)
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2021-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("filed to create file err: %v\n", err)
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("filed to write file err: %v\n", err)
			return
		}
		newFile.Seek(0, 0)
		fileMeta.Filesha1 = util.FileSha1(newFile)
		fmt.Printf("fileMeta.filesha1: %v\n", fileMeta.Filesha1)
		// meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	} else if r.Method == "GET" {
		//返回上传html页面
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		io.WriteString(w, string(data))
	}
}
func UploadsucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "上传成功")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fmt.Printf("filehash: %s\n", filehash)
	// fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	fmt.Println(string(data))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

// 下载处理handler
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	// 文件大时建议使用流失读取
	data, err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("Content-Disposition", "attachment; filename="+fm.FileName)
	w.Write(data)
}
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filesha1 := r.Form.Get("filehash")
	fmeta := meta.GetFileMeta(filesha1)

	os.Remove(fmeta.Location)
	meta.RemoveFileMeta(filesha1)
	w.WriteHeader(http.StatusOK)
}
