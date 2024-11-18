package files

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// WithTempDir creates a temp directory and provides it to the provided function as an argument.
// The temp dir is removed automatically after the called func.
func WithTempDir(t *testing.T, testFunc func(dirName string)) {
	dirName, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dirName)

	testFunc(dirName)
}

// CreateFileWithContent creates a file at the provided fileName path with the provided content string
func CreateFileWithContent(t *testing.T, fileName string, content string) {
	require.NoError(t, os.WriteFile(fileName, []byte(content), 0666))
}

// Ensuring merging with files with the same name and same content will only
// copy 1 of the files since duplicates should not be copied into the output dir.
func TestFileMerge_SameNameAndSameContent(t *testing.T) {
	content := "TestFileMerge_SameNameAndSameContent"
	fileName := "text.txt"

	WithTempDir(t, func(outputDir string) {
		WithTempDir(t, func(input1 string) {
			WithTempDir(t, func(input2 string) {
				CreateFileWithContent(t, filepath.Join(input1, fileName), content)
				CreateFileWithContent(t, filepath.Join(input2, fileName), content)

				hash1, _, err := getContentHash(filepath.Join(input1, fileName))
				require.NoError(t, err)
				hash2, _, err := getContentHash(filepath.Join(input2, fileName))
				require.NoError(t, err)
				require.Equal(t, hash1, hash2)

				require.NoError(t, MergeFileDirs([]string{input1, input2}, outputDir, false))
				dirs, err := os.ReadDir(outputDir)
				require.NoError(t, err)
				require.Equal(t, 1, len(dirs))

				_, err = os.Stat(filepath.Join(outputDir, fileName))
				require.NoError(t, err)
			})
		})
	})
}

// Ensure that when two files with a different name and the same content are merged that only the first file
// (determined by the input directory order) is copied to the output dir.
func TestFileMerge_DifferentNameAndSameContent(t *testing.T) {
	content := "TestFileMerge_DifferentNameAndSameContent"
	file1 := "file1.txt"
	file2 := "file2.txt"

	WithTempDir(t, func(outputDir string) {
		WithTempDir(t, func(input1 string) {
			WithTempDir(t, func(input2 string) {
				CreateFileWithContent(t, filepath.Join(input1, file1), content)
				CreateFileWithContent(t, filepath.Join(input2, file2), content)

				hash1, _, err := getContentHash(filepath.Join(input1, file1))
				require.NoError(t, err)
				hash2, _, err := getContentHash(filepath.Join(input2, file2))
				require.NoError(t, err)
				require.Equal(t, hash1, hash2)

				require.NoError(t, MergeFileDirs([]string{input1, input2}, outputDir, false))
				dirs, err := os.ReadDir(outputDir)
				require.NoError(t, err)
				require.Equal(t, 1, len(dirs))

				_, err = os.Stat(filepath.Join(outputDir, file1))
				require.NoError(t, err)

				_, err = os.Stat(filepath.Join(outputDir, file2))
				require.Error(t, err)
			})
		})
	})
}

// Ensure that when we merge files that have different names AND different content that they are copied
// into the output directory.
func TestFileMerge_DifferentNameAndDifferentContent(t *testing.T) {
	contentPrefix := "TestFileMerge_DifferentNameAndDifferentContent"
	file1 := "diff1.txt"
	file2 := "diff2.txt"

	WithTempDir(t, func(outputDir string) {
		WithTempDir(t, func(input1 string) {
			WithTempDir(t, func(input2 string) {
				CreateFileWithContent(t, filepath.Join(input1, file1), contentPrefix+"1")
				CreateFileWithContent(t, filepath.Join(input2, file2), contentPrefix+"2")

				hash1, _, err := getContentHash(filepath.Join(input1, file1))
				require.NoError(t, err)
				hash2, _, err := getContentHash(filepath.Join(input2, file2))
				require.NoError(t, err)
				require.NotEqual(t, hash1, hash2)

				require.NoError(t, MergeFileDirs([]string{input1, input2}, outputDir, false))
				dirs, err := os.ReadDir(outputDir)
				require.NoError(t, err)
				require.Equal(t, 2, len(dirs))

				_, err = os.Stat(filepath.Join(outputDir, file1))
				require.NoError(t, err)
				_, err = os.Stat(filepath.Join(outputDir, file2))
				require.NoError(t, err)
			})
		})
	})
}

// Ensure that when the same file name already exists with the same content that it will copy BOTH versions of the file
// even if there is a conflicting filename.
func TestFileMerge_SameNameDifferentContent(t *testing.T) {
	contentPrefix := "TestFileMerge_SameNameDifferentContent"
	fileName := "test.txt"

	WithTempDir(t, func(outputDir string) {
		WithTempDir(t, func(input1 string) {
			WithTempDir(t, func(input2 string) {
				CreateFileWithContent(t, filepath.Join(input1, fileName), contentPrefix+"1")
				CreateFileWithContent(t, filepath.Join(input2, fileName), contentPrefix+"2")

				hash1, _, err := getContentHash(filepath.Join(input1, fileName))
				require.NoError(t, err)
				hash2, _, err := getContentHash(filepath.Join(input2, fileName))
				require.NoError(t, err)

				require.NotEqual(t, hash1, hash2)

				require.NoError(t, MergeFileDirs([]string{input1, input2}, outputDir, false))
				dirs, err := os.ReadDir(outputDir)
				require.NoError(t, err)
				require.Equal(t, 2, len(dirs))

				_, err = os.Stat(filepath.Join(outputDir, fileName))
				require.NoError(t, err)

				hash3, _, err := getContentHash(filepath.Join(outputDir, fileName))
				require.NoError(t, err)
				require.Equal(t, hash1, hash3)

				duplicateFileName := "test (1).txt"
				_, err = os.Stat(filepath.Join(outputDir, duplicateFileName))
				require.NoError(t, err)

				hash4, _, err := getContentHash(filepath.Join(outputDir, duplicateFileName))
				require.NoError(t, err)
				require.Equal(t, hash2, hash4)
			})
		})
	})
}

// Ensure that when the retain modified flag is set, the copied files retain the same
// modified date as the original.
func TestRetainModifiedTime(t *testing.T) {
	WithTempDir(t, func(outputDir string) {
		WithTempDir(t, func(input1 string) {
			WithTempDir(t, func(input2 string) {
				CreateFileWithContent(t, filepath.Join(input1, "text.txt"), "content")
				original, err := os.Stat(filepath.Join(input1, "text.txt"))
				require.NoError(t, err)

				// Sleep for 2 seconds to force a longer difference in time between the initial file
				// and its copied file
				time.Sleep(2 * time.Second)

				require.NoError(t, MergeFileDirs([]string{input1, input2}, outputDir, true))

				copied, err := os.Stat(filepath.Join(outputDir, "text.txt"))
				require.NoError(t, err)

				require.Equal(t, original.ModTime(), copied.ModTime())
			})
		})
	})
}

// Ensure that when the retain modified flag is set, the copied files does not have its modified date changed
// and is after the original files modified date.
func TestDontRetainModifiedTime(t *testing.T) {
	WithTempDir(t, func(outputDir string) {
		WithTempDir(t, func(input1 string) {
			WithTempDir(t, func(input2 string) {
				CreateFileWithContent(t, filepath.Join(input1, "text.txt"), "content")
				original, err := os.Stat(filepath.Join(input1, "text.txt"))
				require.NoError(t, err)

				// Sleep for 2 seconds to force a longer difference in time between the initial file
				// and its copied file
				time.Sleep(2 * time.Second)

				require.NoError(t, MergeFileDirs([]string{input1, input2}, outputDir, false))

				copied, err := os.Stat(filepath.Join(outputDir, "text.txt"))
				require.NoError(t, err)

				require.Less(t, original.ModTime(), copied.ModTime())
			})
		})
	})
}
