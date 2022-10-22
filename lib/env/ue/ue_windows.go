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
	if errors.Is(err, registry.ErrNotExist) {
		return false, nil
	}
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
	current, err = filepath.Abs(current)
	if err != nil {
		return false, errors.Wrap(err, "could not make current path absolute")
	}

	return current == installPath, nil
}

func hasOtherInstall() (bool, error) {
	_, err := openSetupKey(registry.QUERY_VALUE)
	if errors.Is(err, registry.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "could not open the UE setup registry key")
	}
	return true, nil
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

func disableUninstaller() error {
	key, err := openSetupKey(registry.ALL_ACCESS)
	if err != nil {
		return errors.Wrap(err, "could not open the UE setup registry key")
	}
	err = key.SetStringValue("UninstallString", "")
	if err != nil {
		return errors.Wrap(err, "could not empty the uninstall string")
	}

	return nil
}

func getUninstallString() (string, error) {
	key, err := openSetupKey(registry.ALL_ACCESS)
	if err != nil {
		return "", errors.Wrap(err, "could not open the UE setup registry key")
	}
	uninstallString, _, err := key.GetStringValue("UninstallString")
	if err != nil {
		return "", errors.Wrap(err, "could not get the registry value")
	}
	return uninstallString, nil
}

func reenableUninstall(uninstallString string) error {
	fmt.Println("Reenabling uninstaller")
	key, err := openSetupKey(registry.ALL_ACCESS)
	if err != nil {
		return errors.Wrap(err, "could not open the UE setup registry key")
	}
	err = key.SetStringValue("UninstallString", uninstallString)
	if err != nil {
		return errors.Wrap(err, "could not empty the uninstall string")
	}
	return nil
}
