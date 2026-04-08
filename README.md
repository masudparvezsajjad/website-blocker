# AdultBlocker (macOS MVP)

A lightweight, terminal-based adult content blocker built with Go.
Designed for macOS, this tool blocks unwanted websites system-wide using `/etc/hosts` and serves a local block page.

---

## 🚀 Features

- 🔒 System-wide domain blocking (works across all browsers)
- 🧠 Simple rule-based domain management
- 🌐 Local block page (`127.0.0.1`)
- ⚡ Fast and lightweight (pure Go, no heavy dependencies)
- 🛠 CLI-based control (install, enable, disable, status)
- 🧩 Easy to extend (future DNS / Network Extension support)

---

## 🏗 How It Works

1. Blocked domains are mapped to `127.0.0.1` via `/etc/hosts`
2. A local HTTP server runs on `127.0.0.1:8088`
3. When a blocked site is visited:

- DNS resolves to localhost
- Browser hits local server
- A custom "blocked" page is shown

> ⚠️ Note: Some HTTPS sites may show a browser warning instead of the custom page. This is expected in the MVP version.

---

## 📦 Installation

### 1. Clone & Build

```bash
git clone https://github.com/yourname/adult-blocker.git
cd adult-blocker
go build -o mps-blocker ./cmd/blocker
```

### 2. Install (requires sudo)

```bash
sudo ./mps-blocker install
```

---

## ▶️ Usage

### Enable Blocking

```bash
sudo ./blocker enable
```

### Start Block Page Server

```bash
sudo ./mps-blocker daemon
```

### Check Status

```bash
sudo ./mps-blocker status
```

### Disable Blocking

```bash
sudo ./mps-blocker disable
```

### Uninstall

```bash
sudo ./mps-blocker uninstall
```

---

## 🌍 Manage Domains

### Add Domain

```bash
sudo ./mps-blocker add-domain example.com
```

### Remove Domain

```bash
sudo ./mps-blocker remove-domain example.com
```

---

## ⚙️ Configuration

Config file location:

```bash
/Library/Application Support/AdultBlocker/config.json
```

Example:

```json
{
  "enabled": true,
  "block_page_port": 8088,
  "blocked_domains": ["yoursite.com", "youradultvideos.com"],
  "redirect_ipv4": "127.0.0.1",
  "redirect_ipv6": "::1"
}
```

---

## 🔄 Auto Start (launchd)

To run automatically on macOS startup:

1. Copy binary:

```bash
sudo cp ./mps-blocker /usr/local/bin/mps-blocker
```

1. Create plist:

```bash
sudo nano /Library/LaunchDaemons/com.adultblocker.daemon.plist
```

1. Load service:

```bash
sudo launchctl bootstrap system /Library/LaunchDaemons/com.adultblocker.daemon.plist
```

---

## 🔐 Security Considerations

This app is safe when used correctly:

### ✔ Safe Practices

- Runs only on `127.0.0.1` (no external exposure)
- No remote API or network access
- No shell injection (safe file handling)
- Requires `sudo` only for system changes

### ⚠ Things to Avoid

- Do not bind server to `0.0.0.0`
- Do not add remote control without authentication
- Do not execute untrusted scripts or configs

---

## ⚠ Limitations (MVP)

- ❌ Easy to bypass (hosts file can be edited)
- ❌ No protection against VPN / DNS-over-HTTPS
- ❌ HTTPS block page not always shown
- ❌ No category-based filtering yet

---

## 🛣 Roadmap

### Version 2

- Password-protected disable/uninstall
- Domain import (blocklists)
- Logging (SQLite)
- Auto launchd setup

### Version 3

- Local DNS server
- Stronger anti-bypass mechanisms

### Version 4

- macOS Network Extension (DNS Proxy)
- Advanced filtering + category detection

---

## 🧠 Future Vision

- Cross-platform (macOS, Windows, Linux)
- System-level enforcement (WFP / Network Extension)
- AI-powered content classification
- Parental control / accountability features

---

## 🤝 Contributing

Contributions are welcome!

- Fork the repo
- Create a feature branch
- Submit a PR

---

## 📄 License

MIT License (or your preferred license)

---

## 🙌 Disclaimer

This tool is intended for **personal productivity and self-control purposes**.
It is not a full security solution and should not be relied upon for enterprise-grade filtering.

---
