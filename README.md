# SMEI

Satisfactory Modding Environment Installer. Automatically sets up the Satisfactory Modding Environment, and dependencies, for you. Work in progress.

## Usage

As with the regular Satisfactory modding development environment, **only Windows is supported**.

1. Set up a github account linked as an Epic Games Developer Account by [following the modding documentation](https://docs.ficsit.app/satisfactory-modding/latest/Development/BeginnersGuide/dependencies.html#_link_your_github_as_an_epic_games_developer_account).
2. Sign up for a [WWise account](https://www.audiokinetic.com/en/products/wwise/)
3. Download the latest version of SMEI from the [Releases](https://github.com/Feyko/SMEI/releases) page.

### Installing a Modding Environment

1. Open a powershell terminal in the folder you downloaded the installer to.
2. Run `.\SMEI install --target <path to where you want the project to live>` and follow its prompts

### Integrating Wwise

1. Have an existing modding project set up 
2. Open a powershell terminal in the folder you downloaded the installer to.
3. Run `.\SMEI integrate --target <path to existing starer project>` and follow its prompts

### Configuring

Configuration interface is WIP. You can change some behaviors, such as skipping UE install or Visual Studio install, by editing `%APPDATA%\SMEI\config.yaml`.

## Troubleshooting

- Temporary files and config files are located in `%APPDATA%\SMEI\` and `%LOCALAPPDATA%\SMEI\`.
- If you forget your password, delete the directories mentioned above to reset it.

## Development

### Dependencies

- [Go 1.19](https://go.dev/doc/install)
- IDE of Choice. Goland or VSCode suggested.

## Building

```bash
go build
```

Will produce `SMEI.exe` in the repo root directory.

Consider enabling the `smei-developer-mode` config option in `%APPDATA%\SMEI\config.yaml` to skip the password entry process.

With this setting enabled you can quickly test changes via the following example Powershell command:

```powershell
go build; ./smei install --target="C:\Git\SMEI_TEST"
```

### Testing

Testing is somewhat troublesome since the installer is meant to be run on a fresh system. The best way to test is to run the installer on a VM, or a fresh install of Windows.
