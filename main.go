package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	inputPath     string
	outputPath    string
	yearNamedPath string

	yearPathNames = map[string]string{}

	tmdbRegex = regexp.MustCompile(`{tmdb-\d*}`)
	yearRegex = regexp.MustCompile(`\(\d*\)`)
)

func main() {
	flag.StringVar(&inputPath, "input", "", "The path a folder of media named using plex naming")
	flag.StringVar(&outputPath, "output", "", "The path to output renamed files to")
	flag.StringVar(&yearNamedPath, "yearPath", "", "The path to a folder of folders with the right years to compare against")

	flag.Parse()

	if inputPath == "" {
		log.Println("Input path must not be blank")
	}

	if outputPath == "" {
		log.Println("Output path must not be blank")
	}

	if yearNamedPath != "" {
		slog.Debug("Parsing years folder")

		files, err := os.ReadDir(yearNamedPath)
		if err != nil {
			slog.Error("Failed to read input directory", "err", err)
			return
		}

		for _, file := range files {
			if !file.IsDir() || strings.HasPrefix(file.Name(), ".") {
				continue
			}

			name := file.Name()
			slog.Debug("Processing Year Folder", "folder", name)

			yearParts := yearRegex.FindAllString(name, 1)
			if len(yearParts) == 1 {
				year := yearParts[0]
				name = strings.Replace(name, year, "", -1)
				name = strings.TrimSpace(name) // Clean up whitespace
				yearPathNames[name] = year
				slog.Debug("Add year", "folder", name, "year", year)
			} else {
				continue
			}
		}
	}

	slog.Debug("Processing folders", "path", inputPath)

	files, err := os.ReadDir(inputPath)
	if err != nil {
		slog.Error("Failed to read input directory", "err", err)
		return
	}

	slog.Debug("Found Folders", "folders", len(files))

	for _, file := range files {
		if !file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		slog.Debug("Processing folder", "folder", file.Name())

		name := file.Name()
		// Remove tmdb suffix
		name = tmdbRegex.ReplaceAllString(name, "")

		// Trim whitespace
		name = strings.TrimSpace(name)

		// Add year if possible
		if yearNamedPath != "" && len(yearPathNames) > 0 {
			year, ok := yearPathNames[name]
			if ok {
				name = name + " " + year
			}
		}

		dirInput := filepath.Join(inputPath, file.Name())
		dirOutput := filepath.Join(outputPath, name)

		err := os.Rename(dirInput, dirOutput)
		if err != nil {
			slog.Error("Error renaming folder", "from", dirInput, "to", dirOutput)
			continue
		}

		slog.Info("Renaming folder", "from", file.Name(), "to", name)
	}
}
