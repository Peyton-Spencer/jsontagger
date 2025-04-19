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
	
	// Add usage information
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "JSonTagger: A tool for managing JSON tags in Go struct fields\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  jsontagger -file path/to/file.go [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nFeatures:\n")
		fmt.Fprintf(os.Stderr, "  - Converts existing JSON tags between snake_case and camelCase\n")
		fmt.Fprintf(os.Stderr, "  - Generates new JSON tags for fields without them\n")
		fmt.Fprintf(os.Stderr, "  - Preserves tag options like 'omitempty'\n")
		fmt.Fprintf(os.Stderr, "  - Works with fields that have other non-JSON tags\n")
	}
	
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
	jsonTagRegex := regexp.MustCompile(`json:"([^"]+)"`)
	
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
			// Skip fields without names (embedded types)
			if len(field.Names) == 0 {
				continue
			}
			
			fieldName := field.Names[0].Name
			
			// Case 1: Field has no tag - add a new JSON tag
			if field.Tag == nil {
				var newFieldName string
				if config.UseSnake {
					newFieldName = strcase.ToSnake(fieldName)
				} else {
					newFieldName = strcase.ToLowerCamel(fieldName)
				}
				
				// Create a new tag with just the JSON part
				field.Tag = &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("`json:\"%s\"`", newFieldName),
				}
				modified = true
				continue
			}
			
			// Case 2: Field has a tag - check if it contains a JSON part
			tag := field.Tag.Value
			tagStr := strings.Trim(tag, "`")
			
			matches := jsonTagRegex.FindStringSubmatch(tagStr)
			
			// If no JSON tag is present, add one
			if len(matches) < 2 {
				var newFieldName string
				if config.UseSnake {
					newFieldName = strcase.ToSnake(fieldName)
				} else {
					newFieldName = strcase.ToLowerCamel(fieldName)
				}
				
				// Add JSON tag to existing tags
				newTagStr := tagStr + fmt.Sprintf(` json:"%s"`, newFieldName)
				field.Tag.Value = "`" + newTagStr + "`"
				modified = true
				continue
			}
			
			// Case 3: JSON tag exists - modify it as before
			jsonTag := matches[1]
			
			// Split off any options like ,omitempty
			parts := strings.Split(jsonTag, ",")
			existingFieldName := parts[0]
			options := parts[1:]
			
			// Transform the field name according to config
			var newFieldName string
			if config.UseSnake {
				newFieldName = strcase.ToSnake(existingFieldName)
			} else {
				newFieldName = strcase.ToLowerCamel(existingFieldName)
			}
			
			// No change needed
			if newFieldName == existingFieldName {
				continue
			}
			
			// Reconstruct the JSON tag
			newTagParts := []string{newFieldName}
			newTagParts = append(newTagParts, options...)
			newJsonTag := strings.Join(newTagParts, ",")
			
			// Replace the JSON tag in the original tag string
			newTagStr := jsonTagRegex.ReplaceAllString(tagStr, fmt.Sprintf(`json:"%s"`, newJsonTag))
			
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