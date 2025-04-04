package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type folder struct {
	path     string
	decimal  int
	restName string
}

type Config struct {
	OldPrefix   string
	NewPrefix   string
	StartFrom   int
	Dir        string
	DryRun     bool
	DigitCount int // Number of digits after the decimal point
}

// renameDirectories performs the actual renaming logic
func renameDirectories(cfg Config) error {
	// Validate inputs
	if cfg.DigitCount < 1 {
		return fmt.Errorf("digit count must be at least 1")
	}

	// If we're renumbering, validate that we won't exceed the maximum
	if cfg.StartFrom > 0 {
		// Count how many folders we need to renumber
		folderCount := 0
		err := filepath.Walk(cfg.Dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() || path == cfg.Dir {
				return nil
			}
			baseName := filepath.Base(path)
			parts := strings.SplitN(baseName, " ", 2)
			if len(parts) != 2 {
				return nil
			}
			numberParts := strings.Split(parts[0], ".")
			if len(numberParts) != 2 || numberParts[0] != cfg.OldPrefix {
				return nil
			}
			folderCount++
			return nil
		})
		if err != nil {
			return fmt.Errorf("error counting folders: %v", err)
		}

		// Check if the last number would exceed the maximum
		lastNumber := cfg.StartFrom + folderCount - 1
		maxNumber := int(math.Pow10(cfg.DigitCount)) - 1
		if lastNumber > maxNumber {
			return fmt.Errorf("renumbering would exceed xx.%s (last number would be %d)", strings.Repeat("9", cfg.DigitCount), lastNumber)
		}
	}

	// Collect all matching folders first
	var folders []folder
	err := filepath.Walk(cfg.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Skip the root directory
		if path == cfg.Dir {
			return nil
		}

		baseName := filepath.Base(path)
		parts := strings.SplitN(baseName, " ", 2)
		if len(parts) == 2 {
			numberParts := strings.Split(parts[0], ".")
			if len(numberParts) == 2 && numberParts[0] == cfg.OldPrefix {
				decimal, err := strconv.Atoi(numberParts[1])
				if err != nil {
					return nil
				}

				folders = append(folders, folder{path: path, decimal: decimal, restName: parts[1]})
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %v", err)
	}

	// Sort folders by their current decimal number
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].decimal < folders[j].decimal
	})

	// Use oldPrefix as newPrefix if not specified but startFrom is
	if cfg.StartFrom > 0 && cfg.NewPrefix == "" {
		cfg.NewPrefix = cfg.OldPrefix
	}

	// Renumber and rename the folders
	newDecimal := cfg.StartFrom
	if newDecimal == 0 {
		// If no start number provided, keep original decimal numbers
		for _, f := range folders {
			format := fmt.Sprintf("%%s.%%0%dd %%s", cfg.DigitCount)
			newName := fmt.Sprintf(format, cfg.NewPrefix, f.decimal, f.restName)
			newPath := filepath.Join(filepath.Dir(f.path), newName)

			// Skip if the source and target paths are identical
			if f.path == newPath {
				fmt.Printf("Skipping %s (already has correct name)\n", f.path)
				continue
			}

			if cfg.DryRun {
				fmt.Printf("Would rename: %s -> %s\n", f.path, newPath)
			} else {
				fmt.Printf("Renaming: %s -> %s\n", f.path, newPath)
				if err := os.Rename(f.path, newPath); err != nil {
					fmt.Printf("Error renaming %s: %v\n", f.path, err)
				}
			}
		}
	} else {
		// Renumber starting from the specified number
		for _, f := range folders {
			format := fmt.Sprintf("%%s.%%0%dd %%s", cfg.DigitCount)
			newName := fmt.Sprintf(format, cfg.NewPrefix, newDecimal, f.restName)
			newPath := filepath.Join(filepath.Dir(f.path), newName)

			// Skip if the source and target paths are identical
			if f.path == newPath {
				fmt.Printf("Skipping %s (already has correct name)\n", f.path)
				newDecimal++
				continue
			}

			if cfg.DryRun {
				fmt.Printf("Would rename: %s -> %s\n", f.path, newPath)
			} else {
				fmt.Printf("Renaming: %s -> %s\n", f.path, newPath)
				if err := os.Rename(f.path, newPath); err != nil {
					fmt.Printf("Error renaming %s: %v\n", f.path, err)
				}
			}
			newDecimal++
		}
	}

	return nil
}

func main() {
	var oldPrefix string
	var newPrefix string
	var startFrom int
	var dir string
	var dryRun bool
	var digitCount int

	flag.StringVar(&oldPrefix, "from", "", "Original prefix number (e.g., '10')")
	flag.StringVar(&newPrefix, "to", "", "New prefix number (e.g., '20')")
	flag.IntVar(&startFrom, "start", 0, "Start renumbering from this number")
	flag.StringVar(&dir, "dir", ".", "Directory to process")
	flag.BoolVar(&dryRun, "dry-run", false, "Preview changes without making them")
	flag.IntVar(&digitCount, "digits", 2, "Number of digits after the decimal point (default: 2)")
	flag.Parse()

	if oldPrefix == "" {
		fmt.Println("Please provide the -from prefix")
		flag.Usage()
		os.Exit(1)
	}

	// If -start is not provided and -to is not provided, show error
	if startFrom == 0 && newPrefix == "" {
		fmt.Println("Please provide the -to prefix when not using -start")
		flag.Usage()
		os.Exit(1)
	}

	config := Config{
		OldPrefix:   oldPrefix,
		NewPrefix:   newPrefix,
		StartFrom:   startFrom,
		Dir:        dir,
		DryRun:     dryRun,
		DigitCount: digitCount,
	}

	if err := renameDirectories(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
