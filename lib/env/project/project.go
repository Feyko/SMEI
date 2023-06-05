package project

import (
	"SMEI/lib/secret"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/mircearoata/wwise-cli/lib/wwise"
	"github.com/mircearoata/wwise-cli/lib/wwise/client"
	"github.com/mircearoata/wwise-cli/lib/wwise/product"
	"github.com/pkg/errors"
)

type Info struct {
	Location string
	Git      *GitInfo
}

type GitInfo struct {
	Branch   string
	Commit   string
	UpToDate bool
}

func Clone(targetPath string) error {
	fmt.Printf("Cloning starter project to '%s'...\n", targetPath)
	_, err := git.PlainClone(filepath.Join(targetPath, "SatisfactoryModLoader"), false, &git.CloneOptions{
		URL:      "https://github.com/SatisfactoryModding/SatisfactoryModLoader",
		Progress: os.Stdout,
	})
	return err
}

func targetPathToUProjectPath(targetPath string) string {
	return filepath.Join(targetPath, "SatisfactoryModLoader", "FactoryGame.uproject")
}

func makeUBTArguments(targetPath string) []string {
	return []string{
		"-projectfiles",
		"-game",
		"-rocket",
		"-progress",
		fmt.Sprintf("-project=%v", targetPathToUProjectPath(targetPath)),
	}
}

func GenerateProjectFiles(targetPath, UEPath string) error {
	fmt.Println("Generating Visual Studio project files...")
	UBTPath := filepath.Join(UEPath, "Engine", "Binaries", "DotNET", "UnrealBuildTool.exe")
	arguments := makeUBTArguments(targetPath)
	cmd := exec.Command(UBTPath, arguments...)
	fmt.Println(cmd)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("generation command failed: %v", err)
	}

	return nil
}

func BuildAll(targetPath, UEPath string) error {
	fmt.Println("Building Development Editor...")
	err := BuildDevEditor(targetPath, UEPath)
	if err != nil {
		return err
	}
	fmt.Println("Building Shipping...")
	err = BuildShipping(targetPath, UEPath)
	// TODO build dedicated servers
	return err
}

func BuildShipping(targetPath, UEPath string) error {
	return Build(targetPath, UEPath, true)
}

func BuildDevEditor(targetPath, UEPath string) error {
	return Build(targetPath, UEPath, false)
}

func Build(targetPath, UEPath string, shipping bool) error {
	buildScript := filepath.Join(UEPath, "Engine", "Build", "BatchFiles", "Build.bat")
	arguments := makeBuildArguments(targetPath, shipping)
	fmt.Println(buildScript, arguments)
	cmd := exec.Command(buildScript, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func makeBuildArguments(targetPath string, shipping bool) []string {
	var r []string
	r = append(r, makeTargetArguments(shipping)...)
	r = append(r, "-Target="+targetPathToUProjectPath(targetPath))
	r = append(r, "-WaitMutex", "-FromMsBuild")
	return r
}

func makeTargetArguments(shipping bool) []string {
	if shipping {
		return []string{"FactoryGame", "Win64", "Shipping"}
	}
	return []string{"FactoryGameEditor", "Win64", "Development"}
}

func Install(targetPath string, UEPath string, auth WwiseAuth) error {
	var err error
	err = Clone(targetPath)
	if err != nil {
		return fmt.Errorf("could not clone the project: %v", err)
	}

	err = InstallWWise(targetPath, auth)
	if err != nil {
		return errors.Wrap(err, "could not move the Wwise install")
	}

	err = GenerateProjectFiles(targetPath, UEPath)
	if err != nil {
		return fmt.Errorf("could not generate the VS project files: %v", err)
	}

	err = BuildAll(targetPath, UEPath)
	if err != nil {
		return fmt.Errorf("could not build the project: %v", err)
	}

	return nil
}

type WwiseAuth struct {
	Email    secret.String
	Password secret.String
}

func InstallWWise(targetPath string, auth WwiseAuth) error {
	fmt.Println("Installing Wwise files...")
	wwiseClient := client.NewWwiseClient()

	err := wwiseClient.Authenticate(string(auth.Email), string(auth.Password))
	if err != nil {
		return errors.Wrap(err, "authentication error. check your Wwise credentials")
	}

	sdk := product.NewWwiseProduct(wwiseClient, "wwise")
	sdkProductVersion, err := sdk.GetVersion("2021.1.8.7831")
	if err != nil {
		return errors.Wrap(err, "could not get SDK version")
	}

	sdkVersionInfo, err := sdkProductVersion.GetInfo()
	if err != nil {
		return errors.Wrap(err, "could not get SDK version info")
	}

	files := sdkVersionInfo.FindFilesByGroups([]product.GroupFilter{
		{GroupID: "Packages", GroupValues: []string{"SDK"}},
		{GroupID: "DeploymentPlatforms", GroupValues: []string{"Windows_vc140", "Windows_vc150", "Windows_vc160", "Mac", "Linux", ""}},
	})

	for _, file := range files {
		fmt.Printf("Downloading %v from %v\n", file.Name, file.URL)
		err = sdkProductVersion.DownloadOrCache(file)
		if err != nil {
			return errors.Wrapf(err, "could not download file %v", file.Name)
		}
	}

	fmt.Println("Integrating Wwise files...")
	err = wwise.IntegrateWwiseUnreal(targetPathToUProjectPath(targetPath), "2021.1.8.2285", wwiseClient)
	if err != nil {
		return errors.Wrap(err, "integration failed")
	}

	return nil
}
