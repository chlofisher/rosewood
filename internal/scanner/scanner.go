package scanner

import (
	"os"
	"fmt"
	"log"
	"io/fs"
	"strings"
	"path/filepath"
	"crypto/sha256"
	"github.com/dhowden/tag"
	"chlofisher.com/rosewood/internal/library"
	"chlofisher.com/rosewood/internal/db"
)

var audioExtensions = map[string]bool{
    ".mp3":  true,
    ".flac": true,
    ".ogg":  true,
    ".wav":  true,
    ".m4a":  true,
}

func isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return audioExtensions[ext] // True if ext is a recognised audio extension
}

type FileScanner struct {
	Music *db.MusicStore
}

func (s *FileScanner) Scan(root string) error {
	err := filepath.WalkDir(root, s.handleFile)
	return err
}

func (s *FileScanner) handleFile(path string, d fs.DirEntry, err error) error {
	if err != nil {
		log.Printf("Skipping %s: %v\n", path, err)
		return nil
	}

	if d.IsDir() {
		return nil
	}

	if isAudioFile(path) {
		err := s.handleAudioFile(path)	
		if err != nil {
			log.Printf("Error processing %s: %v\n", path, err)
		}
		return nil
	}

	return nil
}

func (s *FileScanner) handleAudioFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	meta, err := tag.ReadFrom(file)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%x", hashPath(path))
	song := &library.Song{
		ID: id,
		Path: path,
		Title: meta.Title(),
		Artist: meta.Artist(),
		Album: meta.Album(),
	}

	return s.Music.Insert(song)
}

func hashPath(path string) []byte {
	h := sha256.New()
	h.Write([]byte(path))
	return h.Sum(nil)
}

