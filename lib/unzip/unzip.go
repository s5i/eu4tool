package unzip

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
)

// Meta opens a save zip file and extracts contents of meta file
// into a byte slice.
func Meta(path string) ([]byte, error) {
	return unzip(path, "meta")
}

// Gamestate opens a save zip file and extracts contents of gamestate
// file into a byte slice.
func Gamestate(path string) ([]byte, error) {
	return unzip(path, "gamestate")
}

// AI opens a save zip file and extracts contents of ai file
// into a byte slice.
func AI(path string) ([]byte, error) {
	return unzip(path, "ai")
}

func unzip(path, innerFile string) ([]byte, error) {
	rc, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("zip.OpenReader failed: %v", err)
	}
	defer rc.Close()

	for _, f := range rc.Reader.File {
		if f.FileHeader.Name != innerFile {
			continue
		}
		frc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("f.Open failed: %v", err)
		}
		defer frc.Close()

		b, err := ioutil.ReadAll(frc)
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadAll failed: %v", err)
		}

		return b, nil
	}

	return nil, fmt.Errorf("%s not found in zip", innerFile)
}
