# VK-Quotes

VK-Quotes is a small console program for storing and managing quotes. Each quote is stored with its text, author, language, and date.

## Features

- Add, update, and delete quotes
- Show all quotes
- Find quotes by text
- Display statistics by author and language
- Read quotes in random order
- Find similar quotes
- Backup and restore your quote database
- Plays background music in the CLI

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/vk-quotes.git
   cd vk-quotes
   ```
2. Build the project:
   ```sh
   go build
   ```
3. Run the program:
   ```sh
   ./vk-quotes
   ```

## Usage

When you run the program, you’ll see a command-line interface. Use the following commands:

- `add` or `a` — Add a new quote
- `update <id>` or `u <id>` — Update a quote by ID
- `delete <id>` or `d <id>` — Delete a quote by ID
- `find` or `f` — Find a quote by text
- `showall` or `s` — Show all quotes
- `stats` — Show statistics
- `read` or `r` — Read quotes in random order
- `similarquotes` or `sim` — Find similar quotes
- `resetids` — Reset all quote IDs
- `quit` or `q` — Backup and exit

Example:
```sh
=> add
=> update 3
=> delete 2
=> stats
=> quit
```

## Project Structure

- `main.go` — Entry point
- `pkg/` — Main package code
  - `audio/` — Audio playback
  - `cmd/` — Command-line interface
  - `config/` — Configuration and constants
  - `db/` — Database logic (CRUD, backup, statistics, etc.)
  - `util/` — Utility functions
- `QUOTES/quotes.json` — Local database file

## Dependencies

- [github.com/faiface/beep](https://github.com/faiface/beep) — Audio playback
- [github.com/peterh/liner](https://github.com/peterh/liner) — Command-line input
- [github.com/jdkato/prose](https://github.com/jdkato/prose) — Text processing

## License

MIT License (or your chosen license)