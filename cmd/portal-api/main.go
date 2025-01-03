package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krishkumar84/bdcoe-golang-portal/pkg/config"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/auth"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/test"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/users"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/middleware"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/storage/mongodb"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/contest"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/question"
	"github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/testcase"
	// "github.com/krishkumar84/bdcoe-golang-portal/pkg/http/handler/users"
)

func main() {

	// load config

	cfg := config.MustLoad()


	//database
	//storage
    storage, err := mongodb.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Database connected",cfg.DatabaseName)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JwtSecret)

	// Setup routes
	router := http.NewServeMux()

	router.HandleFunc("GET/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to BDCOE Portal API server is dockerized up and running"))
	})

	router.Handle("GET /api/user/test", 
    authMiddleware.Authenticate(
        http.HandlerFunc(test.TestUserRoute),
    ),
)

router.Handle("GET /api/admin/test", 
    authMiddleware.Authenticate(
        authMiddleware.RequireAdmin(
            http.HandlerFunc(test.TestAdminRoute),
        ),
    ),
)

	router.HandleFunc("POST /api/signup",users.New(storage))
	router.HandleFunc("POST /api/login",auth.Login(storage,cfg.JwtSecret))
	router.HandleFunc("POST /api/contest",contest.CreateContest(storage))
	router.HandleFunc("POST /api/question",question.CreateQuestion(storage))
	router.HandleFunc("POST /api/testcase",testcase.CreateTestCase(storage))
	router.HandleFunc("GET /api/contest",contest.GetAllContests(storage))
    // router.HandleFunc("GET /api/users/{id}",users.GetById(storage))
	// router.HandleFunc("GET /api/users",users.GetAll(storage))
    
	//start server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

    fmt.Println("Server is running on port", cfg.Addr)
  
     done := make(chan os.Signal,1)

	 signal.Notify(done, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

	}()

	<-done

	slog.Info("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx) ; err != nil {

		slog.Error("Server Shutdown Failed",slog.String("error",err.Error()))
	}

	slog.Info("Server ShutDown Properly")
}