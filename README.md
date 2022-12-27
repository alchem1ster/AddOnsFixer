# AddOnsFixer
### _Fixes incorrectly unpacked addOns for you, quickly_

- Fixes read-only permissions for addOns
- Finds broken addOns and fixes them
- ✨Magic ✨

## Summary

<img src="https://i.ibb.co/ScymHmN/info.png" alt="drawing" width="500"/> 

## Attention

Before running this script for the first time, I STRONGLY RECOMMEND making a backup of your `Interface` folder: something bad is unlikely to happen, but just in case

## Installation

- Download the latest release binary (`AddOnsFixer.exe`) or build it yourself
- Place the file in the game folder, next to Wow.exe
- Run `AddOnsFixer.exe`

## Testing

- Place folders from the `tests` directory to your `./Interface/AddOns`
- Run `AddOnsFixer.exe`
- Check result inplace

## Manual build

Build from PowerShell for x86:

```powershell
$Env:GOARCH=386
go build -ldflags "-s -w" -o AddOnsFixer.exe
```

Build from CMD for x86:

```cmd
set GOARCH=386
go build -ldflags "-s -w" -o AddOnsFixer.exe
```
