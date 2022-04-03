package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
)

//
type App struct {
	Router      *mux.Router
	Middlewares *Middleware
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
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
	a.Middlewares = &Middleware{}
	a.initializeRoutes()
}

//设置路由
func (a *App) initializeRoutes() {
	// a.Router.HandleFunc("/api/shorten", a.createShortlink).Methods("POST")
	// a.Router.HandleFunc("/api/info", a.getShortlinkInfo).Methods("GET")
	// a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")

	m := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)

	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortlink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortlinkInfo)).Methods("GET")
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

//创建短链接
func (a *App) createShortlink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithErr(w, StatusError{http.StatusBadRequest, fmt.Errorf("json格式错误 %v", r.Body)})
		return
	}

	if err := validator.Validate(req); err != nil {
		respondWithErr(w, StatusError{http.StatusBadRequest, fmt.Errorf("参数校验错误 %v", req)})
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
	panic(s)
}

//短链跳转
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
