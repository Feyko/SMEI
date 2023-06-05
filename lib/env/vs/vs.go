package vs

import (
	"SMEI/lib/colors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Info struct {
	Components []string
}

func Install(path string, avoidVsReinstall bool) error {
	if avoidVsReinstall {
		// TODO move this to a better part of the process
		colors.SequenceColor.Println("Skipping installing Visual Stuido due to user-selected config option")
		return nil
	}

	colors.SequenceColor.Printf("Installing Visual Studio at: %s\n", path)

	targetPath, err := filepath.Abs(path)

	filename, err := downloadInstaller()
	if err != nil {
		return fmt.Errorf("could not download the VS installer: %v", err)
	}

	configString, err := makeConfigString(targetPath)
	if err != nil {
		return fmt.Errorf("could not make the VS installer config string: %v", err)
	}

	configFilename := filename + ".conf.json"

	err = os.WriteFile(configFilename, configString, 0666)
	if err != nil {
		return fmt.Errorf("could not create the VS installer configuration file: %v", err)
	}

	cmd := exec.Command(filename, "--wait", "--in", configFilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil && !isRebootExitCode(err) {
		return fmt.Errorf("error while running the VS installer: %v", err)
	}

	return nil
}

func makeConfigString(targetPath string) ([]byte, error) {
	config, err := makeConfig(targetPath)
	if err != nil {
		return nil, fmt.Errorf("could not make the config: %v", err)
	}

	r, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("could not marshal the json: %v", err)
	}
	return r, nil
}

func makeConfig(targetPath string) (map[string]interface{}, error) {
	config := defaultConfigObject()
	config["installPath"] = filepath.Join(targetPath, "VisualStudio")
	return config, nil
}

func downloadInstaller() (string, error) {
	link := "https://aka.ms/vs/17/release/vs_community.exe"
	filename := filepath.Join(os.TempDir(), "vs_Community.exe")
	resp, err := http.Get(link)
	if err != nil {
		return "", fmt.Errorf("could not get the installer file: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not download the installer file: %v", err)
	}

	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return "", fmt.Errorf("could not save the installer file: %v", err)
	}

	return filename, nil
}

func defaultConfigObject() map[string]interface{} {
	return map[string]interface{}{
		"productId":      "Microsoft.VisualStudio.Product.Community",
		"channelUri":     "https://aka.ms/vs/17/release/channel",
		"addProductLang": []string{"en-US"},
		"add": []string{
			"Microsoft.VisualStudio.Workload.NativeDesktop",
			"Microsoft.VisualStudio.Workload.NativeGame",
			"Microsoft.Net.Component.4.8.SDK",
			"Microsoft.VisualStudio.Component.Windows10SDK.20348",
		},
		"passive":   true,
		"force":     true,
		"norestart": true,
	}
}

func isRebootExitCode(err error) bool {
	code, ok := err.(*exec.ExitError)
	return !ok || code.ExitCode() == 3010
}
