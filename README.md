# roge

_command line interface for Roblox package versioning system_ \
\
[![GitHub Release](https://img.shields.io/github/v/release/robinskaba/roge?style=flat-square)](https://github.com/robinskaba/roge/releases/latest)

Roge enables Roblox package development from popular IDEs instead of native Roblox Studio. Its aim is to serve as a middle ground between Studio and fully external workflows like Rojo and Wally. With Roge, you can build packages in editors like VSCode and immediately work with them in Roblox Studio without manual transitions.

## Installation

1. Install [the latest version of roge](https://github.com/robinskaba/roge/releases/latest) for your OS.
2. Open a terminal in the directory of the roge executable and run the following command:

```bash
roge setup
```

This installs roge in your user's program directory and adds roge to your PATH. After reopening the terminal, you should be able to use roge. Verify the installation by running `roge --version`.

Use `roge update` to update roge to the latest version.

## Setup

### Requirements

- Roblox [Open Cloud API key](https://create.roblox.com/dashboard/credentials?activeTab=ApiKeysTab)
  - Assets: Read, Write
  - Legacy Assets: Manage
- your Roblox user ID (found in your profile URL: `roblox.com/users/<id>/profile`)

### Configuration

```bash
roge config set --api-key <api_key> --user-id <user_id> --global
```

Use `--local` instead of `--global` to scope configuration to a specific repository.

## Usage

### Cloning a package

```bash
roge clone <asset_id>
```

### Publishing a package

```bash
roge init
roge push
```

- add `--name`, `--description` to overwrite default or current name / description
- `push` creates a new package if it's a new repository
- the package entry point is deduced automatically by looking for an `init.luau` file or a `.luau` file matching the directory name

### Pulling a package

```bash
roge pull
```

- all `luau` files will be overwritten by the latest package version
- requires asset ID to be set manually or as a result of a push command

#### Nested modules

Any nested modules are required to name their main module files `init.luau`:

<table>
<tr>
<th>Local</th>
<th>Roblox Studio</th>
</tr>
<tr>
<td>
<pre>
MainModuleName
├── init.luau
└── SubModule
    ├── init.luau
    └── Nested.luau
</pre>
</td>
<td>
<pre>
MainModuleName
└── SubModule
    └── Nested
</pre>
</td>
</tr>
</table>

## Commands

- `init` - initialize a repository in the current directory
- `config` - `--global`/`--local`
  - `config set` - set API key and user ID (`--api-key`, `--user-id`)
  - `config list` - show current configuration
- `asset`
  - `asset set` - update asset config (`--id`)
  - `asset view` - show current asset configuration
  - `asset reset` - reset asset configuration to defaults
- `push` - publish a new version (`--name`, `--description`)
- `pull` - overwrite local files with the latest version from Roblox
- `clone <asset_id>` - clone a package into a new local directory
- `log` - list all versions with timestamps
- `update` - update the roge binary to the latest version

Use `roge <command> help` for further information.

## Recommended workflow

Roge works best alongside the [Luau Language Server extension](https://marketplace.visualstudio.com/items?itemName=JohnnyMorganz.luau-lsp) in VSCode. For future compatibility, it's best to write all packages as if they were nested. The language server provides hints when importing nested modules, and its syntax works in Roblox Studio as well.

```
MyModule
│   .roge
│   init.luau                  // entry file
│   CoreUtils.luau             // next to main module entry file
└───NestedModule
    │   init.luau              // nested under main module entry file
    │   NestedUtils.luau       // next to nested module entry file
```

In the main module entry file, you can then use the following imports and their methods:

```luau
local CoreUtils = require("./MyModule/CoreUtils")
local NestedModule = require("./MyModule/NestedModule")
```

Since this setup does not require Rojo and we are dealing only with module script packages, it is recommended to set `"luau-lsp.sourcemap.autogenerate": false`, as sourcemaps are not used at all.
