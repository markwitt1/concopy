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
	// Use patterns from .concopyuse if it exists, otherwise from command-line arguments
	patterns, err := readConcopyUsePatterns(".")
	if err != nil {
		fmt.Println("Error reading .concopyuse:", err)
		return
	}
	if len(patterns) == 0 {
		flag.Parse()
		patterns = flag.Args()
		if len(patterns) == 0 {
			patterns = append(patterns, ".") // Default to "."
		}
	}

	if err := concopy(".", patterns); err != nil {
		fmt.Println("Error:", err)
	}
}

func concopy(rootDir string, patterns []string) error {
	var buffer bytes.Buffer
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the .git directory
		if info.IsDir() && (info.Name() == ".git" || strings.Contains(path, "/.git/")) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			gitIgnore, err := loadGitIgnore(path)
			if err != nil {
				fmt.Println("Error loading .gitignore:", err)
				return nil
			}
			if gitIgnore != nil && gitIgnore.MatchesPath(path) {
				return filepath.SkipDir
			}
		} else if shouldInclude(path, patterns, rootDir) && !isIgnoredByGitignore(path, rootDir) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			fmt.Println(path)
			buffer.WriteString(path + ":\n")
			buffer.Write(data)
			buffer.WriteString("\n\n")
		}
		return nil
	})

	if err != nil {
		return err
	}

	return clipboard.WriteAll(buffer.String())
}

func loadGitIgnore(dir string) (*gitignore.GitIgnore, error) {
	gitignorePath := filepath.Join(dir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return nil, nil
	}
	return gitignore.CompileIgnoreFile(gitignorePath)
}

func isIgnoredByGitignore(filePath, rootDir string) bool {
	for dir := filePath; dir != rootDir; dir = filepath.Dir(dir) {
		gitIgnore, _ := loadGitIgnore(filepath.Dir(dir))
		if gitIgnore != nil && gitIgnore.MatchesPath(filePath) {
			return true
		}
	}
	return false
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

// readConcopyUsePatterns reads the .concopyuse file and returns the patterns
func readConcopyUsePatterns(rootDir string) ([]string, error) {
	concopyUsePath := filepath.Join(rootDir, ".concopyuse")
	content, err := ioutil.ReadFile(concopyUsePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No .concopyuse file, not an error
		}
		return nil, err
	}

	patterns := strings.Split(strings.TrimSpace(string(content)), "\n")
	for i, pattern := range patterns {
		patterns[i] = strings.TrimSpace(pattern)
	}
	return patterns, nil
}
