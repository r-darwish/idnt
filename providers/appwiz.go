//go:build windows

package providers

import (
	"github.com/google/shlex"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"strings"
)

type AppWiz struct {
}

func (a *AppWiz) GetName() string {
	return "Application Wizard"
}

func (a *AppWiz) getApplications(rootKey registry.Key, registryKey string) ([]Application, error) {
	key, err := registry.OpenKey(
		rootKey,
		registryKey,
		registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	var result []Application

	keys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return nil, err
	}

	for _, subKey := range keys {
		app, err := a.getApp(rootKey, registryKey+"\\"+subKey)
		if err == nil {
			result = append(result, app)
		}
	}

	return result, nil
}

func (a *AppWiz) GetApplications() ([]Application, error) {
	var result []Application

	mainApplications, err := a.getApplications(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	if err != nil {
		return nil, err
	}
	result = append(result, mainApplications...)

	userApplications, err := a.getApplications(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	if err != nil {
		return nil, err
	}
	result = append(result, userApplications...)

	wow64Applications, err := a.getApplications(registry.LOCAL_MACHINE, "SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	if err != nil {
		return nil, err
	}
	result = append(result, wow64Applications...)

	return result, nil
}

type AppWizExtendedInfo struct {
	uninstallString string
}

func (a *AppWiz) getApp(rootKey registry.Key, keyName string) (Application, error) {
	var result Application
	result.Provider = a

	key, err := registry.OpenKey(
		rootKey,
		keyName,
		registry.QUERY_VALUE)
	if err != nil {
		return result, err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	name, _, err := key.GetStringValue("DisplayName")
	if err != nil {
		return result, err
	}
	result.Name = name

	uninstallString, _, err := key.GetStringValue("UninstallString")
	if err != nil {
		return result, err
	}

	result.ExtraInfo = make(map[string]string)
	storeValue := func(field string, valueKey string) {
		value, _, err2 := key.GetStringValue(valueKey)
		if err2 == nil && value != "" {
			result.ExtraInfo[field] = value
		}
	}

	storeValue("URL", "URLInfoAbout")
	storeValue("Publisher", "Publisher")
	storeValue("Version", "DisplayVersion")

	extendedInfo := AppWizExtendedInfo{uninstallString: uninstallString}
	result.ExtendedInfo = &extendedInfo

	return result, nil
}

func (a AppWiz) RemoveApplication(application *Application) error {
	extendedInfo := application.ExtendedInfo.(*AppWizExtendedInfo)
	commandLine := []string{extendedInfo.uninstallString}
	var err error
	if !strings.HasPrefix(commandLine[0], "C:\\Program Files") {
		commandLine, err = shlex.Split(strings.Replace(extendedInfo.uninstallString, "\\", "\\\\", -1))
		if err != nil {
			return err
		}
	}
	command := exec.Command(commandLine[0], commandLine[1:]...)
	return command.Run()
}
