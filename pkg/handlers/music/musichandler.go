package music

import (
	"XAXAtonMTC/pkg/handlers"
	"XAXAtonMTC/pkg/music"
	"XAXAtonMTC/pkg/packetsender"
	"XAXAtonMTC/pkg/room"
	"XAXAtonMTC/pkg/sender"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	RoomsRepo    room.RoomsRepo
	SendersRepo  sender.SendersRepo
	PacketSender packetsender.PacketSender
	SongsRepo    music.SongsRepo
}

func NewHandler(roomsRepo room.RoomsRepo, sendersRepo sender.SendersRepo, songsRepo music.SongsRepo) *Handler {
	return &Handler{
		RoomsRepo:   roomsRepo,
		SendersRepo: sendersRepo,
		SongsRepo:   songsRepo,
	}
}

func (handler *Handler) Music(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomID, err := handlers.UInt32FromMapString(vars, "room_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	musicID, err := handlers.UInt32FromMapString(vars, "music_id")
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	if handler.PacketSender == nil {
		song, err := handler.SongsRepo.SelectByID(musicID)
		if err != nil {
			handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
			return
		}

		splitter, err := music.NewFileSplitter(song.MusicName, song.AuthorName)
		if err != nil {
			handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
			return
		}

		handler.PacketSender = splitter
	}

	userIDs, err := handler.RoomsRepo.UsersID(roomID)
	if err != nil {
		handlers.HttpJSONErr(w, err, http.StatusBadRequest)
		return
	}

	for i := 0; i < music.BitratePacketCount; i++ {
		packet, err := handler.PacketSender.NextPacket()
		if err != nil {
			handlers.HttpJSONErr(w, err, http.StatusInternalServerError)
			continue
		}

		handler.SendersRepo.SendPacket(userIDs, func(uint32) bool {
			return true
		}, packet)

		if !packet.IsNext {
			break
		}
	}

	w.WriteHeader(http.StatusOK)
}
