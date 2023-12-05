package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	gitignore "github.com/sabhiram/go-gitignore"
)

func main() {
	flag.Parse()
	patterns := flag.Args()

	if len(patterns) == 0 {
		// If no pattern is provided, default to "."
		patterns = append(patterns, ".")
	}

	if err := concopy(".", patterns); err != nil {
		fmt.Println("Error:", err)
	}
}

func concopy(rootDir string, patterns []string) error {
	gitIgnore, err := gitignore.CompileIgnoreFile(filepath.Join(rootDir, ".gitignore"))
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || gitIgnore.MatchesPath(path) || !shouldInclude(path, patterns, rootDir) || strings.Contains(path, "/.git/") || strings.HasSuffix(path, "/.git") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Println(path)
		buffer.WriteString(path + ":\n")
		buffer.Write(data)
		buffer.WriteString("\n\n")

		return nil
	})

	if err != nil {
		return err
	}

	return clipboard.WriteAll(buffer.String())
}

func shouldInclude(filePath string, patterns []string, rootDir string) bool {
	relPath, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		return false
	}

	if len(patterns) == 0 {
		return true // If no patterns specified, include all files
	}

	for _, pattern := range patterns {
		if pattern == "." {
			// Special case for '.', include all files
			return true
		}

		// Check if the pattern is a directory
		patternInfo, err := os.Stat(pattern)
		if err == nil && patternInfo.IsDir() {
			if strings.HasPrefix(relPath, pattern) {
				return true
			}
		}

		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			fmt.Println("Invalid pattern:", err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
