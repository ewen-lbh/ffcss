# ffcss

---

**Warning** This project is in a "I'm figuring out the interface" phase. The README is meant to represent what the final product will look like.

---

A CLI interface to apply and configure [Firefox CSS themes](https://reddit.com/r/FirefoxCSS) (also known as userChrome.css themes).

## Installation

```shell
# Install the latest release by downloading the binary on Github
curl -LO https://github.com/ewen-lbh/ffcss/releases/latest/download/ffcss
# Make sure the file is marked as executable
chmod a+x ffcss
# Move it to a folder that's in your path (so you can type `ffcss` anywhere), eg.
mv ffcss ~/.local/bin/ffcss
```

## Usage

```docopt
ffcss - Apply and configure FirefoxCSS themes

Usage:
    ffcss configure KEY [VALUE]
    ffcss use THEME_NAME
    ffcss reapply
    ffcss init [FORMAT]

Where:
    KEY         a setting key (see firefox's about:config)
    THEME_NAME  a theme name or URL (see README.md)
    FORMAT      one of "json", "yaml" or "toml"
```

#### The `use` command

Synopsis: `ffcss use THEME_NAME`

ffcss will first search for a folder at `~/.config/ffcss/themes/THEME_NAME`. If not found, it will try to download the theme:

If `THEME_NAME` is of the form `OWNER/REPO`:

- It'll try to download a folder named `chrome` from the repository `github.com/OWNER/REPO`

If `THEME_NAME` is of the form `DOMAIN.TLD/PATH`:

- It'll download the zip file at `https://DOMAIN.TLD/PATH`
  
If `THEME_NAME` is of the form `NAME`:

- It'll download the zip file at the URL found in this repo's `themes.toml`  file

And if `THEME_NAME` is an URL:

- It'll download the zip file at `THEME_NAME`

Some config keys need to be changed before applying a theme. `toolkit.legacyUserProfileCustomizations.stylesheets` must be set to `true` for _all_ themes, but most require their own extra config keys. Thos can be set in the project's `ffcss.json`, but, don't worry, if the theme you use do not include a `ffcss.json` file, it might be in this repository's `themes.toml`,

### The `config` command

Synopsis: `ffcss config KEY [VALUE]`

Much simpler than the `use` command, this one just adds convinience to set `about:config` keys. If `VALUE` is not provided, ffcss will output the specified `KEY`'s current value.

### The `reapply` command

Synopsis: `ffcss reapply`

This is the same as doing `ffcss use` with the current theme, useful when firefox updates.

### The `init` command

Synopsis: `ffcss init [FORMAT]`

Creates a [`ffcss` manifest file](#creating-a-firefoxcss-theme) in the current directory, either `ffcss.json`, `ffcss.yaml` or `ffcss.toml` (dependeding on the value of `FORMAT`, which defaults to `yaml`.)

See where the arguments are placed:

```yaml
manifest_version: 1
files:
    # Declare what files to copy over...
    - chrome/**
config:
    # Add your configuration keys here...
```

## Creating a FirefoxCSS theme

So that your users can benefit from the simple installation process provided by ffcss, you can add a `ffcss.json`, `ffcss.yaml` or `ffcss.toml` file in the root of your project and declare additional configuration you might need. Note that `toolkit.legacyUserProfileCustomizations.stylesheets` is set to `true` automatically, no need to declare it.

By default, all files from your respository's `chrome` folder are copied over to the user's profile directory. If you want to only copy certain files, you can set `files`, an array of [glob patterns](https://globster.xyz/). `files` can also be an object where the keys are operating systems (`windows`, `linux` or `macos`) and the values are arrays of glob patterns.

A configuration file example for [@MiguelRAvila](https://github.com/MiguelRAvila)'s [SimplerentFox](https://github.com/MiguelRAvila/SimplerentFox):

```json
{
    "manifest_version": 1,
    "files": {
        "linux": [
            "Linux/userChrome__WithURLBar.comp.css",
            "Linux/userContent.css"
        ],
        "windows": [
            "Windows/userChrome__WithURLBar.css",
            "Windows/userContent.css"
        ]
    },
    "config": {
        "layers.acceleration.force-enabled": true,
        "gfx.webrender.all": true,
        "svg.context-properties.content.enabled": true,
    }
}
```

The same file using YAML syntax:

```yaml
manifest_version: 1
files:
    linux:
        - Linux/userChrome__WithURLBar.comp.css
        - Linux/userContent.css
    windows:
        - Windows/userChrome__WithURLBar.css
        - Windows/userContent.css
config:
    layers.acceleration.force-enabled: true
    gfx.webrender.all: true
    svg.context-properties.content.enabled: true
```

And using TOML syntax:

```toml
manifest_version = 1

[files]
linux = [
    "Linux/userChrome__WithURLBar.comp.css",
    "Linux/userContent.cs"
]
windows = [
    "Windows/userChrome__WithURLBar.css",
    "Windows/userContent.css"
]

[config]
"layers.acceleration.force-enabled" = true
"gfx.webrender.all" = true
"svg.context-properties.content.enabled" = true
```
