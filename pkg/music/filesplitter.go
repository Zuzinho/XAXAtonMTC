package music

import (
	"XAXAtonMTC/pkg/packetsender"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	musicPath          = "./../../music/"
	BitratePacketCount = 2000
)

type FileSplitter struct {
	buf      []byte
	metadata []byte
	offset   int
}

func NewFileSplitter(songName, authorName string) (*FileSplitter, error) {
	file, err := os.Open(fmt.Sprintf("%s - %s.mp3", authorName, songName))
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	metadata, err := json.Marshal(struct {
		AuthorName string `json:"author_name"`
		MusicName  string `json:"music_name"`
	}{
		AuthorName: authorName,
		MusicName:  songName,
	})
	if err != nil {
		return nil, err
	}

	log.Println("metadata: ", string(metadata))

	return &FileSplitter{
		buf:      buf,
		metadata: metadata,
		offset:   0,
	}, nil
}

func (splitter *FileSplitter) NextPacket() (*packetsender.Packet, error) {
	if len(splitter.buf) <= packetsender.PacketSize+splitter.offset {
		return packetsender.NewPacket(splitter.buf[splitter.offset:], splitter.metadata, packetsender.SONG, false), nil
	}

	buf := splitter.buf[splitter.offset : splitter.offset+packetsender.PacketSize]

	splitter.offset += packetsender.PacketSize

	return packetsender.NewPacket(buf, splitter.metadata, packetsender.SONG, true), nil
}
