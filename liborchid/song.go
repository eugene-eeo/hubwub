package liborchid

import "os"
import "path/filepath"
import "github.com/faiface/beep/mp3"
import "github.com/dhowden/tag"

func FindSongs(dir string) (songs []*Song, err error) {
	songs = []*Song{}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".mp3" {
			songs = append(songs, NewSong(path))
		}
		return nil
	})
	return
}

type Song struct {
	path string
}

func NewSong(path string) *Song {
	return &Song{path: path}
}

func (s *Song) Name() string {
	u := filepath.Base(s.path)
	ext := filepath.Ext(u)
	return u[:len(u)-len(ext)]
}

func (s *Song) file() (*os.File, error) {
	return os.Open(s.path)
}

func (s *Song) Stream() (*Stream, error) {
	f, err := s.file()
	if err != nil {
		return nil, err
	}
	stream, format, _ := mp3.Decode(f)
	return NewStream(stream, format), nil
}

func (s *Song) Tags() (tag.Metadata, error) {
	f, err := s.file()
	if err != nil {
		return nil, err
	}
	return tag.ReadFrom(f)
}
