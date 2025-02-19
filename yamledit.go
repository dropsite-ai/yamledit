package yamledit

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Read decodes the YAML value at the dot-notation path into out.
// out can be any type (string, struct, map, etc.) that matches the YAML content.
func Read(yamlBytes []byte, dotPath string, out interface{}) error {
	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return err
	}
	keys := strings.Split(dotPath, ".")
	if len(root.Content) == 0 {
		return fmt.Errorf("empty YAML document")
	}
	node, err := findNode(root.Content[0], keys)
	if err != nil {
		return err
	}
	return node.Decode(out)
}

// Update replaces the YAML value at the dot-notation path with newValue (which can be any type)
// and returns the updated YAML as bytes.
func Update(yamlBytes []byte, dotPath string, newValue interface{}) ([]byte, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return nil, err
	}
	keys := strings.Split(dotPath, ".")
	if len(root.Content) == 0 {
		return nil, fmt.Errorf("empty YAML document")
	}
	if err := updateNode(root.Content[0], keys, newValue); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&root); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
	// Unmarshal returns a document node; use its first content node.
	if len(node.Content) > 0 {
		return node.Content[0], nil
	}
	return nil, fmt.Errorf("could not create YAML node from value")
}
