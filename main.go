package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "dnfDownloader <package-name>",
	Short: "Download RPMs for a package and its dependencies",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]
		run(packageName)
	},
}

func run(packageName string) {
	// Create output directory at $PWD/out/<package-name>
	outputDir, err := createOutputDirectory(packageName)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	fmt.Printf("Output directory created: %s\n", outputDir)

	// Resolve package dependencies using dnf repoquery.
	deps, err := resolveDependencies(packageName)
	if err != nil {
		log.Fatalf("Error resolving dependencies for %s: %v", packageName, err)
	}

	// Download each dependency package to the current working directory (cache).
	downloadPackages(deps)

	// Find all RPM files in the cache and move them to the output directory.
	moveRPMFiles(outputDir)
}

// createOutputDirectory creates $PWD/out/<package-name> for storing RPM files.
func createOutputDirectory(packageName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %v", err)
	}
	outputDir := filepath.Join(cwd, "out", packageName)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating output directory %s: %v", outputDir, err)
	}
	return outputDir, nil
}

// resolveDependencies uses dnf repoquery to fetch the package dependencies.
func resolveDependencies(pkgName string) ([]string, error) {
	cmd := exec.Command("dnf", "repoquery", "--resolve", "--requires", pkgName)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(output)), nil
}

// downloadPackages iterates over packages and downloads each using dnf download.
func downloadPackages(packages []string) {
	for _, pkg := range packages {
		if err := downloadPackage(pkg); err != nil {
			log.Printf("Failed to download %s: %v", pkg, err)
		} else {
			fmt.Printf("Successfully downloaded %s\n", pkg)
		}
	}
}

// downloadPackage downloads the given package.
func downloadPackage(pkgName string) error {
	cmd := exec.Command("dnf", "download", pkgName)
	return cmd.Run()
}

// moveRPMFiles searches for all RPM files in the current directory and moves them to outputDir.
func moveRPMFiles(outputDir string) {
	rpmFiles, err := filepath.Glob("*.rpm")
	if err != nil {
		log.Fatalf("Error searching for RPM files: %v", err)
	}
	for _, rpmFile := range rpmFiles {
		destPath := filepath.Join(outputDir, filepath.Base(rpmFile))
		if err := os.Rename(rpmFile, destPath); err != nil {
			log.Printf("Failed to move %s to %s: %v", rpmFile, outputDir, err)
		} else {
			fmt.Printf("Moved %s to %s\n", rpmFile, outputDir)
		}
	}
}
