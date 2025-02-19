package yamledit

import (
	"reflect"
	"testing"
)

// Person is a sample struct used for testing.
type Person struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

const sampleYAML = `
person:
  name: John Doe
  age: 30
settings:
  theme: dark
  notifications: true
`

func TestRead(t *testing.T) {
	// Test reading a full struct.
	var p Person
	if err := Read([]byte(sampleYAML), "person", &p); err != nil {
		t.Fatalf("Failed to read person node: %v", err)
	}
	if p.Name != "John Doe" || p.Age != 30 {
		t.Errorf("Unexpected person struct: %+v", p)
	}

	// Test reading a simple string.
	var name string
	if err := Read([]byte(sampleYAML), "person.name", &name); err != nil {
		t.Fatalf("Failed to read person.name: %v", err)
	}
	if name != "John Doe" {
		t.Errorf("Expected name to be 'John Doe', got %q", name)
	}

	// Test reading a boolean.
	var notifications bool
	if err := Read([]byte(sampleYAML), "settings.notifications", &notifications); err != nil {
		t.Fatalf("Failed to read settings.notifications: %v", err)
	}
	if notifications != true {
		t.Errorf("Expected notifications to be true, got %v", notifications)
	}
}

func TestUpdate(t *testing.T) {
	// Update person.age to 35.
	updatedYAML, err := Update([]byte(sampleYAML), "person.age", 35)
	if err != nil {
		t.Fatalf("Failed to update person.age: %v", err)
	}

	// Verify updated age.
	var age int
	if err = Read(updatedYAML, "person.age", &age); err != nil {
		t.Fatalf("Failed to read updated person.age: %v", err)
	}
	if age != 35 {
		t.Errorf("Expected age 35, got %d", age)
	}

	// Update settings.theme to "light".
	updatedYAML, err = Update(updatedYAML, "settings.theme", "light")
	if err != nil {
		t.Fatalf("Failed to update settings.theme: %v", err)
	}

	// Verify updated theme.
	var theme string
	if err := Read(updatedYAML, "settings.theme", &theme); err != nil {
		t.Fatalf("Failed to read updated settings.theme: %v", err)
	}
	if theme != "light" {
		t.Errorf("Expected theme 'light', got %q", theme)
	}

	// Ensure person.name remains unchanged.
	var name string
	if err := Read(updatedYAML, "person.name", &name); err != nil {
		t.Fatalf("Failed to read person.name: %v", err)
	}
	if name != "John Doe" {
		t.Errorf("Expected person.name to be 'John Doe', got %q", name)
	}
}

func TestNonExistentKey(t *testing.T) {
	// Attempt to read a key that doesn't exist.
	errExpected := "key nonexistent not found"
	var out string
	if err := Read([]byte(sampleYAML), "person.nonexistent", &out); err == nil {
		t.Fatalf("Expected error when reading non-existent key, got nil")
	} else if err.Error() != errExpected {
		t.Errorf("Expected error %q, got %q", errExpected, err.Error())
	}

	// Attempt to update a key that doesn't exist.
	if _, err := Update([]byte(sampleYAML), "person.nonexistent", 40); err == nil {
		t.Fatalf("Expected error when updating non-existent key, got nil")
	}
}

func TestUpdateWholeNode(t *testing.T) {
	// Replace the entire person node with a new struct.
	newPerson := Person{
		Name: "Jane Smith",
		Age:  25,
	}
	updatedYAML, err := Update([]byte(sampleYAML), "person", newPerson)
	if err != nil {
		t.Fatalf("Failed to update person node: %v", err)
	}

	// Read back the updated person node.
	var p Person
	if err := Read(updatedYAML, "person", &p); err != nil {
		t.Fatalf("Failed to read updated person node: %v", err)
	}
	if !reflect.DeepEqual(p, newPerson) {
		t.Errorf("Expected person %+v, got %+v", newPerson, p)
	}
}
