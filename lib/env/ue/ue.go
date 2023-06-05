package ue

import (
	"SMEI/config"
	"SMEI/lib/colors"
	"SMEI/lib/env/gh"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
)

const orgName = "SatisfactoryModdingUE"
const repoName = "UnrealEngine"
const installerName = "UnrealEngine-CSS-Editor-Win64.exe"
const CacheFolder = "UE-Installer"

type Info struct {
	Version  string
	Location string
}

func Install(installDir, installerDir string, avoidUeReinstall bool) error {
	cached, err := installerIsCached()
	if err != nil {
		return errors.Wrap(err, "could not check if the installer is cached")
	}
	if !cached {
		colors.Sequence.Println("UE installer is not cached, downloading it. This will require GitHub authentication.")
		err = downloadInstaller(installerDir)
		if err != nil {
			return errors.Wrap(err, "could not download the installer")
		}
	} else {
		fmt.Printf("UE installer is cached in '%s'\n", getInstallerPath())
		fmt.Println("TODO: Need detection for how old the cached UE version is and/or if there is a newer version available")
	}

	err = runInstallerIfRequired(installerDir, installDir, avoidUeReinstall)
	if err != nil {
		return fmt.Errorf("could not run the Unreal Engine installer: %v", err)
	}

	return nil
}

func getInstallerPath() string {
	return filepath.Join(config.ConfigDir, CacheFolder, installerName)
}

func installerIsCached() (bool, error) {
	_, err := os.Stat(getInstallerPath())
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func downloadInstaller(path string) error {
	ctx := context.Background()
	client, err := gh.AuthedClient(ctx)
	if err != nil {
		return errors.Wrap(err, "error making a GitHub auth client")
	}

	err = ensureGithubAccess(ctx, client)
	if err != nil {
		return fmt.Errorf("could not ensure GitHub access: %v", err)
	}

	assetsToDownload, err := getAssetsToDownload(ctx, client)
	if err != nil {
		return fmt.Errorf("could not get the assets to download: %v", err)
	}

	err = os.MkdirAll(path, 0666)
	if err != nil {
		return fmt.Errorf("could not create the directories for the path '%v': %v", path, err)
	}

	for _, asset := range assetsToDownload {
		err := downloadAsset(ctx, client, asset, path)
		if err != nil {
			return fmt.Errorf("could not download asset '%v': %v", asset.GetName(), err)
		}
	}
	return nil
}

func getAssetsToDownload(ctx context.Context, client *github.Client) ([]*github.ReleaseAsset, error) {
	latestAssets, err := getLatestReleaseAssets(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("could not get the latest release's assets: %v", err)
	}

	assetsToDownload, err := filterAssets(latestAssets)
	if err != nil {
		return nil, fmt.Errorf("could not filter the assets to download: %v", err)
	}
	return assetsToDownload, nil
}

func getLatestReleaseAssets(ctx context.Context, client *github.Client) ([]*github.ReleaseAsset, error) {
	release, _, err := client.Repositories.GetLatestRelease(ctx, orgName, repoName)
	if err != nil {
		return nil, fmt.Errorf("could not get the latest release: %v", err)
	}

	assets, _, err := client.Repositories.ListReleaseAssets(ctx, orgName, repoName, release.GetID(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not list the release's assets")
	}

	return assets, nil
}

func ensureGithubAccess(ctx context.Context, client *github.Client) error {
	_, _, err := client.Repositories.Get(ctx, orgName, repoName)
	if err != nil {
		return fmt.Errorf("could not get the repo: %v\nNOTE: This is most likely because you haven't joined the Epic Games organization. Please refer to the docs:\nhttps://docs.ficsit.app/satisfactory-modding/latest/Development/BeginnersGuide/dependencies.html#_unreal_engine_4_custom_engine", err)
	}

	return nil
}

func downloadAsset(ctx context.Context, client *github.Client, asset *github.ReleaseAsset, dir string) error {
	assetName := asset.GetName()

	colors.Sequence.Printf("Downloading asset %v\n", assetName)
	data, err := getAssetData(ctx, client, asset)
	if err != nil {
		return fmt.Errorf("could not get data for asset '%v': %v", assetName, err)
	}

	colors.Sequence.Printf("Writing asset '%v' to disk\n", assetName)
	err = writeAssetFile(dir, assetName, data)
	if err != nil {
		return fmt.Errorf("could not write asset '%v' to disk: %v", assetName, err)
	}
	return nil
}

func getAssetData(ctx context.Context, client *github.Client, asset *github.ReleaseAsset) ([]byte, error) {
	content, _, err := client.Repositories.DownloadReleaseAsset(ctx, orgName, repoName, asset.GetID(), http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("could not start downloading asset '%v': %v", asset.GetName(), err)
	}
	defer content.Close()

	all, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, fmt.Errorf("could not read the content of asset '%v': %v", asset.GetName(), err)
	}

	return all, nil
}

func writeAssetFile(targetDir, assetName string, data []byte) error {
	filename := filepath.Join(targetDir, assetName)
	return os.WriteFile(filename, data, 0666)
}

func runInstallerIfRequired(installerDir, installDir string, avoidUeReinstall bool) error {
	reinstall := false
	other, err := hasOtherInstall()
	if err != nil {
		return errors.Wrap(err, "could not check if an Unreal Engine install already exists")
	}

	if other {
		reinstall, err = isReinstall(installDir)
		if err != nil {
			return errors.Wrap(err, "could not check if this is a reinstall")
		}
		if reinstall {
			colors.Sequence.Printf("An install already exists in '%s', this is a reinstall\n", installDir)
		} else {
			colors.Sequence.Println("The existing install appears unrelated TODO @feyko better message?")
		}
	}

	if reinstall && avoidUeReinstall {
		colors.Sequence.Println("Skipping installing Unreal Engine again due to user-selected config option")
		return nil
	}

	if other && !reinstall {
		err = disableUninstaller()
		if err != nil {
			return errors.Wrap(err, "could not disable the uninstaller")
		}
	}

	return runInstaller(installDir, installerDir)
}

func runInstaller(installDir, installerDir string) error {
	colors.Sequence.Println("Running the UE installer")
	filename := filepath.Join(installerDir, installerName)

	cmd := exec.Command(filename,
		"/SILENT",
		"/NORESTART",
		fmt.Sprintf(`/DIR=%v`, installDir),
	)

	return cmd.Run()
}

func filterAssets(assets []*github.ReleaseAsset) ([]*github.ReleaseAsset, error) {
	r := make([]*github.ReleaseAsset, 0)
	for _, asset := range assets {
		wanted, err := isAssetWanted(asset)
		if err != nil {
			return nil, fmt.Errorf("could not check if an asset is wanted:%v", err)
		}

		if wanted {
			r = append(r, asset)
		}
	}
	return r, nil
}

func isAssetWanted(asset *github.ReleaseAsset) (bool, error) {
	regex := `\.(exe|bin)`
	matched, err := regexp.MatchString(regex, asset.GetName())
	if err != nil {
		return false, fmt.Errorf("could not use regex on the asset name: %v", err)
	}
	return matched, nil
}
