package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "gojsontagger_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file with mixed case JSON tags
	testFile := filepath.Join(tempDir, "test.go")
	testContent := `package test

type Person struct {
	ID        int    ` + "`json:\"person_id\"`" + `
	FirstName string ` + "`json:\"firstName\"`" + `
	LastName  string ` + "`json:\"lastName,omitempty\"`" + `
	Address   string ` + "`json:\"home_address\"`" + `
}
`
	if err := ioutil.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test camelCase conversion
	camelConfig := Config{
		Filename: testFile,
		UseCamel: true,
	}
	if err := processFile(camelConfig); err != nil {
		t.Fatalf("processFile failed with camelCase: %v", err)
	}

	// Read the modified file
	camelContent, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}

	// Verify the tags were converted to camelCase
	expectedCamel := `package test

type Person struct {
	ID        int    ` + "`json:\"personId\"`" + `
	FirstName string ` + "`json:\"firstName\"`" + `
	LastName  string ` + "`json:\"lastName,omitempty\"`" + `
	Address   string ` + "`json:\"homeAddress\"`" + `
}
`
	if string(camelContent) != expectedCamel {
		t.Errorf("Camel case conversion failed.\nExpected:\n%s\nGot:\n%s", expectedCamel, string(camelContent))
	}

	// Now test snake_case conversion
	snakeConfig := Config{
		Filename: testFile,
		UseSnake: true,
	}
	if err := processFile(snakeConfig); err != nil {
		t.Fatalf("processFile failed with snake_case: %v", err)
	}

	// Read the modified file
	snakeContent, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}

	// Verify the tags were converted to snake_case
	expectedSnake := `package test

type Person struct {
	ID        int    ` + "`json:\"person_id\"`" + `
	FirstName string ` + "`json:\"first_name\"`" + `
	LastName  string ` + "`json:\"last_name,omitempty\"`" + `
	Address   string ` + "`json:\"home_address\"`" + `
}
`
	if string(snakeContent) != expectedSnake {
		t.Errorf("Snake case conversion failed.\nExpected:\n%s\nGot:\n%s", expectedSnake, string(snakeContent))
	}
}