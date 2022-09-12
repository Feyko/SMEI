package project

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/mircearoata/wwise-cli/lib/wwise"
	"github.com/mircearoata/wwise-cli/lib/wwise/client"
	"github.com/mircearoata/wwise-cli/lib/wwise/product"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
)

func Clone(targetPath string) error {
	_, err := git.PlainClone(filepath.Join(targetPath, "SatisfactoryModLoader"), false, &git.CloneOptions{
		URL:          "https://github.com/SatisfactoryModding/SatisfactoryModLoader",
		Progress:     os.Stdout,
		SingleBranch: true,
	})
	return err
}

func getUEPath(targetPath string, local bool) string {
	if local {
		return getLocalUEPath(targetPath)
	}
	return getGlobalUEPath()
}

func getLocalUEPath(targetPath string) string {
	return filepath.Join(targetPath, "Unreal Engine - CSS")
}

func getGlobalUEPath() string {
	return filepath.Join(os.ExpandEnv("$ProgramW6432"), "Unreal Engine - CSS")
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

func GenerateProjectFiles(targetPath string, local bool) error {
	UEPath := getUEPath(targetPath, local)

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

func BuildAll(targetPath string, local bool) error {
	err := BuildDevEditor(targetPath, local)
	if err != nil {
		return err
	}
	err = BuildShipping(targetPath, local)
	return err
}

func BuildShipping(targetPath string, local bool) error {
	return Build(targetPath, local, true)
}

func BuildDevEditor(targetPath string, local bool) error {
	return Build(targetPath, local, false)
}

func Build(targetPath string, local, shipping bool) error {
	UEPath := getUEPath(targetPath, local)
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

func Install(targetPath string, local bool) error {
	var err error
	err = Clone(targetPath)
	if err != nil {
		return fmt.Errorf("could not clone the project: %v", err)
	}

	err = InstallWWise(targetPath)
	if err != nil {
		return errors.Wrap(err, "could not move the Wwise install")
	}

	err = GenerateProjectFiles(targetPath, local)
	if err != nil {
		return fmt.Errorf("could not generate the VS project files: %v", err)
	}

	err = BuildAll(targetPath, local)
	if err != nil {
		return fmt.Errorf("could not build the project: %v", err)
	}

	return nil
}

func InstallWWise(targetPath string) error {
	wwiseClient := client.NewWwiseClient()

	email := viper.GetString("wwise-email")
	if email == "" {
		return errors.New("wwise-email was not set")
	}

	password := viper.GetString("wwise-password")
	if password == "" {
		return errors.New("wwise-password was not set")
	}

	err := wwiseClient.Authenticate(email, password)
	if err != nil {
		return errors.Wrap(err, "authentication error")
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

	err = wwise.IntegrateWwiseUnreal(targetPathToUProjectPath(targetPath), "2021.1.8.2285", wwiseClient)
	if err != nil {
		return errors.Wrap(err, "integration failed")
	}

	return nil
}
