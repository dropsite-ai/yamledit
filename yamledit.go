package yamledit

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse parses YAML bytes into a document node and returns a pointer to the document.
func Parse(yamlBytes []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &doc); err != nil {
		return nil, err
	}
	if len(doc.Content) == 0 {
		return nil, fmt.Errorf("empty YAML document")
	}
	return &doc, nil
}

// Encode encodes the entire document node back into YAML bytes.
func Encode(doc *yaml.Node) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ReadNode decodes the YAML value at the dot-notation path from the document node into out.
// The document node should be obtained via Parse.
func ReadNode(doc *yaml.Node, dotPath string, out interface{}) error {
	if len(doc.Content) == 0 {
		return fmt.Errorf("empty YAML document")
	}
	keys := strings.Split(dotPath, ".")
	// Use the document’s first content node (usually the mapping root).
	mappingNode := doc.Content[0]
	node, err := findNode(mappingNode, keys)
	if err != nil {
		return err
	}
	return node.Decode(out)
}

// UpdateNode updates the YAML value at the dot-notation path in the document node with newValue.
// It modifies the document in place.
func UpdateNode(doc *yaml.Node, dotPath string, newValue interface{}) error {
	if len(doc.Content) == 0 {
		return fmt.Errorf("empty YAML document")
	}
	keys := strings.Split(dotPath, ".")
	mappingNode := doc.Content[0]
	return updateNode(mappingNode, keys, newValue)
}

// Read is a wrapper that accepts YAML bytes, parses them, and reads the value at dotPath.
func Read(yamlBytes []byte, dotPath string, out interface{}) error {
	doc, err := Parse(yamlBytes)
	if err != nil {
		return err
	}
	return ReadNode(doc, dotPath, out)
}

// Update is a wrapper that accepts YAML bytes, parses them, updates the value at dotPath,
// and returns the updated YAML as bytes.
func Update(yamlBytes []byte, dotPath string, newValue interface{}) ([]byte, error) {
	doc, err := Parse(yamlBytes)
	if err != nil {
		return nil, err
	}
	if err := UpdateNode(doc, dotPath, newValue); err != nil {
		return nil, err
	}
	return Encode(doc)
}

// findNode recursively traverses the node tree using dot-notation keys and returns the target node.
func findNode(node *yaml.Node, keys []string) (*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected a mapping node, got kind %v", node.Kind)
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == keys[0] {
			if len(keys) == 1 {
				return node.Content[i+1], nil
			}
			return findNode(node.Content[i+1], keys[1:])
		}
	}
	return nil, fmt.Errorf("key %s not found", keys[0])
}

// updateNode navigates to the target node (using dot-notation keys) and replaces it with newValue.
func updateNode(node *yaml.Node, keys []string, newValue interface{}) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got kind %v", node.Kind)
	}
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == keys[0] {
			if len(keys) == 1 {
				newNode, err := createYAMLNode(newValue)
				if err != nil {
					return err
				}
				node.Content[i+1] = newNode
				return nil
			}
			return updateNode(node.Content[i+1], keys[1:], newValue)
		}
	}
	return fmt.Errorf("key %s not found", keys[0])
}

// createYAMLNode marshals newValue and then unmarshals it into a yaml.Node.
// This lets us update with any kind of value.
func createYAMLNode(newValue interface{}) (*yaml.Node, error) {
	data, err := yaml.Marshal(newValue)
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}
	// Unmarshal returns a document node; return its first content.
	if len(node.Content) > 0 {
		return node.Content[0], nil
	}
	return nil, fmt.Errorf("could not create YAML node from value")
}
