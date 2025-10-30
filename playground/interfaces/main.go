package main

import (
	"fmt"
	"strings"
)

// ============================================
// ANTI-PATTERN: Single factory returning interface
// ============================================

type Storage interface {
	Save(data string) error
	Load() (string, error)
}

type DiskStorage struct {
	path string
}

func (d *DiskStorage) Save(data string) error {
	fmt.Printf("Saving to disk: %s\n", d.path)
	return nil
}

func (d *DiskStorage) Load() (string, error) {
	return "disk data", nil
}

type MemoryStorage struct {
	cache map[string]string
}

func (m *MemoryStorage) Save(data string) error {
	fmt.Println("Saving to memory")
	return nil
}

func (m *MemoryStorage) Load() (string, error) {
	return "memory data", nil
}

// BAD: Single factory that returns interface based on parameter
func NewStorage(storageType string) Storage {
	if storageType == "disk" {
		return &DiskStorage{path: "/tmp/data"}
	}
	return &MemoryStorage{cache: make(map[string]string)}
}

// ============================================
// RECOMMENDED: Separate factories returning structs
// ============================================

// GOOD: Each factory returns its own concrete type
func NewDiskStorage(path string) *DiskStorage {
	return &DiskStorage{path: path}
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{cache: make(map[string]string)}
}

// Business logic still accepts interfaces (from Part 1!)
func ProcessData(storage Storage, data string) error {
	return storage.Save(data)
}

// ============================================
// UNAVOIDABLE CASE: Parser returning different token types
// ============================================

// When parsing, you genuinely don't know what type you'll get
type Token interface {
	Type() string
	Value() string
}

type NumberToken struct {
	val string
}

func (n *NumberToken) Type() string  { return "NUMBER" }
func (n *NumberToken) Value() string { return n.val }

type StringToken struct {
	val string
}

func (s *StringToken) Type() string  { return "STRING" }
func (s *StringToken) Value() string { return s.val }

type OperatorToken struct {
	val string
}

func (o *OperatorToken) Type() string  { return "OPERATOR" }
func (o *OperatorToken) Value() string { return o.val }

// UNAVOIDABLE: Parser must return interface because the type varies
func ParseNextToken(input string) Token {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}

	// We don't know what we'll return until we parse
	if input[0] >= '0' && input[0] <= '9' {
		return &NumberToken{val: string(input[0])}
	}
	if input[0] == '"' {
		return &StringToken{val: "hello"}
	}
	return &OperatorToken{val: string(input[0])}
}

// ============================================
// ERROR EXCEPTION: Errors return interface
// ============================================

type NotFoundError struct {
	resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("resource not found: %s", e.resource)
}

type ValidationError struct {
	field   string
	message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.field, e.message)
}

// Functions return error interface, but different concrete types
func FindUser(id string) (*UserData, error) {
	if id == "" {
		// Returns ValidationError
		return nil, &ValidationError{field: "id", message: "cannot be empty"}
	}
	if id == "999" {
		// Returns NotFoundError
		return nil, &NotFoundError{resource: "user"}
	}
	// Returns nil (which satisfies error interface)
	return &UserData{ID: id, Name: "John"}, nil
}

type UserData struct {
	ID   string
	Name string
}

// ============================================
// DEMONSTRATION
// ============================================

func main() {
	fmt.Println("=== Anti-Pattern: Single Factory Returning Interface ===")
	storage1 := NewStorage("disk")
	storage2 := NewStorage("memory")
	fmt.Printf("disk path of storage1: %s\n", storage1.(*DiskStorage).path) // Hidden behind interface)
	fmt.Printf("Type 1: %T\n", storage1)                                    // Hidden behind interface
	fmt.Printf("Type 2: %T\n", storage2)                                    // Hidden behind interface
	fmt.Println("Problem: Caller can't access concrete type features")

	fmt.Println("\n=== Recommended: Separate Factories Returning Structs ===")
	disk := NewDiskStorage("/var/data")
	memory := NewMemoryStorage()
	fmt.Printf("Disk type: %T (concrete!)\n", disk)
	fmt.Printf("Memory type: %T (concrete!)\n", memory)
	fmt.Printf("Disk path: %s (can access specific fields!)\n", disk.path)

	// But we can still pass them to functions that accept interfaces
	ProcessData(disk, "some data")
	ProcessData(memory, "other data")

	fmt.Println("\n=== Unavoidable Case: Parser ===")
	fmt.Println("Parser MUST return interface - we don't know the type until runtime:")
	token1 := ParseNextToken("5")
	token2 := ParseNextToken(`"hello"`)
	token3 := ParseNextToken("+")
	fmt.Printf("Token 1: %s = %s\n", token1.Type(), token1.Value())
	fmt.Printf("Token 2: %s = %s\n", token2.Type(), token2.Value())
	fmt.Printf("Token 3: %s = %s\n", token3.Type(), token3.Value())

	fmt.Println("\n=== Error Exception: Always Returns Interface ===")
	// error is an interface, but different concrete types are returned
	_, err1 := FindUser("")
	_, err2 := FindUser("999")
	_, err3 := FindUser("123")

	fmt.Printf("Error 1: %v (type: %T)\n", err1, err1)
	fmt.Printf("Error 2: %v (type: %T)\n", err2, err2)
	fmt.Printf("Error 3: %v (type: %T)\n", err3, err3)

	fmt.Println("\n=== Summary ===")
	fmt.Println("✓ Prefer: Separate factories returning concrete structs")
	fmt.Println("✗ Avoid: Single factory with type parameter returning interface")
	fmt.Println("⚠ Exceptions: Parsers (runtime-determined types) and errors")
}
