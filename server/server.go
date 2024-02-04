package server

import (
	"crypto/tls"
	"kakeibo-api/db"
	"kakeibo-api/domain"
	"log"
	loginRepo "kakeibo-api/login/repository"
	loginServices "kakeibo-api/login/services"
	userRepo "kakeibo-api/user/repository"
	userServices "kakeibo-api/user/services"
	"net/http"
	"os"
	"time"
)

type Server struct {
	UserService  domain.UserService
	LoginService domain.LoginService
}

func NewAdminServer() *Server {
	ur := userRepo.NewUserRepository(db.DataStore)
	loginSrvc := loginServices.NewLoginService(loginRepo.NewLoginRepository(db.DataStore), ur)
	userSrvc := userServices.NewUserService(ur)

	return &Server{
		LoginService: loginSrvc,
		UserService:  userSrvc,
	}
}

var serverStart time.Time

func Start() {
	srv := NewAdminServer()
	router := srv.NewRouter()

	server := &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}

	log.Printf("Server started at: %s", server.Addr)

	serverStart = time.Now()

	if os.Getenv("ENABLE_TLS") == "true" {
		log.Fatal(server.ListenAndServeTLS(os.Getenv("SSL_CRT_PATH"), os.Getenv("SSL_KEY_PATH")))
	} else {
		//Globally disabling SSL certificate check
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		log.Fatal(server.ListenAndServe())
	}
}
