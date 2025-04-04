package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func createTestDirs(t *testing.T, baseDir string, names []string) {
	for _, name := range names {
		err := os.Mkdir(filepath.Join(baseDir, name), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", name, err)
		}
	}
}

func getDirNames(t *testing.T, dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

func TestRenameWithNewPrefix(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"10.01 Projects",
		"10.02 Documents",
		"10.03 Archive",
	}
	createTestDirs(t, tempDir, testDirs)

	// Run the program
	cfg := Config{
		OldPrefix:   "10",
		NewPrefix:   "20",
		Dir:        tempDir,
		DigitCount: 2,
	}

	err = renameDirectories(cfg)
	if err != nil {
		t.Fatalf("Failed to rename directories: %v", err)
	}

	// Check results
	expected := []string{
		"20.01 Projects",
		"20.02 Documents",
		"20.03 Archive",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestRenumberWithSamePrefix(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"10.01 Projects",
		"10.02 Documents",
		"10.03 Archive",
	}
	createTestDirs(t, tempDir, testDirs)

	// Run the program with renumbering
	cfg := Config{
		OldPrefix:   "10",
		StartFrom:   5,
		Dir:        tempDir,
		DigitCount: 2,
	}

	err = renameDirectories(cfg)
	if err != nil {
		t.Fatalf("Failed to rename directories: %v", err)
	}

	// Check results
	expected := []string{
		"10.05 Projects",
		"10.06 Documents",
		"10.07 Archive",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestRenumberWithGaps(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories with gaps in numbers
	testDirs := []string{
		"20.01 Some folder",
		"20.02 Another folder",
		"20.10 Yet another folder",
	}
	createTestDirs(t, tempDir, testDirs)

	// Run the program with renumbering
	cfg := Config{
		OldPrefix:   "20",
		StartFrom:   14,
		Dir:        tempDir,
		DigitCount: 2,
	}

	err = renameDirectories(cfg)
	if err != nil {
		t.Fatalf("Failed to rename directories: %v", err)
	}

	// Check results
	expected := []string{
		"20.14 Some folder",
		"20.15 Another folder",
		"20.16 Yet another folder",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestSkipIdenticalPaths(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"20.14 Some folder",
		"20.15 Another folder",
		"20.16 Yet another folder",
	}
	createTestDirs(t, tempDir, testDirs)

	// Try to renumber with the same numbers
	cfg := Config{
		OldPrefix:   "20",
		StartFrom:   14,
		Dir:        tempDir,
		DigitCount: 2,
	}

	err = renameDirectories(cfg)
	if err != nil {
		t.Fatalf("Failed to rename directories: %v", err)
	}

	// Check results - should be exactly the same as input
	expected := []string{
		"20.14 Some folder",
		"20.15 Another folder",
		"20.16 Yet another folder",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestPreventExceeding99(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"20.01 First",
		"20.02 Second",
		"20.03 Third",
	}
	createTestDirs(t, tempDir, testDirs)

	// Try to renumber starting from 98
	cfg := Config{
		OldPrefix:   "20",
		StartFrom:   98,
		Dir:        tempDir,
		DigitCount: 2,
	}

	err = renameDirectories(cfg)
	if err == nil {
		t.Error("Expected an error when renumbering would exceed xx.99, but got none")
	} else if !strings.Contains(err.Error(), "would exceed xx.99") {
		t.Errorf("Expected error about exceeding xx.99, got: %v", err)
	}

	// Verify that no files were renamed
	expected := []string{
		"20.01 First",
		"20.02 Second",
		"20.03 Third",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestLongerDecimalNumbers(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"20.01 First",
		"20.02 Second",
		"20.10 Third",
	}
	createTestDirs(t, tempDir, testDirs)

	// Try to renumber with 4 digits
	cfg := Config{
		OldPrefix:   "20",
		NewPrefix:   "90",
		StartFrom:   1,
		Dir:        tempDir,
		DigitCount: 4,
	}

	err = renameDirectories(cfg)
	if err != nil {
		t.Fatalf("Failed to rename directories: %v", err)
	}

	// Check results
	expected := []string{
		"90.0001 First",
		"90.0002 Second",
		"90.0003 Third",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestLongerDecimalNumbersValidation(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"20.0001 First",
		"20.0002 Second",
		"20.0003 Third",
	}
	createTestDirs(t, tempDir, testDirs)

	// Try to renumber with a number that would exceed the maximum
	cfg := Config{
		OldPrefix:   "20",
		StartFrom:   9998,
		Dir:        tempDir,
		DigitCount: 4,
	}

	err = renameDirectories(cfg)
	if err == nil {
		t.Error("Expected an error when renumbering would exceed the maximum, but got none")
	} else if !strings.Contains(err.Error(), "would exceed xx.9999") {
		t.Errorf("Expected error about exceeding xx.9999, got: %v", err)
	}

	// Verify that no files were renamed
	expected := []string{
		"20.0001 First",
		"20.0002 Second",
		"20.0003 Third",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}

func TestDryRun(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "johnny-decimal-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directories
	testDirs := []string{
		"10.01 Projects",
		"10.02 Documents",
		"10.03 Archive",
	}
	createTestDirs(t, tempDir, testDirs)

	// Run the program with dry-run
	cfg := Config{
		OldPrefix:   "10",
		NewPrefix:   "20",
		Dir:        tempDir,
		DryRun:     true,
		DigitCount: 2,
	}

	err = renameDirectories(cfg)
	if err != nil {
		t.Fatalf("Failed to rename directories: %v", err)
	}

	// Check that no files were actually renamed
	expected := []string{
		"10.01 Projects",
		"10.02 Documents",
		"10.03 Archive",
	}
	actual := getDirNames(t, tempDir)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d directories, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if i >= len(actual) {
			t.Errorf("Missing expected directory: %s", expected[i])
			continue
		}
		if expected[i] != actual[i] {
			t.Errorf("Expected directory %s, got %s", expected[i], actual[i])
		}
	}
}
