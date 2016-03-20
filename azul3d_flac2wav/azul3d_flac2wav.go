// flac2wav is a tool which converts FLAC files to WAV files.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"azul3d.org/engine/audio"
	_ "azul3d.org/engine/audio/flac" // Add audio decoder
	"azul3d.org/engine/audio/wav"
)

// flagForce specifies if file overwriting should be forced, when a WAV file of
// the same name already exists.
var flagForce bool

func init() {
	flag.BoolVar(&flagForce, "f", false, "Force overwrite.")
}

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		if err := flac2wav(path); err != nil {
			log.Fatal(err)
		}
	}
}

// flac2wav converts the provided FLAC file to a WAV file.
func flac2wav(path string) error {
	fr, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fr.Close()

	// Open FLAC file.
	dec, _, err := audio.NewDecoder(fr)
	if err != nil {
		return err
	}

	// Create WAV file.
	wavPath := trimExt(path) + ".wav"
	if !flagForce {
		exists, err := fileExists(wavPath)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("the file %q exists already", wavPath)
		}
	}
	fw, err := os.Create(wavPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	// Create WAV encoder.
	enc, err := wav.NewEncoder(fw, dec.Config())
	if err != nil {
		return err
	}
	defer enc.Close()

	// Copy samples from the FLAC decoder to the WAV encoder.
	if _, err = audio.Copy(enc, dec); err != nil {
		return err
	}

	return nil
}

// fileExists reports whether the given file or directory exists or not.
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// trimExt returns filePath without its extension.
func trimExt(filePath string) string {
	ext := path.Ext(filePath)
	return filePath[:len(filePath)-len(ext)]
}
