# Slidedeck

This directory contains Marp slides for the BlobTube project.

## Prerequisites

Make sure you're in the nix development environment:
```bash
nix develop
```

## Previewing Slides

### Option 1: VS Code Extension (Recommended)
1. Install the [Marp for VS Code](https://marketplace.visualstudio.com/items?itemName=marp-team.marp-vscode) extension
2. Open `slides.md` in VS Code
3. Click the preview button or use `Ctrl+Shift+V`

### Option 2: CLI Watch Mode
```bash
cd slidedeck
marp slides.md --watch
```
This will open a browser window that auto-refreshes on changes.

## Building Slides

### Export to HTML
```bash
cd slidedeck
marp slides.md -o output.html
```

### Export to PDF
```bash
cd slidedeck
marp slides.md -o output.pdf
```

### Export to PowerPoint
```bash
cd slidedeck
marp slides.md -o output.pptx
```

## Marp Syntax Reference

### Creating Slides
Use `---` to separate slides:
```markdown
---

# New Slide

Content here
```

### Slide Directives
Add directives in HTML comments:
```markdown
<!-- _class: lead -->
# Centered Slide

<!-- _backgroundColor: #123456 -->
# Custom Background

<!-- _paginate: false -->
# No Page Number
```

### Frontmatter Options
Configure at the top of slides.md:
```yaml
---
marp: true
theme: gaia
paginate: true
header: 'BlobTube'
footer: '2026'
---
```

## Resources
- [Marp Official Documentation](https://marpit.marp.app/)
- [Marp CLI Documentation](https://github.com/marp-team/marp-cli)
- [Gaia Theme](https://github.com/marp-team/marp-core/tree/main/themes)
