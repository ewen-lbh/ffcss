# ffcss

---

**Warning** This project is in a "I'm figuring out the interface" phase. The README is meant to represent what the final product will look like.

---

A CLI interface to apply and configure [Firefox CSS themes](https://reddit.com/r/FirefoxCSS) (also known as userChrome.css themes).

## Installation

### Pre-compiled binary 

```sh
# Install the latest release by downloading the binary on Github
curl -LO https://github.com/ewen-lbh/ffcss/releases/latest/download/ffcss
# Make sure the file is marked as executable
chmod a+x ffcss
# Move it to a folder that's in your path (so you can type `ffcss` anywhere), eg.
mv ffcss ~/.local/bin/ffcss
```

### Compile it yourself

```sh
git clone https://github.com/ewen-lbh/ffcss
cd ffcss
make 
make tests # optional, to make sure everything works
make install
```

## Usage

```docopt
ffcss - Apply and configure FirefoxCSS themes

Usage:
    ffcss use THEME_NAME
    ffcss reapply
    ffcss init 

Where:
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

- It'll download the zip file at the URL found in this repo's `themes/*.yaml`  files

And if `THEME_NAME` is an URL:

- It'll download the zip file at `THEME_NAME`

Some config keys need to be changed before applying a theme. `toolkit.legacyUserProfileCustomizations.stylesheets` must be set to `true` for _all_ themes, but most require their own extra config keys. Those can be set in the project's `ffcss.yaml`, but, don't worry, if the theme you use do not include a `ffcss.yaml` file, it might be in this repository's `themes/*.yaml` files

<!-- ### The `config` command

Synopsis: `ffcss config KEY [VALUE]`

Much simpler than the `use` command, this one just adds convenience to set `about:config` keys. If `VALUE` is not provided, ffcss will output the specified `KEY`'s current value. -->

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

You can use `userChrome`, `userContent` and `user.js` keys to specify where those files are in your repo. You can use **case-insensitive** [glob patterns][globster]

If not specified, their default values are either `userChrome.css` (or `userContent.css`, or `user.js`) or `null` when the default file is not found.

Note that keys declared in `config` will be appended to the copied `user.js`.

You can also declare other files in `assets`, an array of [glob patterns][globster]. They take precedence over the others keys, since they get copied last.

You can use `{{ os }}`, which will get replaced with one of `windows`, `linux` or `macos`, and `{{ variant }}`, which will get replaced by the variant the user has chosen.

All files will get copied to `<user's profile folder>/chrome/`. You can change the destination folder (relative to `<user's profile folder>`) with `copy to`:

```yaml
ffcss: 0 # signals that no compatibility is ensured (since 0.X.X versions can contain breaking changes, see semver)

repository: https://github.com/muckSponge/MaterialFox
config:
  svg.context-properties.content.enabled: true
  browser.tabs.tabClipWidth: 83
  materialFox.reduceTabOverflow: true
  security.insecure_connection_text.enabled: true

assets: chrome/**
copy to: ./ # relative to user's profile folder
```

without the `copy to`, files would get copied to `<user's profile folder>/chrome/chrome/...`, as `chrome/` will be a part of the file names.

### Variants

Some themes allow users to choose between different variations. To declare them, add an object with key `variants`, that maps variant names to a configuration object, overriding `config`, `userContent`, etc. for that variation. An additional `description` key is available and will be shown to users when selecting a variant.

Note that overriding `config` only overrides values set, it does not remove configuration keys that have been set globally: with the following manifest:

```yaml
ffcss: 1 # signals that the manifest works with ffcss versions 1.X.X

config:
    one.property: yes
    another.property: buckaroo

variants:
    blue:
        config:
            one.property: false
```

choosing the variant "blue" will apply the following config:

```yaml
one.property: false
another.property: buckaroo
```

### Example

A configuration file example for [@MiguelRAvila](https://github.com/MiguelRAvila)'s [SimplerentFox](https://github.com/MiguelRAvila/SimplerentFox):

```yaml
ffcss: 1

config:
    layers.acceleration.force-enabled: true
    gfx.webrender.all: true
    svg.context-properties.content.enabled: true

userContent: ./{{ os }}/userContent.css

variants:
    OneLine:
        description: Puts everything onto a single line
        userChrome:
            ./{{ os }}/userChrome__OneLine.css
    WithURLBar:
        description: Include the URL bar
        userChrome:
            ./{{ os }}/userChrome__WithURLBar.css
    WithoutURLBar:
        description: Do not include a URL bar
        userChrome:
            ./{{ os }}/userChrome__WithoutURLBar.css
```

#### Using {{ variant }}

The above manifest could be simplified to:


```yaml
ffcss: 1

config:
    layers.acceleration.force-enabled: true
    gfx.webrender.all: true
    svg.context-properties.content.enabled: true

userContent: ./{{ os }}/userContent.css
userChrome: ./{{ os }}/userChrome__{{ variant }}.css

variants:
    OneLine:
        description: Puts everything onto a single line
    WithURLBar:
        description: Include the URL bar
    WithoutURLBar:
        description: Do not include a URL bar
```

If you don't want to write descriptions for variants and don't need to override anything, use `{}` as the variant's value:

```yaml
ffcss: 1

userChrome: ./userChrome--{{ variant }}.css

variants:
    blue: {}
    yellow: {}
    red: {}
```

[globster]: https://globster.xyz/
