package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Test case 1: Template with data
	dataFile := filepath.Join(tmpDir, "data.json")
	templateFile := filepath.Join(tmpDir, "template.tpl")

	dataContent := `{"name": "World", "count": 42}`
	templateContent := `Hello {{.name}}! Count: {{.count}}`

	if err := os.WriteFile(dataFile, []byte(dataContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(templateFile, []byte(templateContent), 0o644); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	err := renderTemplate(dataFile, templateFile, &output)
	if err != nil {
		t.Fatalf("renderTemplate failed: %v", err)
	}

	expected := "Hello World! Count: 42"
	if output.String() != expected {
		t.Errorf("Expected %q, got %q", expected, output.String())
	}
}

func TestRenderTemplateNoData(t *testing.T) {
	tmpDir := t.TempDir()

	// Test case 2: Template without data
	templateFile := filepath.Join(tmpDir, "simple.tpl")
	templateContent := `Hello from template!`

	if err := os.WriteFile(templateFile, []byte(templateContent), 0o644); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	err := renderTemplate("", templateFile, &output)
	if err != nil {
		t.Fatalf("renderTemplate failed: %v", err)
	}

	expected := "Hello from template!"
	if output.String() != expected {
		t.Errorf("Expected %q, got %q", expected, output.String())
	}
}

func TestLoadData(t *testing.T) {
	tmpDir := t.TempDir()

	// Test case 1: Valid JSON file
	dataFile := filepath.Join(tmpDir, "valid.json")
	content := `{"key": "value", "number": 123}`
	if err := os.WriteFile(dataFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := loadData(dataFile)
	if err != nil {
		t.Fatalf("loadData failed: %v", err)
	}

	if data["key"] != "value" {
		t.Errorf("Expected key=value, got %v", data["key"])
	}
	if data["number"] != float64(123) {
		t.Errorf("Expected number=123, got %v", data["number"])
	}
}

func TestLoadDataEmpty(t *testing.T) {
	// Test case 2: Empty filename returns empty map
	data, err := loadData("")
	if err != nil {
		t.Fatalf("loadData with empty filename failed: %v", err)
	}

	if len(data) != 0 {
		t.Errorf("Expected empty map, got %v", data)
	}
}

func TestLoadDataInvalidFile(t *testing.T) {
	// Test case 3: Non-existent file
	_, err := loadData("nonexistent.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadDataInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// Test case 4: Invalid JSON
	invalidFile := filepath.Join(tmpDir, "invalid.json")
	content := `{"invalid": json}`
	if err := os.WriteFile(invalidFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := loadData(invalidFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
