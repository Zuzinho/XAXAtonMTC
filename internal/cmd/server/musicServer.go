package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
)

const (
	PacketRequestSong = iota
	PacketCoverArt
	PacketCoverArtEnd
	PacketMusicStreamSize
	PacketMusicStream
	PacketEndConnection
)

const PORT = ":3017"

type Packet struct {
	m_packetType uint32
	m_dataSize   uint32
	m_data       [256 - (4 * 2)]byte
}

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
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

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
	
	//// Main loop
	//clientDisconnect := false
	//for {
	//
	//	// Client sent a disconnect packet
	//	if clientDisconnect {
	//		break
	//	}
	//
	//	// Wait for a packet from the client
	//	var clientPacket Packet
	//	if ReadPacket(c, &clientPacket) {
	//		switch clientPacket.m_packetType {
	//
	//		case PacketRequestSong:
	//			songRequested := string(clientPacket.m_data[0 : clientPacket.m_dataSize-1]) // Drop terminating '\n', hence the m_dataSize - 1
	//			fmt.Println("Song req: ", songRequested)
	//
	//			SendCoverArt(c, songRequested)
	//			SendMusic(c, songRequested)
	//		case PacketEndConnection:
	//			fmt.Println("Disconnect packet received")
	//			clientDisconnect = true
	//		}
	//	}
	//}

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

	// delete from room
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
	listeners[cmd.RoomID] = roomTemp
}

func CreateRoom(c net.Conn, cmd Command) {
	// Creating room, and adding self
	usersC := make(map[int32]string)
	usersC[cmd.UserID] = fmt.Sprintf("%s", c.RemoteAddr().String())
	listeners[cmd.RoomID] = room{
		users: usersC,
		song:  cmd.Song}
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

	go func(map[int32]net.Conn) {
		for {
			select {
			// select and channel to end stream
			//case <-stop: break
			default:
				// ToDo Stream data
				for _, conn := range connections {
					// TODO
					conn.Write([]byte("NEXTCHUNK"))
					// go conn.Write([]byte("NEXTCHUNK"))
				}
			}
		}

	}(connections)
}

func StopStream(cmd Command) {
	for _, conn := range listeners[cmd.RoomID].connections {
		conn.Close()
	}
	//TODO stop stream first send to channel
}

func ReadPacket(conn net.Conn, packet *Packet) bool {

	buf := make([]byte, 256)
	packetRead := false
	// Read the incoming connection into the buffer
	for {
		packetLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Couldn't read packet")
			break
		}

		if packetLen == 256 {
			packetRead = true
			break // All data from packet received
		}
	}

	if packetRead {
		// Read it back to the Packet data structure
		packet.m_packetType = binary.LittleEndian.Uint32(buf[0:4])
		packet.m_dataSize = binary.LittleEndian.Uint32(buf[4:8])
		if packet.m_dataSize != 0 {
			copy(packet.m_data[0:256-(4*2)], buf[8:256])
		}
	}

	return packetRead
}

func SendPacket(conn net.Conn, packet *Packet) {

	// Transfer the content of the Packet to buf and transfer it over connection
	buf := make([]byte, 256)
	binary.LittleEndian.PutUint32(buf[0:], packet.m_packetType) // Packet type
	binary.LittleEndian.PutUint32(buf[4:], packet.m_dataSize)   // Data size
	copy(buf[8:256], packet.m_data[0:packet.m_dataSize])        // Data load

	conn.Write(buf)
}

func SendCoverArt(conn net.Conn, songRequested string) {

	// Open the song cover art request
	var coverArtFileName string
	coverArtFileName = "sfx/coverArt_" + songRequested
	fmt.Println("Opening cover art: ", coverArtFileName)

	coverArt, err := ioutil.ReadFile(coverArtFileName)
	if err != nil {
		fmt.Println("Couldn't open cover art, error: ", err.Error())
	}

	var coverArtLen = uint32(len(coverArt))
	var currentCoverArtPos uint32 = 0
	var dataLoadMaxSize uint32 = 256 - (4 * 2)

	for currentCoverArtPos < coverArtLen {
		var dataSize = dataLoadMaxSize

		if currentCoverArtPos+dataLoadMaxSize > coverArtLen {
			dataSize = coverArtLen - currentCoverArtPos
		}

		var coverArtPack Packet
		coverArtPack.m_packetType = PacketCoverArt
		coverArtPack.m_dataSize = dataSize
		copy(coverArtPack.m_data[0:coverArtPack.m_dataSize], coverArt[currentCoverArtPos:currentCoverArtPos+dataSize])
		SendPacket(conn, &coverArtPack)
		currentCoverArtPos += dataSize
	}

	var coverArtEndPack Packet
	coverArtEndPack.m_packetType = PacketCoverArtEnd
	coverArtEndPack.m_dataSize = 0
	SendPacket(conn, &coverArtEndPack)
}

func SendMusic(conn net.Conn, songRequested string) {

	musicFileName := "sfx/" + songRequested
	musicContent, err := ioutil.ReadFile(musicFileName)
	if err != nil {
		fmt.Println("Couldn't read music file")
	}

	var musicFileLen = uint32(len(musicContent))

	var musicSizePack Packet
	musicSizePack.m_packetType = PacketMusicStreamSize
	musicSizePack.m_dataSize = 4
	binary.LittleEndian.PutUint32(musicSizePack.m_data[0:], musicFileLen)
	SendPacket(conn, &musicSizePack)

	fmt.Println("Sending music...")

	var dataLoadMaxSize uint32 = 256 - (4 * 2)
	var currentMusicPos uint32 = 0
	for currentMusicPos < musicFileLen {
		var dataSize = dataLoadMaxSize

		if currentMusicPos+dataLoadMaxSize > musicFileLen {
			dataSize = musicFileLen - currentMusicPos
		}

		var musicPack Packet
		musicPack.m_packetType = PacketMusicStream
		musicPack.m_dataSize = dataSize
		copy(musicPack.m_data[0:musicPack.m_dataSize], musicContent[currentMusicPos:currentMusicPos+dataSize])
		SendPacket(conn, &musicPack)
		fmt.Print("*")
		currentMusicPos += dataSize
	}

	fmt.Println("")
	fmt.Println("All music has been sent!")
}

func handleServerConnection(c net.Conn) {
	// we create a decoder that reads directly from the socket
	d := json.NewDecoder(c)
	var msg Command // TODO
	err := d.Decode(&msg)
	fmt.Println(msg, err)

	c.Close()
}
