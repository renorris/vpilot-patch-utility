package patcher

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os/user"
	"path/filepath"
	"strings"
)

// PatchFile represents the structure of the entire YAML configuration file.
type PatchFile struct {
	Name                string `yaml:"name"`
	ExecutableDirectory string `yaml:"executable_directory"`
	ExecutablePath      string
	ConfigFilePath      string
	ExpectedSum         string            `yaml:"expected_sum"`
	SimplePatches       []SimplePatch     `yaml:"simple_patches"`
	UserstringPatches   []UserstringPatch `yaml:"userstring_patches"`
	ConfigPatches       ConfigPatches     `yaml:"config_patches"`
}

// SimplePatch represents a single entry in the simple_patches list.
type SimplePatch struct {
	Name   string `yaml:"name"`
	Offset uint32 `yaml:"offset"` // string to preserve hex format
	Data   []byte `yaml:"data"`
}

// UserstringPatch represents a single entry in the userstring_patches list.
type UserstringPatch struct {
	Name       string `yaml:"name"`
	HeapOffset uint32 `yaml:"heap_offset"` // string to preserve hex format
	Value      string `yaml:"value"`
}

// ConfigPatches represents the configuration patches section.
type ConfigPatches struct {
	NetworkStatusURL string   `yaml:"network_status_url"`
	CachedServers    []string `yaml:"cached_servers"`
}

func ParsePatchfile(filedata []byte) (patchFile *PatchFile, err error) {
	yamlDecoder := yaml.NewDecoder(bytes.NewReader(filedata))

	patchFile = &PatchFile{}
	if err = yamlDecoder.Decode(patchFile); err != nil {
		return
	}

	// Replace {USERNAME} with current user
	var currentUser *user.User
	if currentUser, err = user.Current(); err != nil {
		return
	}
	patchFile.ExecutableDirectory = strings.Replace(patchFile.ExecutableDirectory, "{HOME}", currentUser.HomeDir, -1)

	// Initialize individual file paths
	patchFile.ExecutablePath = patchFile.ExecutableDirectory + string(filepath.Separator) + "vPilot.exe"
	patchFile.ConfigFilePath = patchFile.ExecutableDirectory + string(filepath.Separator) + "vPilotConfig.xml"

	return
}
