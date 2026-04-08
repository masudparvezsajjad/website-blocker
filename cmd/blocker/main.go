package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/masudparvezsajjad/website-blocker/internal/blockpage"
	"github.com/masudparvezsajjad/website-blocker/internal/config"
	"github.com/masudparvezsajjad/website-blocker/internal/hosts"
	platform "github.com/masudparvezsajjad/website-blocker/internal/platform/darwin"
	"github.com/masudparvezsajjad/website-blocker/internal/reflectpause"
	"github.com/masudparvezsajjad/website-blocker/internal/util"
)

var (
	cliTitle   = color.New(color.FgHiCyan, color.Bold)
	cliHeading = color.New(color.FgCyan, color.Bold)
	cliSuccess = color.New(color.FgHiGreen)
	cliDanger  = color.New(color.FgHiRed)
	cliInfo    = color.New(color.FgHiBlue)
	cliMuted   = color.New(color.Faint)
	cliErr     = color.New(color.FgHiRed, color.Bold)
	cliWarn    = color.New(color.FgYellow)
	cliKey     = color.New(color.FgCyan)
)

func main() {
	if err := run(); err != nil {
		cliErr.Fprintf(os.Stderr, "Error: ")
		fmt.Fprintln(os.Stderr, err)
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
	cliTitle.Println("AdultBlocker · macOS")
	cliMuted.Println("Block sites via /etc/hosts and a local block page server.")
	fmt.Println()
	cliHeading.Println("Commands")
	w := tabwriter.NewWriter(color.Output, 0, 0, 2, ' ', 0)
	rows := [][2]string{
		{"sudo blocker install", "Set up config, backup hosts"},
		{"sudo blocker enable", "Turn blocking on"},
		{"sudo blocker disable", "Turn blocking off"},
		{"sudo blocker status", "Show config and block list"},
		{"sudo blocker uninstall", "Remove hosts entries and config"},
		{"sudo blocker add-domain <domain>", "Add a domain to the list"},
		{"sudo blocker remove-domain <domain>", "Remove a domain"},
		{"sudo blocker daemon", "Run the local block page (keep running)"},
	}
	for _, row := range rows {
		cliInfo.Fprint(w, "  ", row[0], "\t")
		cliMuted.Fprint(w, row[1], "\n")
	}
	_ = w.Flush()
	fmt.Println()
	cliMuted.Println("Most commands require root (sudo).")
}

func printKV(key string, value interface{}) {
	cliKey.Fprintf(color.Output, "  %-22s ", key+":")
	fmt.Fprintln(color.Output, value)
}

func printBoolKV(key string, val bool) {
	cliKey.Fprintf(color.Output, "  %-22s ", key+":")
	if val {
		cliSuccess.Fprintln(color.Output, "yes")
	} else {
		cliDanger.Fprintln(color.Output, "no")
	}
}

func printCountKV(key string, n int) {
	cliKey.Fprintf(color.Output, "  %-22s ", key+":")
	if n > 0 {
		cliSuccess.Fprintf(color.Output, "%d\n", n)
	} else {
		cliDanger.Fprintf(color.Output, "%d\n", n)
	}
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

	cliSuccess.Println("Install complete.")
	fmt.Println()
	printKV("Config", config.ConfigPath())
	printKV("Hosts backup", "/etc/hosts.adultblocker.bak")
	fmt.Println()
	cliMuted.Println("Run sudo blocker enable to activate blocking.")
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

	cliSuccess.Println("Blocking enabled.")
	fmt.Println()
	cliMuted.Print("Start the block page in another terminal: ")
	cliInfo.Println("sudo blocker daemon")
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

	if err := reflectpause.Run(15, "turn off blocking"); err != nil {
		if errors.Is(err, reflectpause.ErrAborted) {
			cliWarn.Println("Nothing changed.")
			return nil
		}
		return err
	}

	cfg.Enabled = false
	if err := config.Save(cfg); err != nil {
		return err
	}

	if err := hosts.Remove(); err != nil {
		return err
	}

	cliDanger.Println("Blocking disabled.")
	return nil
}

func uninstall() error {
	if err := platform.EnsureRoot(); err != nil {
		return err
	}

	if err := reflectpause.Run(15, "uninstall the blocker"); err != nil {
		if errors.Is(err, reflectpause.ErrAborted) {
			cliWarn.Println("Nothing changed.")
			return nil
		}
		return err
	}

	if err := hosts.Remove(); err != nil {
		return err
	}

	_ = os.Remove(config.ConfigPath())

	cliSuccess.Println("Uninstall complete.")
	fmt.Println()
	cliMuted.Println("Managed hosts entries removed.")
	cliMuted.Println("Config file removed.")
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

	fmt.Println()
	cliTitle.Println("Status")
	fmt.Println()
	printBoolKV("Blocking on", cfg.Enabled)
	printBoolKV("Hosts entries applied", present)
	printCountKV("Blocked domains", len(cfg.BlockedDomains))
	printKV("Config path", config.ConfigPath())
	printKV("Block page port", cfg.BlockPagePort)
	if len(cfg.BlockedDomains) > 0 {
		fmt.Println()
		cliHeading.Println("Blocked domains")
		for _, d := range cfg.BlockedDomains {
			cliMuted.Print("  · ")
			cliSuccess.Println(d)
		}
	}
	fmt.Println()
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

	cliSuccess.Print("Added to block list: ")
	cliSuccess.Println(domain)
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

	if err := reflectpause.Run(15, fmt.Sprintf("remove %q from your block list", domain)); err != nil {
		if errors.Is(err, reflectpause.ErrAborted) {
			cliWarn.Println("Nothing changed.")
			return nil
		}
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

	cliDanger.Print("Removed from block list: ")
	cliDanger.Println(domain)
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

	if !cfg.Enabled {
		cliDanger.Println("Blocking is off in config — enable for hosts-based protection.")
	}
	cliSuccess.Print("Block page server · ")
	cliInfo.Printf("http://127.0.0.1:%d\n", cfg.BlockPagePort)
	cliMuted.Println("(Press Ctrl+C to stop.)")
	fmt.Println()
	server := &blockpage.Server{Port: cfg.BlockPagePort}
	return server.Start()
}
