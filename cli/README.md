# Mailify CLI

Mailify CLI is a powerful command-line tool for email validation and mail server information retrieval. Built on top of the [mailify](https://github.com/adarsh-jaiss/mailify) library, it provides an easy-to-use interface for validating email addresses and getting mail server details.

## Features

- Single email address validation
- Bulk email validation using Excel files
- Mail server lookup for domains
- Mail server lookup for email addresses
- Simple flag-based interface

## Installation

### Using Go

```bash
go install github.com/yourusername/mailify-cli@latest
```

### Using Binary Releases

Download the latest binary for your platform from the [releases page](https://github.com/yourusername/mailify-cli/releases).

## Usage

The CLI tool uses a flag-based interface where all operations are performed using the main `mailify` command with different flags.

### Required Flag

- `-s, --sender`: Sender email address (required for all operations)

### Operation Flags

You can use one of the following operation flags per command:

- `-v, --validate`: Validate a single email address
- `-e, --excel`: Process and validate emails from an Excel file
- `-d, --domain`: Get mail servers for a domain
- `-r, --receipient`: Get mail servers for a recipient email

### Examples

1. **Validate a single email address**
```bash
mailify -s your@email.com -v user@example.com
```

2. **Bulk validate emails from Excel file**
```bash
mailify -s your@email.com -e emails.xlsx
```

3. **Get mail servers for a domain**
```bash
mailify -s your@email.com -d example.com
```

4. **Get mail servers for an email address**
```bash
mailify -s your@email.com -r user@example.com
```

### Help

```bash
mailify --help
```

## Development Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/mailify-cli.git
cd mailify-cli
```

2. Install dependencies
```bash
go mod download
```

3. Build the project
```bash
go build -o mailify
```

## Release Process

### Using GoReleaser

1. Install GoReleaser:
```bash
brew install goreleaser
```

2. Create `.goreleaser.yml`:
```yaml
before:
  hooks:
    - go mod tidy
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - amd64
      - arm64
    binary: mailify
archives:
  - replacements:
      darwin: Darwin
      linux: Linux
      windows: Windows
      amd64: x86_64
checksum:
  name_template: 'checksums.txt'
snapshot:
  name_template: "{{ incpatch .Version }}-next"
changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
```

3. Release Steps:
```bash
# Create and push a new tag
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0

# Set GitHub token
export GITHUB_TOKEN="your-github-token"

# Create release
goreleaser release --rm-dist
```

## Excel File Format

When using the `-e, --excel` flag, your Excel file should:
- Have a column containing email addresses
- Be in `.xlsx` format
- The tool will create a new column with validation results

## Error Handling

- If no operation flag is specified, the tool will show an error message
- All operations will return meaningful error messages if something goes wrong
- The tool validates inputs before processing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please:
1. Check the [GitHub Issues](https://github.com/yourusername/mailify-cli/issues)
2. Create a new issue if your problem isn't already reported