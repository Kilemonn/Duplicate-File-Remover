package files

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func MergeFileDirs(inputDirs []string, outputDir string, retainModifiedTime bool) (err error) {
	var hashes map[string]bool = make(map[string]bool)
	for _, dir := range inputDirs {
		fmt.Printf("Merging files from dir [%s] into [%s].\n", dir, outputDir)
		err = mergeFileDir(hashes, dir, outputDir, retainModifiedTime)
	}

	return
}

func mergeFileDir(hashes map[string]bool, inputDir string, outputDir string, retainModifiedTime bool) (err error) {
	dirs, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Encountered error when reading directory [%s]. %s.\n", inputDir, err.Error())
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			inputFile := filepath.Join(inputDir, dir.Name())
			hash, bytes, err := getContentHash(inputFile)
			if err != nil {
				fmt.Printf("Failed to generate content hash for file [%s].\n", inputFile)
			}

			if _, exists := hashes[hash]; !exists {
				outputFile := filepath.Join(outputDir, getOutputPath(outputDir, dir.Name()))
				fmt.Printf("Entry for file [%s] (hash %s) does not exist. Creating copy of file to [%s].\n", dir.Name(), hash, outputFile)
				os.WriteFile(outputFile, bytes, os.ModeAppend)
				if retainModifiedTime {
					retainModifiedTimeOfFile(dir, outputFile)
				}
				hashes[hash] = true
			} else {
				fmt.Printf("Entry for file [%s] (hash %s) already exists. File does not need to be copied.\n", inputFile, hash)
			}
		}
	}

	return
}

// Check for colliding file names and make sure we return a filename that exists in the output directory.
// Iterating over file names "test (i).txt" where i will increment for each file with the same name.
func getOutputPath(dir string, file string) string {
	if _, err := os.Stat(filepath.Join(dir, file)); os.IsNotExist(err) {
		return file
	}

	ext := ""
	prefix := ""
	if dot := strings.LastIndex(file, "."); dot != -1 {
		ext = file[dot:]
		prefix = file[:dot]
	} else {
		prefix = file
	}

	for i := 1; ; i++ {
		tempFileName := fmt.Sprintf("%s (%d)%s", prefix, i, ext)
		if _, err := os.Stat(filepath.Join(dir, tempFileName)); os.IsNotExist(err) {
			return tempFileName
		}
	}
}

func retainModifiedTimeOfFile(dir fs.DirEntry, outputFile string) {
	originalFileInfo, err := dir.Info()
	if err != nil {
		fmt.Printf("Failed to retrieve info for file %s. Skipping modification date change for this file.\n", dir.Name())
	} else {
		os.Chtimes(outputFile, time.Time{}, originalFileInfo.ModTime())
	}
}

func getContentHash(filePath string) (hash string, bytes []byte, err error) {
	bytes, err = os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to get content hash for file [%s].\n", filePath)
		return
	}

	hashBytes := sha512.Sum512(bytes)
	hash = hex.EncodeToString(hashBytes[:])
	return
}

func CreateOutputDirectory(outputDir string) (err error) {
	dir, err := filepath.Abs(outputDir)
	if err != nil {
		fmt.Printf("Failed to get absolute path of output directory [%s]. Error: %s.\n", outputDir, err.Error())
	}

	if !filepath.IsAbs(dir) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Failed to retrieve home directory for current user. Error: [%s].\n", err.Error())
			return err
		}
		dir = homeDir + dir
	}

	err = os.Mkdir(dir, fs.FileMode(os.O_APPEND))
	if err != nil && !errors.Is(err, os.ErrExist) {
		fmt.Printf("Failed to create directory [%s]. Error: [%s].\n", outputDir, err.Error())
	} else {
		err = nil
	}
	return
}
