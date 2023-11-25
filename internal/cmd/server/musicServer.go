package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

const PORT = ":3017"

type Command struct {
	UserType string
	UserID   int32
	RoomID   int32
	Command  string
	Song     string
}

type room struct {
	users       map[int32]string // map[personalID]{address}
	song        string
	connections map[int32]net.Conn
	songFiles   []os.DirEntry
	m3u8        []byte
	stop        chan struct{}
}

var listeners map[int32]room // map[RoomID]map[personalID]{address}

func main() {
	//arguments := os.Args
	//if len(arguments) == 1 {
	//	fmt.Println("Please provide a port number!")
	//	return
	//}
	//PORT := ":" + arguments[1]

	l, err := net.Listen("tcp4", PORT)
	fmt.Println("start commands server on port", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	//go musicServer()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	// we create a decoder that reads directly from the socket
	d := json.NewDecoder(c)
	var cmd Command
	err := d.Decode(&cmd)
	if err != nil {
		fmt.Println(err)
		c.Close()
		return
	}

	switch cmd.UserType {
	case "Creator":
		switch cmd.Command {
		case "Create":
			CreateRoom(c, cmd)
		case "Start":
			StartStream(cmd)
		case "Stop":
			StopStream(cmd)
		case "Pause":
			//TODO ???
		}

	case "Listener":
		switch cmd.Command {
		case "Join":
			JoinRoom(c, cmd)
		case "Leave":
			LeaveRoom(cmd)
		}
	}
	c.Close()
}

func LeaveRoom(cmd Command) {
	roomTemp, exist := listeners[cmd.RoomID]
	if !exist {
		fmt.Printf("RoomID %d doesn't exist", cmd.RoomID)
		return
	}

	// if connection to user is open - close it and delete from connections
	if con, exists := roomTemp.connections[cmd.UserID]; exists {
		con.Close()
		delete(roomTemp.connections, cmd.UserID)
	}

	// delete user from room
	delete(roomTemp.users, cmd.UserID)
	listeners[cmd.RoomID] = roomTemp
}

func JoinRoom(c net.Conn, cmd Command) {
	roomTemp, exist := listeners[cmd.RoomID]
	if !exist {
		fmt.Printf("RoomID %d doesn't exist", cmd.RoomID)
		return
	}
	roomTemp.users[cmd.UserID] = fmt.Sprintf("%s", c.RemoteAddr().String())
	c.Write(listeners[cmd.RoomID].m3u8)
	listeners[cmd.RoomID] = roomTemp
}

func CreateRoom(c net.Conn, cmd Command) {
	// Creating room, and adding self
	usersC := make(map[int32]string)
	usersC[cmd.UserID] = fmt.Sprintf("%s", c.RemoteAddr().String())

	songFiles, err := os.ReadDir(fmt.Sprintf("./songs/" + cmd.Song))
	//entries, err := os.ReadDir("./songs/BachGavotteShort")
	if err != nil {
		log.Fatal(err)
	}

	m3u8, err := ioutil.ReadFile(fmt.Sprintf("./songs/" + cmd.Song + "/outputlist.m3u8")) // b has type []byte
	if err != nil {
		log.Fatal(err)
	}

	listeners[cmd.RoomID] = room{
		users:     usersC,
		song:      cmd.Song,
		songFiles: songFiles,
		m3u8:      m3u8,
		stop:      make(chan struct{}),
	}

	c.Write(listeners[cmd.RoomID].m3u8)
}

func StartStream(cmd Command) {
	// Creating connections to users
	connections := make(map[int32]net.Conn, len(listeners[cmd.RoomID].users))
	for _, s := range listeners[cmd.RoomID].users {
		connClient, _ := net.Dial("udp", s)
		//connClient, _ := net.Dial("tcp", s) // TODO Выбрать как подключаемся для стрима
		connections[cmd.UserID] = connClient
	}

	listeners[cmd.RoomID] = room{
		users:       listeners[cmd.RoomID].users,
		song:        listeners[cmd.RoomID].song,
		connections: connections,
	}

	go func(curRoom room) {
		var nextChunk []byte
		var err error
		var name string

	out:
		for i := 0; i < len(curRoom.songFiles)-1; i++ { // do not send the last file - m3u8
			select {
			case <-curRoom.stop: // stop stream and get out
				break out
			default:
				name = curRoom.songFiles[i].Name()
				nextChunk, err = ioutil.ReadFile(fmt.Sprintf("./songs/" + cmd.Song + "/" + name)) // b has type []byte
				if err != nil {
					log.Fatal(err)
				}
				for _, conn := range curRoom.connections {
					conn.Write(nextChunk)
					// go conn.Write(nextChunk)
					// time.Sleep(time.Millisecond*100) // length of chunk
				}
			}
		}
	}(listeners[cmd.RoomID])
}

func StopStream(cmd Command) {
	listeners[cmd.RoomID].stop <- struct{}{} // sending signal to stop transmitting
	// closing connections ???
	for _, conn := range listeners[cmd.RoomID].connections {
		conn.Close()
	}
	//TODO stop stream first send to channel
}

func musicServer() {
	// configure the songs directory name and port
	const songsDir = "songs"
	const port = 8080

	// add a handler for the song files
	http.Handle("/", addHeaders(http.FileServer(http.Dir(songsDir))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", songsDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// http://localhost:8080/BachGavotteShort/outputlist.m3u8
// https://hlsjs-dev.video-dev.org/demo/

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
