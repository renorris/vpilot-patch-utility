package patcher

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"vpilot-patch-utility/config"
	"vpilot-patch-utility/pe"
	"vpilot-patch-utility/pe/userstring"
)

func DoRevert(patchfile *PatchFile) (reverted bool, err error) {
	var dirEntry []os.DirEntry
	if dirEntry, err = os.ReadDir(patchfile.ExecutableDirectory); err != nil {
		fmt.Println("Error reading directory:\n" + err.Error())
		return
	}

	origFileExists := false
	for _, entry := range dirEntry {
		if entry.Name() == "vPilot.exe.orig" {
			origFileExists = true
			break
		}
	}

	if !origFileExists {
		return
	}

	// If the .orig files exist but the correct clean version of vPilot is installed (the user
	// most likely updated vPilot without reverting a patch), remove the .orig files and return.
	if err = VerifyExecutableChecksum(patchfile); err == nil {
		if err = os.Remove(patchfile.ExecutablePath + ".orig"); err != nil {
			fmt.Println("Error deleting backup executable file. Try reinstalling vPilot.\n" + err.Error())
			return
		}
		if err = os.Remove(patchfile.ConfigFilePath + ".orig"); err != nil {
			fmt.Println("Error deleting backup config file. Try reinstalling vPilot.\n" + err.Error())
			return
		}
		return
	}

	fmt.Println("Previous patch detected. Reverting...\n")

	// Revert executable
	if err = CopyFile(patchfile.ExecutablePath+".orig", patchfile.ExecutablePath); err != nil {
		fmt.Println("Error reverting config file. Try reinstalling vPilot.\n" + err.Error())
		return
	}

	// Delete backup executable file
	if err = os.Remove(patchfile.ExecutablePath + ".orig"); err != nil {
		fmt.Println("Error deleting backup executable file. Try reinstalling vPilot.\n" + err.Error())
		return
	}

	fmt.Println("Reverted executable.")

	// Revert config file
	if err = CopyFile(patchfile.ConfigFilePath+".orig", patchfile.ConfigFilePath); err != nil {
		fmt.Println("Error reverting config file. Try reinstalling vPilot.\n" + err.Error())
		return
	}

	// Delete backup config file
	if err = os.Remove(patchfile.ConfigFilePath + ".orig"); err != nil {
		fmt.Println("Error deleting backup config file. Try reinstalling vPilot.\n" + err.Error())
		return
	}

	fmt.Println("Reverted config file.")

	reverted = true

	return
}

func DoConfigFilePatches(patchfile *PatchFile) (err error) {
	// Make config file backup
	if err = CopyFile(patchfile.ConfigFilePath,
		patchfile.ConfigFilePath+".orig"); err != nil {
		return
	}

	// Attempt to open config file
	var configFile *os.File
	if configFile, err = os.OpenFile(patchfile.ConfigFilePath, os.O_RDWR, 0644); err != nil {
		fmt.Println("Error opening config file. Try reinstalling vPilot.\n" + err.Error())
		return
	}
	defer configFile.Close()

	// Copy config file into buffer
	var configFileBytes []byte
	if configFileBytes, err = io.ReadAll(configFile); err != nil {
		fmt.Println("Error copying config file into buffer:\n" + err.Error())
		return
	}

	// Obfuscate patch data before writing

	// Obfuscate network status URL
	var obfuscatedNetworkStatus []byte
	if obfuscatedNetworkStatus, err = config.ObfuscateToBase64([]byte(patchfile.ConfigPatches.NetworkStatusURL), config.ConfigObfuscatorKey); err != nil {
		fmt.Println("Error obfuscating network status:\n" + err.Error())
		return
	}
	patchfile.ConfigPatches.NetworkStatusURL = string(obfuscatedNetworkStatus)

	// Obfuscate each cached server
	for i, cachedServer := range patchfile.ConfigPatches.CachedServers {
		var obfuscatedData []byte
		if obfuscatedData, err = config.ObfuscateToBase64([]byte(cachedServer), config.ConfigObfuscatorKey); err != nil {
			fmt.Println("Error obfuscating cached server:\n" + err.Error())
			return
		}

		patchfile.ConfigPatches.CachedServers[i] = string(obfuscatedData)
	}

	// Update XML data with patch data
	var updatedXMLBytes []byte
	if updatedXMLBytes, err = config.UpdateXML(configFileBytes, patchfile.ConfigPatches.NetworkStatusURL, patchfile.ConfigPatches.CachedServers); err != nil {
		fmt.Println("Error updating XML value:\n" + err.Error())
		return
	}

	// Seek back to 0
	if _, err = configFile.Seek(0, io.SeekStart); err != nil {
		fmt.Println("Error seeking config file:\n" + err.Error())
		return
	}

	// Truncate original config file
	if err = configFile.Truncate(0); err != nil {
		fmt.Println("Error truncating config file:\n" + err.Error())
		return
	}

	// Write the updated XML data
	if _, err = io.Copy(configFile, bytes.NewReader(updatedXMLBytes)); err != nil {
		fmt.Println("Error copying updated XML into config file:\n" + err.Error())
		return
	}

	fmt.Println("Patched config file.\n")

	return
}

func MakeExecutableBackup(patchfile *PatchFile) (err error) {
	executablePath := patchfile.ExecutablePath
	backupExecutableFilePath := patchfile.ExecutablePath + ".orig"
	return CopyFile(executablePath, backupExecutableFilePath)
}

func DoUserstringPatches(patchfile *PatchFile) (err error) {

	executablePath := patchfile.ExecutablePath

	for _, patch := range patchfile.UserstringPatches {
		fmt.Printf("Applying patch \"%s\" ...\n", patch.Name)

		var fileOffset int64
		if fileOffset, err = pe.GetFileOffset(executablePath, patch.HeapOffset); err != nil {
			fmt.Println("Error finding file offset:\n" + err.Error())
			return
		}

		// Open executable file
		var f *os.File
		if f, err = os.OpenFile(executablePath, os.O_RDWR, 0644); err != nil {
			fmt.Println("Error opening executable file:\n" + err.Error())
			return
		}

		// Perform the patch
		if err = userstring.WriteUserString(f, fileOffset, patch.Value); err != nil {
			fmt.Println("Error writing user string into executable:\n" + err.Error())
			return
		}

		// Close file
		if err = f.Close(); err != nil {
			fmt.Println("Error closing executable file:\n" + err.Error())
			return
		}
	}

	return
}

func DoSimplePatches(patchfile *PatchFile) (err error) {
	executablePath := patchfile.ExecutablePath

	// Open executable file
	var f *os.File
	if f, err = os.OpenFile(executablePath, os.O_RDWR, 0644); err != nil {
		fmt.Println("Error opening executable file:\n" + err.Error())
		return
	}
	defer f.Close()

	for _, patch := range patchfile.SimplePatches {
		fmt.Printf("Applying patch \"%s\" ...\n", patch.Name)

		// Seek to patch offset
		if _, err = f.Seek(int64(patch.Offset), io.SeekStart); err != nil {
			fmt.Println("Error seeking executable file:\n" + err.Error())
			return
		}

		// Apply the patch
		if _, err = io.Copy(f, bytes.NewReader(patch.Data)); err != nil {
			fmt.Println("Error writing simple patch to executable file:\n" + err.Error())
			return
		}
	}

	return
}
