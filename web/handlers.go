package main

import (
	"html/template"
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"net/http/httputil"
	"rush/config"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")

	//============这几行代码可以显示程序的位置===============
	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(dir)
	//==================================================

	//如果username或者session有一个取不到，那么就返回登录页面
	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "Dora"}
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing templates home.html error: %s", e)
			return
		}

		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")

	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error: %s", e)
		return
	}

	t.Execute(w, p)
}

func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}  //做了归一化处理
	if err := json.Unmarshal(res, apibody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}

	request(apibody, w, r)
	defer r.Body.Close()
}

func proxyVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//url.Parse就是把字符串url变为URL对象
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func proxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse("http://" + config.GetLbAddr() + ":9000/")
	//在这里，反向代理只是简单的做了域名转换，把8080换为9090
	//一般来说，反向代理还可以重写头部，把header里的东西修改，但这里用不到
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
