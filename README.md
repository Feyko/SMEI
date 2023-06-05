# SMEI

Satisfactory Modding Environment Installer. Automatically sets up the Satisfactory Modding Environment, and dependencies, for you.

## Usage

As with the regular Satisfactory modding development environment, only Windows is supported.

1. Download the latest version from the [Releases](https://github.com/Feyko/SMEI/releases) page.
2. Open a powershell terminal in the folder you downloaded the installer to.
3. Sign up for a WWise account TODO: Link
4. Run the installer with `.\SMEI install --target <path to Satisfactory install>` and follow its prompts

## Troubleshooting

- Temporary files are generated in `%APPDATA%\SMEI\` and `%LOCALAPPDATA%\SMEI\`.
- If you forget your password, delete the directories mentioned above to reset it.

## Development

### Dependencies

- [Go 1.19](https://go.dev/doc/install)
- IDE of Choice. Goland or VSCode suggested.

## Building

```bash
go build
```

Will produce `SMEI.exe` in the current directory.

### Testing

Testing is somewhat troublesome since  the installer is meant to be run on a fresh system. The best way to test is to run the installer on a VM, or a fresh install of Windows.
