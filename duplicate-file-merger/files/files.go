package files

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func MergeFileDirs(inputDirs []string, outputDir string) (err error) {
	var hashes map[string]bool = make(map[string]bool)
	for _, dir := range inputDirs {
		fmt.Printf("Merging files from dir [%s] into [%s].\n", dir, outputDir)
		err = mergeFileDir(hashes, dir, outputDir)
	}

	return
}

func mergeFileDir(hashes map[string]bool, inputDir string, outputDir string) (err error) {
	dirs, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Encountered error when reading directory [%s]. %s.\n", inputDir, err.Error())
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			inputFile := inputDir + "/" + dir.Name()
			hash, bytes, err := getContentHash(inputFile)
			if err != nil {
				fmt.Printf("Failed to generate content hash for file [%s].\n", inputFile)
			}

			if _, exists := hashes[hash]; !exists {
				outputFile := outputDir + "/" + dir.Name()
				fmt.Printf("Entry for file [%s] (hash %s) does not exist. Moving file to %s.\n", dir.Name(), hash, outputFile)
				os.WriteFile(outputFile, bytes, os.ModeAppend)
				hashes[hash] = true
			} else {
				fmt.Printf("Entry for file [%s] (hash %s) already exists. File does not need to be copied.\n", inputFile, hash)
			}
		}
	}

	return
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
