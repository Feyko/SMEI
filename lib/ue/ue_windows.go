package ue

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
)

const ueInstallerID = "F9524AE1-57B8-4EB8-BC4D-951EF9656DC7"

var infoPath = fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\{%v}_is1`, ueInstallerID)

func isReinstall(installPath string) (bool, error) {
	key, err := openSetupKey(registry.QUERY_VALUE)
	if err != nil {
		return false, errors.Wrap(err, "could not open the UE setup registry key")
	}

	current, _, err := key.GetStringValue("InstallLocation")
	if err != nil {
		return false, errors.Wrap(err, "could not get the current install location")
	}
	installPath, err = filepath.Abs(installPath)
	if err != nil {
		return false, errors.Wrap(err, "could not make install path absolute")
	}
	current, err = filepath.Abs(installPath)
	if err != nil {
		return false, errors.Wrap(err, "could not make current path absolute")
	}

	return current == installPath, nil
}

func openSetupKey(access uint32) (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, infoPath, access)
	if err == nil {
		return key, err
	}

	key, err = registry.OpenKey(registry.LOCAL_MACHINE, infoPath, access)
	if err != nil {
		return registry.CURRENT_USER, err
	}
	return key, err
}

func disableUninstaller() (string, error) {

	fmt.Println("Disabling uninstaller")
	key, err := openSetupKey(registry.ALL_ACCESS)
	if err != nil {
		return "", errors.Wrap(err, "could not open the UE setup registry key")
	}
	uninstallString, _, err := key.GetStringValue("UninstallString")
	if err != nil {
		return "", errors.Wrap(err, "could not get the uninstall string")
	}

	err = key.SetStringValue("Uninstall", "")
	if err != nil {
		return "", errors.Wrap(err, "could not empty the uninstall string")
	}

	return uninstallString, nil
}

func reenableUninstall(uninstallString string) error {
	fmt.Println("Reenabling uninstaller")
	key, err := openSetupKey(registry.ALL_ACCESS)
	if err != nil {
		return errors.Wrap(err, "could not open the UE setup registry key")
	}
	err = key.SetStringValue("Uninstall", uninstallString)
	if err != nil {
		return errors.Wrap(err, "could not empty the uninstall string")
	}
	return nil
}
