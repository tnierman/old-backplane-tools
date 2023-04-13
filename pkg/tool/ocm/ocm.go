package ocm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tnierman/backplane-tools/pkg/source/github"
	gogithub "github.com/google/go-github/github"
)

// Tool implements the interface to manage the 'ocm-cli' binary
type Tool struct {
	source *github.GithubTool
}

func NewTool() *Tool {
	t := &Tool{
		source: github.NewGithubTool("openshift-online", "ocm-cli"),
	}
	return t
}

func (t *Tool) Name() string {
	return "ocm"
}

func (t *Tool) Install(rootDir string) error {
	fmt.Println("installing ocm now")
	ocmDir := filepath.Join(rootDir, "ocm")
	err := os.MkdirAll(ocmDir, os.FileMode(0755))
	if err != nil {
		return err
	}

	latest, err := t.source.FetchLatestRelease()
	if err != nil {
		return err
	}
	//fmt.Println("latest: ", latest)
	fmt.Println("GOOS: ", runtime.GOOS)
	fmt.Println("GOARCH: ", runtime.GOARCH)

	var checksumAsset gogithub.ReleaseAsset
	var ocmBinary gogithub.ReleaseAsset

assetLoop:
	for _, asset := range latest.Assets {
		// Exclude assets that do not match system OS
		if !strings.Contains(asset.GetName(), runtime.GOOS) {
			continue assetLoop
		}
		// Exclude assets that do not match system architecture
		if !strings.Contains(asset.GetName(), runtime.GOARCH) {
			continue assetLoop
		}

		if strings.Contains(asset.GetName(), "sha256") {
			if checksumAsset.GetName() != "" {
				return fmt.Errorf("detected duplicate ocm-cli checksum assets")
			}
			checksumAsset = asset
		} else {
			if ocmBinary.GetName() != "" {
				return fmt.Errorf("detected duplicate ocm-cli binary assets")
			}
			ocmBinary = asset
		}
	}

	if checksumAsset.GetName() == "" || ocmBinary.GetName() == "" {
		return fmt.Errorf("failed to find ocm-cli or it's checksum")
	}

	fmt.Println("Downloading")
	err = t.source.DownloadReleaseAssets([]gogithub.ReleaseAsset{checksumAsset, ocmBinary}, ocmDir)
	if err != nil {
		return fmt.Errorf("failed to download one or more assets: %w", err)
	}
	return nil
}

func (t *Tool) Configure() error {
	return nil
}

func (t *Tool) Remove() error {
	return nil
}
