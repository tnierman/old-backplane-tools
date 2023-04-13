package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

type GithubTool struct {
	// Owner specifies the organization or user this tool belongs to
	Owner string

	// Repo specifies the repository of the tool
	Repo string

	// client is used to interact with Github
	client *github.Client
}

func NewGithubTool(owner, repo string) *GithubTool {
	tool := &GithubTool{
		Owner: owner,
		Repo: repo,
		client: github.NewClient(nil),
	}
	return tool
}

// ListReleases returns all releases of the tool from Github
func (t GithubTool) ListReleases(opts *github.ListOptions) ([]*github.RepositoryRelease, error){
	releases, response, err := t.client.Repositories.ListReleases(context.TODO(), t.Owner, t.Repo, &github.ListOptions{})
	if err != nil {
		return []*github.RepositoryRelease{}, err
	}
	err = github.CheckResponse(response.Response)
	if err != nil {
		return []*github.RepositoryRelease{}, err
	}
	return releases, nil
}

// FetchRelease returns the specified release of the tool from Github
func (t GithubTool) FetchRelease(releaseID int64) (*github.RepositoryRelease, error) {
	release, response, err := t.client.Repositories.GetRelease(context.TODO(), t.Owner, t.Repo, releaseID)
	if err != nil {
		return &github.RepositoryRelease{}, err
	}
	err = github.CheckResponse(response.Response)
	if err != nil {
		return &github.RepositoryRelease{}, err
	}
	return release, nil
}

// FetchLatestRelease returns the latest release of the tool from Github
func (t GithubTool) FetchLatestRelease() (*github.RepositoryRelease, error) {
	release, response, err := t.client.Repositories.GetLatestRelease(context.TODO(), t.Owner, t.Repo)
	if err != nil {
		return &github.RepositoryRelease{}, err
	}
	err = github.CheckResponse(response.Response)
	if err != nil {
		return &github.RepositoryRelease{}, err
	}
	return release, nil
}

// DownloadReleaseAssets downloads the provided Github release assets and stores them in the given directory.
// The resulting files will match the assets' names
func (t GithubTool) DownloadReleaseAssets(assets []github.ReleaseAsset, dir string) error {
	for _, asset := range assets {
		fmt.Println("Attempting to download asset: ", asset.GetName())
		reader, redirectURL, err := t.client.Repositories.DownloadReleaseAsset(context.TODO(), t.Owner, t.Repo, asset.GetID())
		if err != nil {
			return err
		}
		defer func() {
			err = reader.Close()
			if err != nil {
				panic(fmt.Sprintf("failed to close reader from GitHub asset '%s'", asset.GetName()))
			}
		}()
		if redirectURL != "" {
			// TODO - Download from the redirect manually
			fmt.Println("redirect is set to: ", redirectURL)
			resp, err := http.Get(redirectURL)
			if err != nil {
				return err
			}
			reader = resp.Body
		}
		filePath := filepath.Join(dir, asset.GetName())
		fmt.Println("filepath: ", filePath)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		err = os.Chmod(file.Name(), os.FileMode(0755))
		if err != nil {
			return err
		}
		_, err = file.ReadFrom(reader)
		if err != nil {
			return err
		}
	}
	return nil
}
