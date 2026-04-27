# portwatch

A lightweight CLI daemon that monitors port usage and alerts on unexpected bindings or conflicts.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon and watch for unexpected port bindings:

```bash
portwatch start --ports 8080,3306,5432 --interval 5s
```

Watch a range of ports and log alerts to a file:

```bash
portwatch start --range 3000-9000 --log /var/log/portwatch.log
```

Run a one-time scan and print current bindings:

```bash
portwatch scan
```

Stop a running daemon:

```bash
portwatch stop
```

### Flags

| Flag | Description | Default |
|------|-------------|----------|
| `--ports` | Comma-separated list of ports to watch | — |
| `--range` | Port range to monitor (e.g. `3000-9000`) | — |
| `--interval` | Poll interval | `10s` |
| `--log` | Log output file | stdout |
| `--alert` | Alert command to run on conflict | — |
| `--pid` | Path to PID file (used by `stop`) | `/tmp/portwatch.pid` |

---

## How It Works

`portwatch` polls the system's active network connections at a configurable interval. When a new process binds to a watched port — or two processes conflict — it emits an alert to stdout, a log file, or a custom command of your choice.

---

## License

MIT © 2024 yourusername
