package patcher

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Entrypoint(patchFileData []byte) {
	// Defer wait for an [ENTER] keypress before exiting process
	defer bufio.NewReader(os.Stdin).ReadBytes('\n')

	fmt.Println("\nvpilot-patch-utility\n")

	var patchFile *PatchFile
	var err error
	if patchFile, err = ParsePatchfile(patchFileData); err != nil {
		log.Println("Error parsing patch file: " + err.Error())
		return
	}

	if IsVPilotRunning(patchFile) {
		return
	}

	// Check if a .orig file exists for the executable. If so, revert the patch and exit
	var reverted bool
	if reverted, err = DoRevert(patchFile); err != nil {
		return
	}

	if reverted {
		fmt.Println("\nAll patches reverted.")
		return
	}

	if err = VerifyExecutableChecksum(patchFile); err != nil {
		return
	}

	PrintPatchInformation(patchFile)

	// Wait for user confirmation before executing patch
	fmt.Println("Press [ENTER] to apply patch, ALT-F4 to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Make executable backup
	if err = MakeExecutableBackup(patchFile); err != nil {
		return
	}

	// Do config patches
	if err = DoConfigFilePatches(patchFile); err != nil {
		return
	}

	// Do userstring patches
	if err = DoUserstringPatches(patchFile); err != nil {
		return
	}

	// Do simple patches
	if err = DoSimplePatches(patchFile); err != nil {
		return
	}

	fmt.Println("\nPatch complete.")
}
