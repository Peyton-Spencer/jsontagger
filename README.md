# Go JSON Tagger

A simple CLI tool to transform JSON struct tags in Go files between snake_case and camelCase formats. This tool uses the [caseconv](https://github.com/peyton-spencer/caseconv) package for case conversion. It can both modify existing JSON tags and generate new tags for fields that don't have them.

## Usage

```bash
# Convert JSON tags to camelCase (default)
jsontagger -file path/to/your/file.go

# Convert JSON tags to snake_case
jsontagger -file path/to/your/file.go -snake

# Explicitly specify camelCase (same as default)
jsontagger -file path/to/your/file.go -camel
```

## Go Tool
use in your project
```
go get -tool github.com/peyton-spencer/jsontagger
```

```
go tool jsontagger -file path/to/your/file.go
```

## Examples

### Modifying existing JSON tags

#### Original struct with snake_case tags:

```go
type User struct {
    UserID      int    `json:"user_id"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name,omitempty"`
}
```

#### After running with `-camel`:

```go
type User struct {
    UserID      int    `json:"userId"`
    FirstName   string `json:"firstName"`
    LastName    string `json:"lastName,omitempty"`
}
```

#### Original struct with camelCase tags:

```go
type Product struct {
    ProductID   int     `json:"productId"`
    Name        string  `json:"name"`
    UnitPrice   float64 `json:"unitPrice,omitempty"`
}
```

#### After running with `-snake`:

```go
type Product struct {
    ProductID   int     `json:"product_id"`
    Name        string  `json:"name"`
    UnitPrice   float64 `json:"unit_price,omitempty"`
}
```

### Generating new JSON tags

#### Original struct without any tags:

```go
type Person struct {
    ID          int
    FirstName   string
    LastName    string
    DateOfBirth string
}
```

#### After running with `-camel`:

```go
type Person struct {
    ID          int    `json:"id"`
    FirstName   string `json:"firstName"`
    LastName    string `json:"lastName"`
    DateOfBirth string `json:"dateOfBirth"`
}
```

#### After running with `-snake`:

```go
type Person struct {
    ID          int    `json:"id"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    DateOfBirth string `json:"date_of_birth"`
}
```

### Working with mixed and non-JSON tags

#### Original struct with mixed tags:

```go
type Project struct {
    ID          int    `db:"project_id"`
    Name        string `validate:"required"`
    Description string
}
```

#### After running with `-camel`:

```go
type Project struct {
    ID          int    `db:"project_id" json:"id"`
    Name        string `validate:"required" json:"name"`
    Description string `json:"description"`
}
```

## Build

```bash
go build
```

## Dependencies

- [github.com/peyton-spencer/caseconv](https://github.com/peyton-spencer/caseconv) - For case conversion functions (using strcase package)

## License

MIT