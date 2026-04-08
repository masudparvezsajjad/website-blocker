package hosts

import (
	"fmt"
	"os"
	"strings"

	"github.com/masudparvezsajjad/website-blocker/internal/config"
)

const (
	HostsPath   = "/etc/hosts"
	BeginMarker = "# BEGIN ADULT_BLOCKER"
	EndMarker   = "# END ADULT_BLOCKER"
)

func ReadHosts() (string, error) {
	b, err := os.ReadFile(HostsPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func BackupHosts() error {
	content, err := os.ReadFile(HostsPath)
	if err != nil {
		return err
	}
	return os.WriteFile(HostsPath+".adultblocker.bak", content, 0644)
}

func buildManagedBlock(cfg *config.Config) string {
	var lines []string
	lines = append(lines, BeginMarker)

	seen := map[string]bool{}
	for _, d := range cfg.BlockedDomains {
		d = strings.TrimSpace(strings.ToLower(d))
		if d == "" || seen[d] {
			continue
		}
		seen[d] = true
		lines = append(lines, fmt.Sprintf("%s %s", cfg.RedirectIPv4, d))
		lines = append(lines, fmt.Sprintf("%s %s", cfg.RedirectIPv6, d))
	}

	lines = append(lines, EndMarker)
	return strings.Join(lines, "\n")
}

func stripManagedSection(content string) string {
	start := strings.Index(content, BeginMarker)
	end := strings.Index(content, EndMarker)

	if start == -1 || end == -1 || end < start {
		return strings.TrimRight(content, "\n")
	}

	end += len(EndMarker)
	before := strings.TrimRight(content[:start], "\n")
	after := strings.TrimLeft(content[end:], "\n")

	switch {
	case before == "" && after == "":
		return ""
	case before == "":
		return after
	case after == "":
		return before
	default:
		return before + "\n" + after
	}
}

func Apply(cfg *config.Config) error {
	content, err := ReadHosts()
	if err != nil {
		return err
	}

	cleaned := stripManagedSection(content)
	block := buildManagedBlock(cfg)

	var final string
	if strings.TrimSpace(cleaned) == "" {
		final = block + "\n"
	} else {
		final = cleaned + "\n\n" + block + "\n"
	}

	return atomicWrite(HostsPath, []byte(final), 0644)
}

func Remove() error {
	content, err := ReadHosts()
	if err != nil {
		return err
	}
	cleaned := stripManagedSection(content)
	if cleaned != "" && !strings.HasSuffix(cleaned, "\n") {
		cleaned += "\n"
	}
	return atomicWrite(HostsPath, []byte(cleaned), 0644)
}

func ManagedEntriesPresent() (bool, error) {
	content, err := ReadHosts()
	if err != nil {
		return false, err
	}
	return strings.Contains(content, BeginMarker) && strings.Contains(content, EndMarker), nil
}

func atomicWrite(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
