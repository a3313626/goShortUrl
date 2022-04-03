package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

//
type App struct {
	Router *mux.Router
}

type shortenReq struct {
	URL                 string `json:"usr" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortlinkResp struct {
	Shortlink string `json:"shortlink"`
}

//初始化定义
func (a *App) Initialize() {
	//log设置
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//定义路由
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//设置路由
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")
}

//创建短链接
func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	if err := validator.Validate(req); err != nil {
		return
	}

	defer r.Body.Close()

	fmt.Printf("%v\n", req)

}

//查询短链详情
func (a *App) getShortlinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	fmt.Printf("%s\n", s)
}

//短链跳转
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
