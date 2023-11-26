package music

type Song struct {
	ID         uint32 `json:"id"`
	AuthorName string `json:"author_name"`
	MusicName  string `json:"music_name"`
}

func NewSong(id uint32, authorName, musicName string) *Song {
	return &Song{
		ID:         id,
		AuthorName: authorName,
		MusicName:  musicName,
	}
}

type SongsRepo interface {
	SelectByID(uint32) (*Song, error)
	SelectAll() ([]*Song, error)
}
