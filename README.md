# Pingen
ahoy!

Simple static site generator for markdown-based bookmarks with wiki-style linking.

## Features
another hello!

- Convert markdown files to static HTML
- Wiki-style internal links: `[[Page Name]]` → HTML links (compatible with marksman LSP)
- One markdown file = one HTML page
- Auto-generated navigation
- Clean, minimal styling
- CLI-focused workflow

## Installation

```bash
go build -o pingen
```

## Usage

```bash
# Generate site with defaults (source: ./, output: ./public)
./pingen

# Specify source and output directories
./pingen -s ./content -o ./public

# Use custom template
./pingen -s ./content -o ./public -t ./my-template.html
```

### Flags

- `-s`, `-source` - Source directory containing markdown files (default: `./`)
- `-o`, `-out` - Output directory for generated HTML (default: `./public`)
- `-t`, `-template` - Custom template file (optional, uses embedded template by default)

## Markdown Format

Each `.md` file becomes one HTML page:

```markdown
# Page Title

## Section Name
- [Bookmark](https://example.com) - description
  - Additional notes
  - Tags: go, cli

## References
See also: [[Other Page]] for related links
```

### Wiki Links

Use `[[Page Name]]` for internal links. The generator will:
1. Convert to proper slug: `page-name.html`
2. Create working HTML links
3. Work seamlessly with marksman LSP in neovim

## Example Workflow

```bash
# 1. Create your markdown files
mkdir content
echo "# Home\n\nSee [[Tools]] for resources" > content/home.md
echo "# Tools\n\n- [Go](https://golang.org)" > content/tools.md

# 2. Generate site
./pingen -s ./content -o ./public

# 3. Serve locally
cd public && python3 -m http.server 8000

# 4. Open browser to http://localhost:8000
```

## Project Structure

```
pingen/
├── main.go              # CLI and orchestration
├── parser.go            # Markdown parsing, wiki link conversion
├── generator.go         # Site generation logic (template embedded at build)
└── templates/
    └── page.html        # HTML template (embedded in binary)
```

## Portability

The binary is fully self-contained:
- Template is embedded at build time using `//go:embed`
- No external files required to run
- Copy `pingen` anywhere and it works
- Optional `-t` flag to override with custom template

## Dependencies

- [goldmark](https://github.com/yuin/goldmark) - Markdown parser (CommonMark compliant)
- Go standard library (including `embed`)

## Tips

- Use H2 (`##`) for sections
- Use H3 (`###`) or bullets for bookmarks
- Files and directories starting with `.` are ignored
- Directory structure is preserved in output
- Works great with neovim + marksman LSP for editing
