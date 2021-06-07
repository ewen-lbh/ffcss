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
    ffcss init 

Where:
    KEY         a setting key (see firefox's about:config)
    THEME_NAME  a theme name or URL (see README.md)
```

#### The `use` command

Synopsis: `ffcss use THEME_NAME`

ffcss will first search for a folder at `~/.config/ffcss/themes/THEME_NAME`. If not found, it will try to download the theme:

If `THEME_NAME` is of the form `OWNER/REPO`:

- It'll try to download a folder named `chrome` from the repository `github.com/OWNER/REPO`

If `THEME_NAME` is of the form `DOMAIN.TLD/PATH`:

- It'll download the zip file at `https://DOMAIN.TLD/PATH`
  
If `THEME_NAME` is of the form `NAME`:

- It'll download the zip file at the URL found in this repo's `themes.yaml`  file

And if `THEME_NAME` is an URL:

- It'll download the zip file at `THEME_NAME`

Some config keys need to be changed before applying a theme. `toolkit.legacyUserProfileCustomizations.stylesheets` must be set to `true` for _all_ themes, but most require their own extra config keys. Those can be set in the project's `ffcss.yaml`, but, don't worry, if the theme you use do not include a `ffcss.yaml` file, it might be in this repository's `themes.yaml`,

### The `config` command

Synopsis: `ffcss config KEY [VALUE]`

Much simpler than the `use` command, this one just adds convenience to set `about:config` keys. If `VALUE` is not provided, ffcss will output the specified `KEY`'s current value.

### The `reapply` command

Synopsis: `ffcss reapply`

This is the same as doing `ffcss use` with the current theme, useful when firefox updates.

### The `init` command

Synopsis: `ffcss init`

Creates a [`ffcss` manifest file](#creating-a-firefoxcss-theme) in the current directory

## Creating a FirefoxCSS theme

So that your users can benefit from the simple installation process provided by ffcss, you can add a `ffcss.yaml` file in the root of your project and declare additional configuration you might need. 

Note that `toolkit.legacyUserProfileCustomizations.stylesheets` is set to `true` automatically, no need to declare it.

### Config

An object mapping `about:config` configuration keys to their values:

```yaml
config:
    svg.context-properties.content.enabled: true
    security.insecure_connection_text.enabled: true
```

### Files

By default, all files from your repository's `chrome` folder are copied over to the user's profile directory.

If you want to only copy certain files, you can set `files`, an array of [glob patterns](https://globster.xyz/). 

`files` can also be an object where the keys are operating systems (`windows`, `linux` or `macos`) and the values are arrays of glob patterns:

```yaml
files:
    - userContent.css
    windows:
    - windows/userChrome.css
    linux:
    - linux/userChrome.css
    macos:
    - macos/userChrome.css
```

If your project is somewhat structured, you can use `{{ os }}`, which will get replaced with one of `windows`, `linux` or `macos`. 

```yaml
files:
    - userContent.css
    - '{{ os }}/userContent.css' # quotes needed if your string starts with {
```

### Variants

Some themes allow users to choose between different variations. Declare the available variants' names in `variants`, an array of strings or objects mapping the name to its description.

Then, in `files`, reference the variant's name with `{{ variant }}`.

A configuration file example for [@MiguelRAvila](https://github.com/MiguelRAvila)'s [SimplerentFox](https://github.com/MiguelRAvila/SimplerentFox):

```yaml
ffcss: 1

config:
    layers.acceleration.force-enabled: true
    gfx.webrender.all: true
    svg.context-properties.content.enabled: true

variants:
    - WithoutURLBar
    - WithURLBar 
    - OneLine: Merged tab & address bars

files:
    - ./{{ os }}/userContent.css
    - ./{{ os }}/userChrome__{{ variant }}.css
```
