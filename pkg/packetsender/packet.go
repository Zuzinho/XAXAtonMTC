package packetsender

type PacketType int

const (
	SONG    PacketType = 0
	MESSAGE PacketType = 1
)

const PacketSize = 256

type Packet struct {
	Data       []byte     `json:"data"`
	MetaData   []byte     `json:"meta_data"`
	PacketType PacketType `json:"packet_type"`
	IsNext     bool       `json:"is_next"`
}

func NewPacket(data, metaData []byte, packetType PacketType, isNext bool) *Packet {
	return &Packet{
		Data:       data,
		MetaData:   metaData,
		PacketType: packetType,
		IsNext:     isNext,
	}
}

type PacketSender interface {
	NextPacket() (*Packet, error)
}
