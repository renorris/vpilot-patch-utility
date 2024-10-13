package pe

import (
	"errors"
	peparser "github.com/saferwall/pe"
	"strings"
)

// GetFileOffset returns the file offset for a given #US offset
func GetFileOffset(filename string, userStringOffset uint32) (offset int64, err error) {
	// Open portable executable file
	var file *peparser.File
	if file, err = peparser.New(filename, &peparser.Options{}); err != nil {
		return
	}

	// parse it
	if err = file.Parse(); err != nil {
		return
	}

	// Close file
	if err = file.Close(); err != nil {
		return
	}

	// Find user string header
	userStringHeaderIndex := -1
	// Find #US metadata stream
	for i, stream := range file.CLR.MetadataStreamHeaders {
		if stream.Name == "#US" {
			userStringHeaderIndex = i
			break
		}
	}

	if userStringHeaderIndex == -1 {
		err = errors.New("didn't find a #US header in the CLR")
		return
	}

	// Find .text section
	if len(file.Sections) < 1 {
		err = errors.New("no file sections found")
		return
	}

	textSection := file.Sections[0]
	if !strings.HasPrefix(string(textSection.Header.Name[:]), ".text") {
		err = errors.New("first prefix should be .text")
		return
	}

	// Calculate offset
	offset = int64(file.CLR.CLRHeader.MetaData.VirtualAddress +
		file.CLR.MetadataStreamHeaders[userStringHeaderIndex].Offset +
		textSection.Header.PointerToRawData -
		textSection.Header.VirtualAddress)

	// Append caller's offset to base offset
	offset += int64(userStringOffset)

	return
}
