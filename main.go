package main

import (
	"fmt"
	"go-netdisc/handler"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc", handler.HTTPInterceptor(handler.UploadsucHandler))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err)
		return
	}
}
