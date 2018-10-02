package main

import (
	"log"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"bytes"
	"rush/config"
)

/*
	web大前端只做两件事
	1. 转发业务请求
		转发分为proxy和api模式
		因为用ajax请求，浏览器都会做跨域限制
		比如渲染出来的page域名是127.0.0.1:8000
		streamServer的域名是127.0.0.1:9000
		不能在js里直接调用另一个域名的东西

	2. 把模板化的东西渲染转发给前端
 */

//这个client做的是代理的事情，把调用api的事都代理了
var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	//偷改url地址，改为loadbalance的host
	u, _ := url.Parse(b.Url)
	u.Host = config.GetLbAddr() + ":" + u.Port()
	newUrl := u.String()

	switch b.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", newUrl, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
		return
	}
}

//这里r不是request，就是response，是api接口返回的response
func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}