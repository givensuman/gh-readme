<div align="center">
  <img src="./assets/logo.png" alt="gh-readme" width="200" />
</div>

# gh-readme

A GitHub CLI extension to render a README in the terminal.

## Install

```bash
gh extension install givensuman/gh-readme
```

## Usage

```
gh-readme <owner/repo> [--ref <branch|tag|sha>]
```

### Examples

```bash
gh readme cli/cli
gh readme charmbracelet/glamour --ref v0.6.0
```

## Flags

| Flag | Description |
|------|-------------|
| `--ref` | Branch, tag, or commit SHA to fetch README from |

## License

[MIT](./LICENSE)
