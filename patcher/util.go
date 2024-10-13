package patcher

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func IsVPilotRunning(patchfile *PatchFile) (running bool) {
	var executableFile *os.File
	var err error
	if executableFile, err = os.OpenFile(patchfile.ExecutablePath, os.O_RDWR, 0644); err != nil {
		fmt.Println("Error opening executable file. Possible causes:\n- vPilot is currently open\n- vPilot is not installed\n\n" + err.Error())
		running = true
		return
	}
	defer executableFile.Close()

	return
}

func PrintPatchInformation(patchfile *PatchFile) {
	fmt.Printf("Loaded %d patches from \"%s\"\n\n",
		len(patchfile.SimplePatches)+len(patchfile.UserstringPatches), patchfile.Name)
}

func VerifyExecutableChecksum(patchfile *PatchFile) (err error) {
	f, err := os.Open(patchfile.ExecutableDirectory + string(filepath.Separator) + "vPilot.exe")
	if err != nil {
		log.Println("Error opening executable:\n" + err.Error())
		return err
	}

	// Calculate SHA1 sum of vpilot file
	h := sha1.New()
	if _, err = io.Copy(h, f); err != nil {
		log.Println(err)
		return
	}

	// Close file
	if err = f.Close(); err != nil {
		log.Println("error closing executable file: " + err.Error())
		return
	}

	if !slices.Equal([]byte(hex.EncodeToString(h.Sum(nil))), []byte(patchfile.ExpectedSum)) {
		fmt.Println("Error validating executable checksum.\nDo you have the right version installed?")
		return errors.New("checksum invalid")
	}

	return
}

func CopyFile(sourcePath, destinationPath string) (err error) {
	// Attempt to open source file
	var sourceFile *os.File
	if sourceFile, err = os.OpenFile(sourcePath, os.O_RDWR, 0644); err != nil {
		fmt.Println("Error opening file for copy. Try reinstalling vPilot.\n" + err.Error())
		return
	}
	defer sourceFile.Close()

	// Make file backup
	fmt.Printf("Copying %s ...\n", filepath.Base(sourcePath))
	var destinationFile *os.File
	if destinationFile, err = os.OpenFile(destinationPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		fmt.Println("Error opening destination file for copying:\n" + err.Error())
		return
	}
	defer destinationFile.Close()

	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		fmt.Println("Error copying file:\n" + err.Error())
		return
	}

	return
}
