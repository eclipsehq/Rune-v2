# Rune V2

Rune V2 is a high-performance, modern Discord self-bot written in Go. It serves as a complete remaster of the original project by **eclipsehq**, rebuilt from the ground up for stability, speed, and maintainability.

Developed with ❤️ by **Light (eclipsehq)** and **Heavenzone**.

## 🚀 Features

- **Cleaner Codebase**: Rewritten using Go for better concurrency and memory management.
- **Actually Runnable**: Focus on stability and ease of deployment.
- **Maintained**: Active development and bug fixes.
- **ANSI Styling**: Beautifully formatted terminal-style responses in Discord using ANSI escape codes.
- **Dynamic Prefix**: Change your bot prefix on the fly.

## 🗺️ Roadmap

- [ ] **Web Panel**: Manage your self-bot through a dedicated web interface.
- [ ] **200+ Commands**: Integration of a massive library of utility, fun, and information commands.
- [ ] **Advanced Logging**: Track interactions and events with more detail.

## 🛠️ Installation

1. Ensure you have [Go](https://go.dev/doc/install) (1.21+) installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/eclipsehq/Rune-v2.git
   cd rune-v2
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build the project:
   ```bash
   go build -o rune
   ```

## ⚙️ Configuration

Edit the `cfg/config.json` file with your credentials:

```json
{
  "token": "YOUR_USER_TOKEN",
  "owner_id": "YOUR_DISCORD_ID",
  "prefix": "&"
}
```

> **Warning**: Never share your token. Rune V2 is a self-bot; using self-bots can lead to account termination by Discord. Use at your own risk.

## 📜 License

This project is licensed under the terms of the license included in the repository.

---