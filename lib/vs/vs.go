package vs

import (
	"encoding/json"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(path string) error {
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
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
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
	temp := os.TempDir()
	resp, err := grab.Get(temp, link)
	if err != nil {
		return "", fmt.Errorf("could not get the installer file: %v", err)
	}
	return resp.Filename, nil
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
		"passive": true,
		"force":   true,
	}
}
