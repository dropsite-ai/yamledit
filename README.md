# yamledit

`yamledit` is a lightweight Go package and CLI tool for reading and updating YAML files using dot-notation paths (e.g., `person.name`) while preserving formatting, comments, and structure.

## Features

- **Read & Update YAML**: Extract and modify values at any depth.
- **Preserves Formatting**: Retains comments, ordering, and structure.
- **CLI Tool**: Quickly edit YAML files from the command line.
- **Go Package**: Programmatically manage YAML data in your projects.

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
```go
import "github.com/dropsite-ai/yamledit"

var name string
yamledit.Read(yamlData, "person.name", &name)

updatedYAML, _ := yamledit.Update(yamlData, "person.age", 35)
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