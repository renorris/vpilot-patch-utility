## vPilot Patch Utility

Patch utility for the vPilot pilot client enabling custom FSD server compatibility.

Currently supported versions:
> `3.11.1 | vPilot.exe SHA1 = 19edcce42b0f9dddc0de0c5cf5c79ed1c7018728`

## Features:

- Arbitrary contiguous binary patching
- Programmatically modify existing strings in #US heap
- Modify obfuscated vPilotConfig.xml values

## Configuration:

See the provided default `patchfile.yml`.

Most FSDers will want to configure the following fields:

1. `name` field. Configures the name of a custom patch. This message will be displayed to anyone applying the patch.
2. `custom fsd-jwt URL` field in the `userstring_patches` section.
   The value for this patch should be updated to include the full "fsd-jwt" auth URL.
   **NOTE: This must be a TLS-enabled HTTPS URL.**
   For an [openfsd](https://github.com/renorris/openfsd) server, this value should look something like the following:
```
https://youropenfsdserver.com/api/v1/fsd-jwt
```

3. `network_status_url` field in the `config_patches` section.
   This field points vPilot to a custom status.txt URL, that should look something like [this](https://status.vatsim.net).
   For an [openfsd](https://github.com/renorris/openfsd) server, the value should look something like this:
```
https://youropenfsdserver.com/api/v1/data/status.txt
```

4. Optionally, you can encode a list of cached server URLs in the `cached_servers` section.
   This effectively ignores `network_status_url` in regards to obtaining an FSD server list. 
   This does NOT, however, ignore the fsd-jwt auth requirement.

## Usage:

One who wishes to build a custom patch and send it to their friends should start by installing the [Go Programming Language](https://go.dev/dl/).
Once installed, clone this repository, edit the patchfile as required, then open a shell/CMD prompt in the repository's directory.

Note: a patchfile is automatically embedded when compiled into a patch .EXE file.
This means that anyone who obtains an .EXE compiled with a given patchfile.yml will automatically have access to those specific patches.

To build a patch for a given patchfile.yml:

Windows:
```
go build -o my-vpilot-patch.exe .
```

UNIX:
```
GOOS=windows GOARCH=amd64 go build -o my-vpilot-patch.exe .
```

Try it out before sending it to anyone, and make any changes as necessary.
The .exe will start by listing the patch information, then will prompt the user to apply them.
**To revert a patch, run the executable again.**

There are no secrets or magical binaries/diffs behind this project. Read the source code!

### Additional Notes

- The default `disable voice server connections` patch inhibits vPilot from initiating any VATSIM voice server connections. i.e. it disables AFV entirely.
- The patcher creates backups of vPilot.exe and vPilotConfig.xml as vPilot.exe.orig and vPilotConfig.xml.orig respectively. If something blows up, one can manually rename these files back to their original names. Worst case: reinstall vPilot.
