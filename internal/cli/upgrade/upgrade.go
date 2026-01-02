package upgrade

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/KashifKhn/haft/internal/logger"
	"github.com/spf13/cobra"
)

var currentVersion = "dev"

func SetCurrentVersion(v string) {
	currentVersion = v
}

func GetCurrentVersion() string {
	return currentVersion
}

type UpgradeResult struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	UpdateAvailable bool   `json:"update_available"`
	Upgraded        bool   `json:"upgraded"`
	Platform        struct {
		OS   string `json:"os"`
		Arch string `json:"arch"`
	} `json:"platform"`
	InstalledPaths []string `json:"installed_paths,omitempty"`
	Error          string   `json:"error,omitempty"`
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Haft to the latest version",
		Long: `Upgrade Haft CLI to the latest version from GitHub releases.

This command will:
  1. Check if a newer version is available
  2. Download the latest release for your platform
  3. Verify the download integrity using checksums
  4. Create a backup of the current binary
  5. Install the new version
  6. Verify the installation works
  7. Rollback automatically if anything fails

The upgrade is safe - if anything goes wrong, your current version
will be automatically restored.`,
		Example: `  # Check and upgrade to latest version
  haft upgrade

  # Only check for updates (don't install)
  haft upgrade --check

  # Force reinstall even if already on latest
  haft upgrade --force

  # Upgrade to a specific version
  haft upgrade --version v0.5.0

  # JSON output for scripting
  haft upgrade --check --json`,
		RunE: runUpgrade,
	}

	cmd.Flags().BoolP("check", "c", false, "Only check for updates without installing")
	cmd.Flags().BoolP("force", "f", false, "Force upgrade even if already on latest version")
	cmd.Flags().StringP("version", "v", "", "Upgrade to a specific version")
	cmd.Flags().Bool("json", false, "Output result as JSON")

	return cmd
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check")
	force, _ := cmd.Flags().GetBool("force")
	targetVersion, _ := cmd.Flags().GetString("version")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	log := logger.Default()
	result := &UpgradeResult{}

	platform, err := GetPlatformInfo()
	if err != nil {
		return outputError(jsonOutput, result, err)
	}
	result.Platform.OS = platform.OS
	result.Platform.Arch = platform.Arch

	currentVersion := GetCurrentVersion()
	result.CurrentVersion = NormalizeVersion(currentVersion)

	if !jsonOutput {
		log.Info("Checking for updates...")
	}

	var latestVersion string
	if targetVersion != "" {
		latestVersion = NormalizeVersion(targetVersion)
	} else {
		latestVersion, err = GetLatestVersion()
		if err != nil {
			return outputError(jsonOutput, result, fmt.Errorf("failed to check for updates: %w", err))
		}
	}
	result.LatestVersion = latestVersion

	isNewer, err := IsNewerAvailable(currentVersion, latestVersion)
	if err != nil {
		return outputError(jsonOutput, result, fmt.Errorf("failed to compare versions: %w", err))
	}
	result.UpdateAvailable = isNewer

	if !isNewer && !force {
		if jsonOutput {
			return outputJSON(result)
		}
		log.Success(fmt.Sprintf("You're already on the latest version (%s)", result.CurrentVersion))
		return nil
	}

	if checkOnly {
		if jsonOutput {
			return outputJSON(result)
		}
		if isNewer {
			log.Info("Current version", "version", result.CurrentVersion)
			log.Info("Latest version", "version", result.LatestVersion)
			log.Success("Update available! Run 'haft upgrade' to install.")
		} else {
			log.Info(fmt.Sprintf("You're on version %s (latest: %s)", result.CurrentVersion, result.LatestVersion))
		}
		return nil
	}

	if !jsonOutput {
		if isNewer {
			log.Info("Current version", "version", result.CurrentVersion)
			log.Info("Latest version", "version", result.LatestVersion)
		} else {
			log.Info("Forcing reinstall of", "version", result.LatestVersion)
		}
		log.Info("Platform", "os", platform.OS, "arch", platform.Arch)
	}

	if !jsonOutput {
		log.Info("Downloading", "version", latestVersion)
	}

	downloadResult, err := DownloadRelease(latestVersion, platform)
	if err != nil {
		return outputError(jsonOutput, result, fmt.Errorf("download failed: %w", err))
	}
	defer CleanupDownload(downloadResult)

	if !jsonOutput {
		log.Debug("Downloaded", "size", formatBytes(downloadResult.Size))
	}

	checksums, err := FetchChecksums(latestVersion)
	if err != nil {
		if !jsonOutput {
			log.Warning("Could not fetch checksums, skipping verification")
		}
	} else if checksums != nil {
		archiveName := platform.GetArchiveName(latestVersion)
		if expectedHash, ok := checksums[archiveName]; ok {
			if err := VerifyChecksum(downloadResult.FilePath, expectedHash); err != nil {
				return outputError(jsonOutput, result, fmt.Errorf("checksum verification failed: %w", err))
			}
			if !jsonOutput {
				log.Debug("Checksum verified")
			}
		}
	}

	if !jsonOutput {
		log.Info("Extracting binary...")
	}

	binaryPath, err := ExtractBinary(downloadResult.FilePath, platform)
	if err != nil {
		return outputError(jsonOutput, result, fmt.Errorf("extraction failed: %w", err))
	}

	installations := FindAllInstallations()
	if len(installations) == 0 {
		execPath, err := GetExecutablePath()
		if err != nil {
			installDir := GetInstallDir()
			installations = []string{fmt.Sprintf("%s/%s", installDir, platform.BinaryName)}
		} else {
			installations = []string{execPath}
		}
	}

	var backups []*BackupInfo
	for _, installPath := range installations {
		if _, err := os.Stat(installPath); err == nil {
			backup, err := CreateBackup(installPath)
			if err != nil {
				if !jsonOutput {
					log.Warning("Failed to create backup", "path", installPath, "error", err)
				}
			} else {
				backups = append(backups, backup)
				if !jsonOutput {
					log.Debug("Created backup", "path", installPath)
				}
			}
		}
	}

	defer func() {
		for _, backup := range backups {
			CleanupBackup(backup)
		}
	}()

	if !jsonOutput {
		log.Info("Installing new version...")
	}

	installed, installErrors := InstallToMultipleLocations(binaryPath, installations)
	result.InstalledPaths = installed

	if len(installed) == 0 {
		for _, backup := range backups {
			if restoreErr := RestoreBackup(backup); restoreErr != nil {
				if !jsonOutput {
					log.Error("Failed to restore backup", "error", restoreErr)
				}
			}
		}

		errMsg := "installation failed at all locations"
		if len(installErrors) > 0 {
			errMsg = fmt.Sprintf("%s: %v", errMsg, installErrors[0])
		}
		return outputError(jsonOutput, result, fmt.Errorf("%s", errMsg))
	}

	if !jsonOutput {
		log.Info("Verifying installation...")
	}

	verifyPath := installed[0]
	if err := VerifyInstallation(verifyPath); err != nil {
		if !jsonOutput {
			log.Warning("Verification failed, restoring backup...")
		}

		for _, backup := range backups {
			if restoreErr := RestoreBackup(backup); restoreErr != nil {
				if !jsonOutput {
					log.Error("Failed to restore backup", "error", restoreErr)
				}
			}
		}

		return outputError(jsonOutput, result, fmt.Errorf("installation verification failed: %w", err))
	}

	result.Upgraded = true

	if jsonOutput {
		return outputJSON(result)
	}

	log.Success(fmt.Sprintf("Successfully upgraded from %s to %s", result.CurrentVersion, result.LatestVersion))

	if len(installed) > 1 {
		log.Info("Updated installations:")
		for _, path := range installed {
			log.Info("  " + path)
		}
	}

	if len(installErrors) > 0 {
		log.Warning("Some installations failed:")
		for _, err := range installErrors {
			log.Warning("  " + err.Error())
		}
	}

	return nil
}

func outputError(jsonOutput bool, result *UpgradeResult, err error) error {
	if jsonOutput {
		result.Error = err.Error()
		return outputJSON(result)
	}
	return err
}

func outputJSON(result *UpgradeResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
