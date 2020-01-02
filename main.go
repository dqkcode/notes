package main

import (
	"net/http"

	"github.com/dqkcode/notes/internal/app/user"

	"github.com/globalsign/mgo"

	MongoDBConf "github.com/dqkcode/notes/internal/app/config/mongodb"

	serverConfig "github.com/dqkcode/notes/internal/app/config/server"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func handler1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("from server!\n"))

}
func main() {

	serverConf, err := serverConfig.Load()
	logrus.Infof("server conf info : %s", serverConf)
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	dbConf, err := MongoDBConf.Load()
	if err != nil {
		panic(err)
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    dbConf.Addrs,
		Database: dbConf.Database,
		Username: dbConf.UserName,
		Password: dbConf.Password,
	})
	if err != nil {
		panic(err)
	}

	userRepo := user.NewMongoDBRepository(session)
	userSrv := user.NewService(userRepo)
	userHandler := user.NewHandler(userSrv)

	//Post
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(userHandler.Register)
	router.Path("/login").Methods(http.MethodPost).HandlerFunc(userHandler.Login)
	//Get
	router.Path("/showinfos").Methods(http.MethodGet).HandlerFunc(userHandler.ShowInfo)

	server := http.Server{
		Addr:         serverConf.Addr,
		ReadTimeout:  serverConf.ReadTimeOut,
		WriteTimeout: serverConf.WriteTimeOut,
		Handler:      router,
	}
	logrus.Infof("server is listinging on: %s", serverConf.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
