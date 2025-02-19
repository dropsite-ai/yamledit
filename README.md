# yamledit

`yamledit` is a lightweight Go package and CLI tool for reading and updating YAML files using dot-notation paths (e.g., `person.name`) while preserving formatting, comments, and structure.

## Features

- **Read & Update YAML**: Extract and modify values at any depth.
- **Preserves Formatting**: Retains comments, ordering, and structure.
- **CLI Tool**: Quickly edit YAML files from the command line.
- **Go Package**: Programmatically manage YAML data in your projects.
- **Advanced Node API**: Work directly with parsed YAML nodes to avoid repeated parsing/encoding.

## Installation

### Go Package
```bash
go get github.com/dropsite-ai/yamledit
```

### Homebrew (macOS & Linux)
```bash
brew tap dropsite-ai/homebrew-tap
brew install yamledit
```

### Download Binary
Get the latest release [here](https://github.com/dropsite-ai/yamledit/releases).

### Build from Source
```bash
git clone https://github.com/dropsite-ai/yamledit.git
cd yamledit
go build -o yamledit cmd/main.go
```

## Usage

### As a Go Package

#### Using the Convenience Wrapper Functions

The wrapper functions accept YAML data as `[]byte`:

```go
import "github.com/dropsite-ai/yamledit"
import "gopkg.in/yaml.v3"

var name string
yamledit.Read(yamlData, "person.name", &name)

updatedYAML, err := yamledit.Update(yamlData, "person.age", 35)
if err != nil {
    // handle error
}
```

#### Using the Advanced Node API

For multiple operations on the same document without repeated unmarshaling/encoding, use the Node API:

```go
import (
    "github.com/dropsite-ai/yamledit"
    "gopkg.in/yaml.v3"
)

// Parse the YAML into a document node.
doc, err := yamledit.Parse(yamlData)
if err != nil {
    // handle error
}

// Read a value.
var name string
if err := yamledit.ReadNode(doc, "person.name", &name); err != nil {
    // handle error
}

// Update a value.
if err := yamledit.UpdateNode(doc, "person.age", 40); err != nil {
    // handle error
}

// Encode the modified document back to YAML bytes.
updatedYAML, err := yamledit.Encode(doc)
if err != nil {
    // handle error
}
```

### CLI Commands

**Read a value:**
```bash
yamledit -op=read -file=config.yaml -path=person.name
```

**Update a value:**
```bash
yamledit -op=update -file=config.yaml -path=person.age -value=35 -out=updated.yaml
```

**Copy a value:**
```bash
value=$(yamledit get --file source.yaml --path source.path)
yamledit set --file target.yaml --path target.path --value "$value"
```

## Contributing

Contributions welcome! Open an issue or submit a pull request.

## License

[MIT License](LICENSE)