package music

import "github.com/jackc/pgx"

type SongsDBRepository struct {
	connConfig pgx.ConnConfig
}

func NewSongsDBRepository(connString string) (*SongsDBRepository, error) {
	conf, err := pgx.ParseConnectionString(connString)
	if err != nil {
		return nil, err
	}

	return &SongsDBRepository{
		connConfig: conf,
	}, nil
}

func (repo *SongsDBRepository) SelectByID(songID uint32) (*Song, error) {
	conn, err := pgx.Connect(repo.connConfig)
	if err != nil {
		return nil, err
	}

	var authorName, musicName string

	err = conn.QueryRow("select author_name, music_name from music where music_id = $1", songID).
		Scan(&authorName, &musicName)
	if err != nil {
		return nil, err
	}

	return NewSong(songID, authorName, musicName), nil
}

func (repo *SongsDBRepository) SelectAll() ([]*Song, error) {
	conn, err := pgx.Connect(repo.connConfig)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query("select * from music")
	if err != nil {
		return nil, err
	}

	songs := make([]*Song, 0)
	for rows.Next() {
		var songID uint32
		var authorName, musicName string

		err = rows.Scan(&songID, &authorName, &musicName)
		if err != nil {
			continue
		}

		songs = append(songs, NewSong(songID, authorName, musicName))
	}

	return songs, nil
}
