# envsync

> Utility to diff and sync `.env` files across environments with secret masking support.

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git
cd envsync && go build -o envsync .
```

---

## Usage

**Diff two `.env` files:**

```bash
envsync diff .env.local .env.production
```

**Sync missing keys from one file to another:**

```bash
envsync sync .env.local .env.production
```

**Mask secrets when outputting diffs:**

```bash
envsync diff .env.local .env.production --mask-secrets
```

Example output:

```
+ DB_HOST=localhost        (missing in production)
~ API_URL                  (value differs)
  SECRET_KEY=***masked***
```

### Flags

| Flag             | Description                          |
|------------------|--------------------------------------|
| `--mask-secrets` | Redact secret values in output       |
| `--dry-run`      | Preview sync changes without writing |
| `--output`       | Output format: `text`, `json`        |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername