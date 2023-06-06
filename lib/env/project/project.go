package project

import (
	"SMEI/config"
	"SMEI/lib/colors"
	"SMEI/lib/credentials"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/mircearoata/wwise-cli/lib/wwise"
	"github.com/mircearoata/wwise-cli/lib/wwise/client"
	"github.com/mircearoata/wwise-cli/lib/wwise/product"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

func projectExists(targetPath string) (bool, error) {
	_, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func Clone(targetPath string) error {
	exists, err := projectExists(targetPath)
	if err != nil {
		return errors.Wrap(err, "could not check if the project already exists")
	}
	if exists {
		colors.Sequence.Printf("Project already exists in '%s', skipping clone\n", targetPath)
		return nil
	} else {
		colors.Sequence.Printf("Cloning starter project to '%s' (this can take many minutes)...\n", targetPath)
		_, err := git.PlainClone(filepath.Join(targetPath, "SatisfactoryModLoader"), false, &git.CloneOptions{
			URL:      "https://github.com/SatisfactoryModding/SatisfactoryModLoader",
			Progress: os.Stdout,
		})
		return err
	}
}

func TargetPathToUProjectPath(targetPath string, useSmlMiddle bool) string {
	if useSmlMiddle {
		return filepath.Join(targetPath, "SatisfactoryModLoader", "FactoryGame.uproject")
	} else {
		return filepath.Join(targetPath, "FactoryGame.uproject")
	}
}

func makeUBTArguments(targetPath string) []string {
	return []string{
		"-projectfiles",
		"-game",
		"-rocket",
		"-progress",
		fmt.Sprintf("-project=%v", TargetPathToUProjectPath(targetPath, true)),
	}
}

func GenerateProjectFiles(targetPath, UEPath string) error {
	colors.Sequence.Println("Generating Visual Studio project files...")
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
	colors.Sequence.Println("Building Development Editor...")
	err := BuildDevEditor(targetPath, UEPath)
	if err != nil {
		return err
	}
	colors.Sequence.Println("Building Shipping...")
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
	r = append(r, "-Target="+TargetPathToUProjectPath(targetPath, true))
	r = append(r, "-WaitMutex", "-FromMsBuild")
	return r
}

func makeTargetArguments(shipping bool) []string {
	if shipping {
		return []string{"FactoryGame", "Win64", "Shipping"}
	}
	return []string{"FactoryGameEditor", "Win64", "Development"}
}

func Install(targetPath string, UEPath string, auth credentials.WwiseAuth) error {
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

func InstallWWise(uprojectPath string, auth credentials.WwiseAuth) error {
	sdkVersion := viper.GetString(config.WwiseSdkVersion_key)
	integrationVersion := viper.GetString(config.WwiseIntegrationVersion_key)
	colors.Sequence.Printf("Downloading Wwise sdk %s files...\n", sdkVersion)
	wwiseClient := client.NewWwiseClient()

	err := wwiseClient.Authenticate(string(auth.Email), string(auth.Password))
	if err != nil {
		return errors.Wrap(err, "authentication error. check your Wwise credentials")
	}

	sdk := product.NewWwiseProduct(wwiseClient, "wwise")
	sdkProductVersion, err := sdk.GetVersion(sdkVersion)
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

	colors.Sequence.Printf("Integrating Wwise %s files...\n", integrationVersion)
	err = wwise.IntegrateWwiseUnreal(uprojectPath, integrationVersion, wwiseClient)
	if err != nil {
		return errors.Wrap(err, "integration failed")
	}

	return nil
}
