package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/masudparvezsajjad/website-blocker/internal/blockpage"
	"github.com/masudparvezsajjad/website-blocker/internal/config"
	"github.com/masudparvezsajjad/website-blocker/internal/hosts"
	platform "github.com/masudparvezsajjad/website-blocker/internal/platform/darwin"
	"github.com/masudparvezsajjad/website-blocker/internal/util"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	if err := platform.EnsureMacOS(); err != nil {
		return err
	}

	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	cmd := os.Args[1]

	switch cmd {
	case "install":
		return install()
	case "enable":
		return enable()
	case "disable":
		return disable()
	case "status":
		return status()
	case "uninstall":
		return uninstall()
	case "add-domain":
		if len(os.Args) < 3 {
			return errors.New("missing domain")
		}
		return addDomain(os.Args[2])
	case "remove-domain":
		if len(os.Args) < 3 {
			return errors.New("missing domain")
		}
		return removeDomain(os.Args[2])
	case "daemon":
		return daemon()
	default:
		printUsage()
		return nil
	}
}

func printUsage() {
	fmt.Println(`AdultBlocker macOS MVP

Commands:
  sudo blocker install
  sudo blocker enable
  sudo blocker disable
  sudo blocker status
  sudo blocker uninstall
  sudo blocker add-domain <domain>
  sudo blocker remove-domain <domain>
  sudo blocker daemon`)
}

func install() error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := hosts.BackupHosts(); err != nil {
		return err
	}

	if cfg.Enabled {
		if err := hosts.Apply(cfg); err != nil {
			return err
		}
	}

	fmt.Println("Installed config at:", config.ConfigPath())
	fmt.Println("Hosts backup created at: /etc/hosts.adultblocker.bak")
	fmt.Println("Install complete.")
	fmt.Println("Run 'sudo blocker enable' to activate blocking.")
	return nil
}

func enable() error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cfg.Enabled = true
	if err := config.Save(cfg); err != nil {
		return err
	}

	if err := hosts.Apply(cfg); err != nil {
		return err
	}

	fmt.Println("Blocking enabled.")
	fmt.Println("Start local server in another terminal with: sudo blocker daemon")
	return nil
}

func disable() error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cfg.Enabled = false
	if err := config.Save(cfg); err != nil {
		return err
	}

	if err := hosts.Remove(); err != nil {
		return err
	}

	fmt.Println("Blocking disabled.")
	return nil
}

func uninstall() error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	if err := hosts.Remove(); err != nil {
		return err
	}

	_ = os.Remove(config.ConfigPath())

	fmt.Println("Managed hosts entries removed.")
	fmt.Println("Config removed.")
	fmt.Println("Uninstall complete.")
	return nil
}

func status() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	present, err := hosts.ManagedEntriesPresent()
	if err != nil {
		return err
	}

	fmt.Println("Enabled:", cfg.Enabled)
	fmt.Println("Config path:", config.ConfigPath())
	fmt.Println("Block page port:", cfg.BlockPagePort)
	fmt.Println("Managed hosts entries present:", present)
	fmt.Println("Blocked domains:", len(cfg.BlockedDomains))
	for _, d := range cfg.BlockedDomains {
		fmt.Println(" -", d)
	}

	return nil
}

func addDomain(raw string) error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	domain := util.NormalizeDomain(raw)
	if !util.IsValidDomain(domain) {
		return fmt.Errorf("invalid domain: %s", raw)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	exists := false
	for _, d := range cfg.BlockedDomains {
		if strings.EqualFold(d, domain) {
			exists = true
			break
		}
	}

	if !exists {
		cfg.BlockedDomains = append(cfg.BlockedDomains, domain)
		sort.Strings(cfg.BlockedDomains)
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	if cfg.Enabled {
		if err := hosts.Apply(cfg); err != nil {
			return err
		}
	}

	fmt.Println("Added domain:", domain)
	return nil
}

func removeDomain(raw string) error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	domain := util.NormalizeDomain(raw)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var updated []string
	for _, d := range cfg.BlockedDomains {
		if !strings.EqualFold(d, domain) {
			updated = append(updated, d)
		}
	}
	cfg.BlockedDomains = updated

	if err := config.Save(cfg); err != nil {
		return err
	}

	if cfg.Enabled {
		if err := hosts.Apply(cfg); err != nil {
			return err
		}
	}

	fmt.Println("Removed domain:", domain)
	return nil
}

func daemon() error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.Enabled {
		if err := hosts.Apply(cfg); err != nil {
			return err
		}
	}

	fmt.Printf("Starting block page server on 127.0.0.1:%d\n", cfg.BlockPagePort)
	server := &blockpage.Server{Port: cfg.BlockPagePort}
	return server.Start()
}
