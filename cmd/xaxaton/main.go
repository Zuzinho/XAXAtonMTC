package main

import (
	"XAXAtonMTC/pkg/env"
	"XAXAtonMTC/pkg/handlers/connection"
	"XAXAtonMTC/pkg/handlers/login"
	"XAXAtonMTC/pkg/handlers/message"
	musichandler "XAXAtonMTC/pkg/handlers/music"
	roomhandler "XAXAtonMTC/pkg/handlers/room"
	userhandler "XAXAtonMTC/pkg/handlers/user"
	"XAXAtonMTC/pkg/middleware"
	"XAXAtonMTC/pkg/music"
	"XAXAtonMTC/pkg/room"
	"XAXAtonMTC/pkg/sender"
	"XAXAtonMTC/pkg/user"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	connStr := env.MustDBConnectionString()

	userRepo, err := user.NewUsersDBRepository(connStr)
	if err != nil {
		log.Fatal(err)
	}

	roomRepo := room.NewRoomsMemoryRepository()

	senderRepo := sender.NewSendersMemoryRepository()

	songsRepo, err := music.NewSongsDBRepository(connStr)
	if err != nil {
		log.Fatal(err)
	}

	loginHandler := login.NewHandler(userRepo, env.MustSmsApiID())

	router.HandleFunc("/login", loginHandler.Login).Methods("POST")
	router.HandleFunc("/login/checksms", loginHandler.CheckSms).Methods("POST")

	rmHandler := roomhandler.NewHandler(userRepo, roomRepo, senderRepo)

	router.HandleFunc("/rooms", rmHandler.Rooms).Methods("GET")
	router.HandleFunc("/rooms", rmHandler.Create).Methods("POST")
	router.HandleFunc("/rooms/{room_id:[0-9]+}/{user_id:[0-9]+}", rmHandler.AddListener).Methods("POST")
	router.HandleFunc("/rooms/{room_id:[0-9]+}/{user_id:[0-9]+}", rmHandler.DeleteListener).Methods("DELETE")
	router.HandleFunc("/rooms/{room_id:[0-9]+}/{user_id:[0-9]+}/{token:[0-9A-z]+}", rmHandler.AddListenerByRef).Methods("POST")

	usHandler := userhandler.NewUserHandler(userRepo)

	router.HandleFunc("/user/{user_id:[0-9]+}", usHandler.Update).Methods("PUT")

	connHandler := connection.NewHandler(roomRepo, senderRepo)

	router.HandleFunc("/connect/{user_id:[0-9]+}", connHandler.Connect)

	musicHandler := musichandler.NewHandler(roomRepo, senderRepo, songsRepo)

	router.HandleFunc("/music/{room_id:[0-9]+}/{music_id:[0-9]+}", musicHandler.Music).Methods("GET")

	messageHandler := message.NewHandler(roomRepo, senderRepo)

	router.HandleFunc("/message/{room_id:[0-9]+}", messageHandler.Message).Methods("POST")

	appRouter := middleware.AccessLogMiddleware(router)
	appRouter = middleware.PanicMiddleware(appRouter)

	log.Println("starting on :8080")
	log.Println(http.ListenAndServe(":8080", appRouter))
}
