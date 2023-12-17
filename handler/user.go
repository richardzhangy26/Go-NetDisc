package handler

import (
	"fmt"
	dblayer "go-netdisc/db"
	"go-netdisc/util"
	"net/http"
	"os"
	"time"
)

const (
	pwd_salt = "*#890"
)

// 处理用户注册规则
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	// username和passwd校验
	if len(username) < 3 || len(passwd) < 5 {
		http.Error(w, "用户名或密码长度不得小于3", http.StatusBadRequest)
		return
	}
	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("注册成功"))
	} else {
		w.Write([]byte("注册失败"))
	}

}
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// data, err := ioutil.ReadFile("./static/view/signin.html")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password + pwd_salt))
	// 1.校验用户名及密码
	pwdChecked := dblayer.UserSignIn(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("登录失败"))
		return
	}
	// 2.生成访问凭证token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("更新token失败"))
		return
	}
	// 3.登陆成功后重定向
	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	// http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}
func GenToken(username string) string {
	// md5(username+timestap+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	// token := r.Form.Get("token")
	// // 2.校验token
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.Write([]byte("token无效"))
	// 	return
	// }
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())

}

func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}
