package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/peyton-spencer/caseconv/strcase"
)

type Config struct {
	Filename  string
	UseSnake  bool
	UseCamel  bool
}

func main() {
	config := Config{}
	
	flag.StringVar(&config.Filename, "file", "", "Go file to process")
	flag.BoolVar(&config.UseSnake, "snake", false, "Convert JSON tags to snake_case")
	flag.BoolVar(&config.UseCamel, "camel", false, "Convert JSON tags to camelCase (default)")
	flag.Parse()
	
	if config.Filename == "" {
		fmt.Println("Error: Missing required -file parameter")
		flag.Usage()
		os.Exit(1)
	}
	
	// If neither flag is set, default to camelCase
	if !config.UseSnake && !config.UseCamel {
		config.UseCamel = true
	}
	
	// If both flags are set, prioritize snake case
	if config.UseSnake && config.UseCamel {
		fmt.Println("Warning: Both -snake and -camel flags set. Using snake_case.")
		config.UseCamel = false
	}
	
	if err := processFile(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func processFile(config Config) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, config.Filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %v", err)
	}
	
	modified := false
	
	// Walk the AST to find struct types
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Type == nil {
			return true
		}
		
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok || structType.Fields == nil {
			return true
		}
		
		// Process struct fields
		for _, field := range structType.Fields.List {
			if field.Tag == nil {
				continue
			}
			
			tag := field.Tag.Value
			// Remove backticks and get the raw tag string
			tagStr := strings.Trim(tag, "`")
			
			// Match JSON tags specifically
			re := regexp.MustCompile(`json:"([^"]+)"`)
			matches := re.FindStringSubmatch(tagStr)
			if len(matches) < 2 {
				continue
			}
			
			// Get the current JSON tag value
			jsonTag := matches[1]
			
			// Split off any options like ,omitempty
			parts := strings.Split(jsonTag, ",")
			fieldName := parts[0]
			options := parts[1:]
			
			// Transform the field name according to config
			var newFieldName string
			if config.UseSnake {
				newFieldName = strcase.ToSnake(fieldName)
			} else {
				newFieldName = strcase.ToLowerCamel(fieldName)
			}
			
			// No change needed
			if newFieldName == fieldName {
				continue
			}
			
			// Reconstruct the JSON tag
			newTagParts := []string{newFieldName}
			newTagParts = append(newTagParts, options...)
			newJsonTag := strings.Join(newTagParts, ",")
			
			// Replace the JSON tag in the original tag string
			newTagStr := re.ReplaceAllString(tagStr, fmt.Sprintf(`json:"%s"`, newJsonTag))
			
			// Update the tag in the AST
			field.Tag.Value = "`" + newTagStr + "`"
			modified = true
		}
		
		return true
	})
	
	if !modified {
		fmt.Println("No JSON tags were modified.")
		return nil
	}
	
	// Format the modified AST
	var buf strings.Builder
	if err := format.Node(&buf, fset, node); err != nil {
		return fmt.Errorf("failed to format output: %v", err)
	}
	
	// Write the result back to the file
	if err := os.WriteFile(config.Filename, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}
	
	fmt.Printf("Successfully updated JSON tags in %s\n", config.Filename)
	return nil
}