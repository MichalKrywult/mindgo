# MindGO

![Go Version](https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green.svg)

MindGO is a lightweight command-line mood tracker written in Go.

Track your daily mood, keep a history of your entries, discover patterns, and gain insights into your well-being — all from your terminal.

## Features

- 📝 **Mood tracking** – Record your mood on a scale of 1–10, with notes and tags.
- 📅 **Mood history** – Browse and sort your previous entries.
- 📊 **Statistics** – Explore averages, mood distribution, and weekday patterns.
- ✏️ **Entry management** – Add, edit, and delete mood entries.
- 💾 **JSON storage** – Keep your data locally.
- 📤 **CSV export** – Export your mood history for further analysis.
- ⚙️ **Configurable storage** – Choose where your data is stored.

## Getting Started

### Requirements

- Go 1.27.1 or newer.

### Installation

Clone the repository:

```bash
git clone https://github.com/MichalKrywult/mindgo.git
cd mindgo
```

### Run

Start the application directly:

```bash
go run .
```

Or build an executable:

```bash
go build -o mindgo .
./mindgo
```

## Configuration

By default, MindGO stores its configuration and mood entries in:

```text
<user-config-dir>/mindgo/
```

You can customize the data file location using command-line flags:

```bash
./mindgo -file work-moods.json
./mindgo -path /path/to/moods.json
```

## Project Structure

```text
mindgo/
├── internal/
│   ├── app/
│   ├── cli/
│   ├── config/
│   ├── domain/
│   ├── export/
│   ├── stats/
│   ├── storage/
│   └── tracker/
├── main.go
└── go.mod
```

MindGO follows a modular architecture, separating domain logic, data persistence, mood tracking, statistics, and user interaction.

## Development

Run tests:

```bash
go test ./...
```

## Contributing

Contributions are welcome! If you'd like to report a bug,
suggest a feature, or contribute code, check out
[CONTRIBUTING.md](CONTRIBUTING.md).

## License

See LICENSE for details.
