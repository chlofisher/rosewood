package scanner

import (
	"os"
	"log"
	"io/fs"
	"strings"
	"path/filepath"
	"github.com/dhowden/tag"
	"chlofisher.com/rosewood/internal/metadata"
	"chlofisher.com/rosewood/internal/server/db"
)

var audioExtensions = map[string]struct{}{
    ".mp3":  {},
    ".flac": {},
    ".ogg":  {},
    ".wav":  {},
    ".m4a":  {},
}

func isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	// Return true if ext is in the set of audioextensions
	if _, ok := audioExtensions[ext]; ok {
		return true
	}
	return false
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
		log.Printf("Skipping %s: %v", path, err)
		return nil
	}

	if d.IsDir() {
		return nil
	}

	if isAudioFile(path) {
		err := s.handleAudioFile(path)	
		if err != nil {
			log.Printf("Error processing %s: %v; skipping.", path, err)
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

	song := &metadata.Song{
		Path: path,
		Title: meta.Title(),
		Artist: meta.Artist(),
		Album: meta.Album(),
	}

	return s.Music.Insert(song)
}

