# website-blocker

**AdultBlocker** is a small macOS CLI that blocks unwanted websites system-wide by updating `/etc/hosts` and serving a simple “blocked” page on localhost.

## Requirements

- macOS (Apple Silicon or Intel)
- `sudo` for changing `/etc/hosts` and installing the binary under `/usr/local` or `/opt/homebrew`

## Install a release (recommended)

GitHub Releases publish `.tar.gz` archives. The install script downloads the right archive for your Mac, installs the `mps-blocker` binary, and tells you what to run next.

**Option A — run the script (review it first):**

```bash
curl -fsSL https://raw.githubusercontent.com/masudparvezsajjad/website-blocker/main/scripts/install.sh | bash
```

**Option B — download, inspect, then run:**

```bash
curl -fsSL -O https://raw.githubusercontent.com/masudparvezsajjad/website-blocker/main/scripts/install.sh
chmod +x install.sh
./install.sh
```

### What `scripts/install.sh` does

| Topic | Behavior |
|--------|----------|
| **OS** | macOS only; other systems exit with an error. |
| **Architecture** | Tries `mps-blocker_darwin_{arm64,amd64}.tar.gz` first; if missing (older releases), falls back to `adult-blocker_darwin_*.tar.gz` and installs the `blocker` binary as `mps-blocker`. |
| **Version** | `VERSION` defaults to `latest` (GitHub “latest” release). Pin a tag: `VERSION=v1.2.3 ./install.sh`. |
| **Download** | `https://github.com/masudparvezsajjad/website-blocker/releases/.../download/.../<asset>.tar.gz` |
| **Install path** | `/opt/homebrew/bin` on Apple Silicon if that directory exists; otherwise `/usr/local/bin`. Uses `sudo install -m 0755`. |

After the script finishes:

```bash
sudo mps-blocker install
sudo mps-blocker enable
sudo mps-blocker daemon
```

Keep `daemon` running in a terminal session, or register it with `launchd` if you want it to start automatically.

## Build from source

Requires [Go](https://go.dev/dl/) (see `go.mod` for the toolchain version).

```bash
git clone https://github.com/masudparvezsajjad/website-blocker.git
cd website-blocker
go build -o mps-blocker ./cmd/blocker
sudo ./mps-blocker install
sudo ./mps-blocker enable
sudo ./mps-blocker daemon
```

## Commands

All commands that change the system need `sudo`.

| Command | Description |
|---------|-------------|
| `install` | Create config if needed, back up `/etc/hosts`, apply rules if blocking is enabled |
| `enable` / `disable` | Turn blocking on or off |
| `status` | Show whether blocking is active |
| `daemon` | HTTP server on localhost for the block page (default port **8088**) |
| `add-domain <host>` / `remove-domain <host>` | Edit the blocked list |
| `uninstall` | Remove blocker-managed host entries and related state |

## How it works

1. Blocked domains are pointed at loopback in `/etc/hosts` (IPv4/IPv6 from config).
2. With `daemon` running, HTTP requests to those names hit the local server and show the block page.
3. **HTTPS** may show a browser certificate warning instead of the custom page; that is a normal limitation of hosts-based blocking.

## Configuration

Config file: **`~/Library/Application Support/AdultBlocker/config.json`**

Relevant fields include `enabled`, `block_page_port`, `blocked_domains`, `redirect_ipv4`, and `redirect_ipv6`. A default file is created on first use.

## Limitations

- Anyone with admin access can edit `/etc/hosts` or disable the tool.
- VPNs, DNS-over-HTTPS, or custom resolvers can bypass hosts-based blocking.
- This is not enterprise-grade or child-proof filtering.

## Security notes

- The block page server should stay bound to localhost only.
- There is no remote control API; configuration is local JSON.

## Contributing

Contributions are welcome: open an issue or a pull request with a short description of the change and how you tested it.

## License

MIT

## Disclaimer

This project is intended for **personal productivity and self-control**. It is provided as-is, without warranty, and is not a complete security or parental-control solution.
