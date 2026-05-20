
# Rune V2

**Rune V2** is a high-performance, modern Discord self-bot written entirely in **Go**. It is a complete ground-up remaster of the original project by **eclipsehq**, rebuilt for maximum stability, speed, and maintainability.

Developed with ❤️ by **Light (eclipsehq)** and **val**.

## 🚀 Features

- **Clean & Modern Codebase**: Fully rewritten in Go for excellent concurrency, low memory usage, and high performance.
- **Lightweight**: Only **2.6 MB** of memory usage with **46 commands** loaded on modern hardware.
- **Production-Ready**: Designed to be stable and easy to deploy.
- **Actively Maintained**: Regular updates and bug fixes.
- **Beautiful ANSI Styling**: Clean, terminal-style formatted responses in Discord using ANSI escape codes.
- **Dynamic Prefix**: Change your command prefix on the fly.
- **40-70+ Commands**: A growing suite of utility, fun, and information commands (exact count fluctuates as we frequently add new ones).

## 🗺️ Goals & Roadmap

Our goals remain the same: deliver the best possible self-bot experience — fast, reliable, and user-friendly.

**Current Roadmap:**
- [ ] **Web Panel**: Full web interface to manage your self-bot remotely.
- [ ] **100+ Commands**: Expanding the command library significantly.
- [ ] **Advanced Logging & Analytics**: Detailed tracking of events and interactions.
- [ ] **Frontend Integration**: Small frontend components coming soon.

## ⚠️ Important Notice

**This project is still in active development.**  
We **do not** want this project to be skidded. Please respect the work of the developers and do not re-upload or claim this as your own.

> **Warning**: Rune V2 is a **self-bot**. Using self-bots violates Discord's Terms of Service and can result in your account being terminated. Use at your own risk. Never share your token.

## 🛠️ Installation

1. Install [Go](https://go.dev/doc/install) (1.21 or higher).
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

Edit `cfg/config.json` with your details:

```json
{
  "token": "YOUR_USER_TOKEN",
  "owner_id": "YOUR_DISCORD_ID",
  "prefix": "&"
}
```

## 📜 License

This project is licensed under the terms included in the repository.

---

**Made with ❤️ by Light & val**  
*Still cooking — stay tuned for more.*
```
