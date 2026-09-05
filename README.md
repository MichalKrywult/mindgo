# MindGO

A simple mood tracker written in Go.

The project is created primarily for learning Go through building a real-life application.

## Features

- Mood entries
- Mood tracking
- Basic statistics
- Persistent storage
- CLI interface
- JSON data storage
- Configurable application data path

The application stores its data in the user's standard configuration directory:

```md
<user-config-dir>/mindgo/moods.json
```

The **mindgo** directory is created automatically when the application starts.

## Project structure

```internal/
├── cli/ # Command-line interface
├── config/ # Application configuration
├── domain/ # Core domain models
├── stats/ # Statistics and calculations
├── storage/ # Data persistence
└── tracker/ # Mood tracking logic
```

## Planned features

- More mood-related statistics
- Improved CLI experience
- Additional commands and options
- Better error handling
- More comprehensive test coverage

## Development

The project is developed incrementally, with each version introducing new Go concepts and features, showing my learning path.
