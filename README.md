<p align="center">
 <img src="https://raw.githubusercontent.com/pashkov256/media/refs/heads/main/cmdpool/cmdpool.svg"/>
</p>
<p align="center">
          <a><img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT"></a>
        <a><img src="https://img.shields.io/badge/Go-1.21+-blue.svg" alt="Go"></a>
         <a href="https://goreportcard.com/report/github.com/pashkov256/cmdpool"> <img src="https://goreportcard.com/badge/github.com/pashkov256/cmdpool"/></a>

<p align="center">
    <em>Run multiple commands simultaneously with real-time output monitoring in a beautiful TUI interface.</em>
</p>

<hr>
</p>

**cmdpool** is a powerful CLI/TUI utility written in Go that allows you to run multiple commands simultaneously while displaying their real-time output in separate terminal panels. Perfect for developers, DevOps engineers, testers, and anyone who needs to monitor parallel tasks in one place.

## ✨ Features

- 🚀 **Parallel Execution**: Run multiple commands simultaneously with real-time monitoring
- 🖥️ **Beautiful TUI Interface**: Modern text-based user interface with automatic screen partitioning
- 🎨 **Color-Coded Output**: Green for stdout, red for stderr with clear visual distinction
- 📊 **Real-Time Status**: Live status indicators: 🟢 Running, ✅ Done, 🔴 Failed
- 🔍 **Log Mode**: Expand any command to full screen for detailed log viewing
- 📜 **Scrollable History**: Navigate through command output history with search functionality
- ⚡ **Live Management**: Stop/restart commands on the fly (r, s keys) without exiting
- ➕ **Dynamic Commands**: Add new commands without leaving the program
- ⚙️ **Configuration Support**: Save and load command sets via `.cmdpool.yml` or TUI interface
- 📈 **Resource Monitoring**: Mini CPU/RAM graphs for each running command
- ⏱️ **Execution Timer**: Track how long each command has been running
- 🎯 **Search & Filter**: Search through logs using `/` like in less

## 🚀 Quick Start

### Basic Usage

Run multiple commands simultaneously:

```bash
cmdpool "ping google.com" "ping github.com" "make build"
```

Each command will be displayed in its own panel with real-time output.

### Interactive TUI Mode

```bash
cmdpool
```

Launch the interactive interface to manage commands visually.

## 📦 Installation

### Using Go

```bash
go install github.com/pashkov256/cmdpool@latest
```

### From Source

```bash
git clone https://github.com/pashkov256/cmdpool.git
cd cmdpool
go build -o cmdpool .
```

## 🛠 Usage

### Command Line Mode

```bash
# Run multiple commands
cmdpool "npm run dev" "go run main.go" "docker compose up"

# Run with configuration file
cmdpool -config .cmdpool.yml

# Run specific command set
cmdpool -set backend
```

### Interactive TUI Mode

```bash
cmdpool
```

Navigate with:

- **Arrow Keys**: Move between panels
- **Enter**: Expand panel to full screen
- **r**: Restart command
- **s**: Stop command
- **+**: Add new command
- **/**: Search in logs
- **q**: Quit

## ⚙️ Configuration

Create a `.cmdpool.yml` file in your project:

```yaml
commands:
  - name: Backend
    cmd: go run main.go
    dir: ./backend

  - name: Frontend
    cmd: npm run dev
    dir: ./frontend

  - name: Database
    cmd: docker compose up db
    dir: ./

  - name: Tests
    cmd: go test ./...
    dir: ./
    auto_restart: true
```

### Configuration Options

| Option         | Description                  | Default           |
| -------------- | ---------------------------- | ----------------- |
| `name`         | Display name for the command | Required          |
| `cmd`          | Command to execute           | Required          |
| `dir`          | Working directory            | Current directory |
| `auto_restart` | Restart on failure           | false             |
| `env`          | Environment variables        | {}                |

## 🎯 Use Cases

### Development Workflow

```bash
cmdpool "go run main.go" "npm run dev" "docker compose up"
```

### DevOps Monitoring

```bash
cmdpool "kubectl logs -f deployment/app" "docker stats" "htop"
```

### Testing & CI

```bash
cmdpool "go test ./..." "npm test" "python -m pytest"
```

### System Administration

```bash
cmdpool "tail -f /var/log/nginx/access.log" "iostat 1" "netstat -i 1"
```

## 🔧 Advanced Features

### Log Mode

Press **Enter** on any panel to expand it to full screen:

- Full-screen command output
- Scrollable history
- Search functionality (`/` key)
- Export logs to file

### Command Management

- **Restart (r)**: Restart a stopped or failed command
- **Stop (s)**: Stop a running command
- **Add (+)**: Add new commands dynamically
- **Remove**: Remove completed commands

### Resource Monitoring

Each panel shows:

- CPU usage mini-graph
- Memory consumption
- Execution time
- Exit status

## 🛠 Contributing

We welcome and appreciate any contributions to cmdpool!
There are many ways you can help us grow and improve:

- **🐛 Report Bugs** — Found an issue? Let us know by opening an issue.
- **💡 Suggest Features** — Got an idea for a new feature? We'd love to hear it!
- **📚 Improve Documentation** — Help us make the docs even clearer and easier to use.
- **💻 Submit Code** — Fix a bug, refactor code, or add new functionality by submitting a pull request.

Before contributing, please take a moment to read our [CONTRIBUTING.md](https://github.com/pashkov256/cmdpool/blob/main/CONTRIBUTING.md) guide.
It explains how to set up the project, coding standards, and the process for submitting contributions.

Together, we can make cmdpool even better! 🚀

## 📋 Requirements

- **Go**: 1.21 or higher
- **Terminal**: ANSI-compatible terminal with support for colors and TUI
- **OS**: Linux, macOS, Windows (with WSL or Git Bash)

## 📜 License

This project is distributed under the **MIT** license.

---

### Thank you to these wonderful people for their contributions!

<a href="https://github.com/pashkov256/cmdpool/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=pashkov256/cmdpool" />
</a>
