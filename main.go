package main

import (
	_ "embed"
	"vpilot-patch-utility/patcher"
)

//go:embed patchfile.yml
var PatchFileData []byte

func main() {
	patcher.Entrypoint(PatchFileData)
}
