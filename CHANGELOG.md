# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- new `firefox` manifest entry: can be used to declare which versions of firefox are compatible with your theme: `version+` means "`version` and up", `up to version` means, well, you guessed it, `version1-version2` is an (inclusive) range and finally just `version` means exactly _that_ version. Use `major.minor` to specify the minor version, or omit it (or use ".x") to ignore it.
- prompt to show the manifest source of the theme you're installing. The prompt can be turned off with a new flag `--skip-manifest-source`
- command _reapply_ to reapply themes after firefox updates. The current themes for each profile are stored in a YAML file at `~/.config/ffcss/currently.yaml`. If ffcss tells you that a profile has no ffcss theme applied, try re-applying it manually with `ffcss use` so that this file gets created & filled.
- hooks to run custom shell commands before and after theme installations, via manifest entries `run.before` and `run.after`. Variants can override them.
- "did you mean ...?" message when a theme with a similar name is found
- command _get_ to download a theme without applying it

### Changed

- `ffcss use <theme name>` is now case-insensitive, punctuation-insensitive (`-`, `_` and `.`), whitespace-insensitive and unicode-insensitive (NFC normalization is applied before searching)

## [0.1.2] - 2021-07-10

### Added

- new theme in registry: [alpenblue](https://github.com/Godiesc/AlpenBlue) by [Godiesc](https://github.com/Godiesc)
- new theme in registry: [australis-tabs](https://github.com/sagars007/Australis-like-tabs-FF-ProtonUI-changes) by [sagars007](https://github.com/sagars007)
- new theme in registry: [blurredfox](https://github.com/manilarome/blurredfox) by [manilarome](https://github.com/manilarome)
- new theme in registry: [compactmode](https://github.com/Godiesc/compactmodefirefoxcss) by [Godiesc](https://github.com/Godiesc)
- new theme in registry: [diamondfucsia](https://github.com/Godiesc/DiamondFucsia) by [Godiesc](https://github.com/Godiesc)
- new theme in registry: [frozenfox](https://github.com/tortious/FrozenFox) by [tortious](https://github.com/tortious)
- new theme in registry: [halo](https://github.com/seirin-blu/Firefox-Halo) by [seirin-blu](https://github.com/seirin-blu)
- new theme in registry: [martinfox](https://github.com/arp242/MartinFox) by [arp242](https://github.com/arp242)
- new theme in registry: [montereyfox](https://github.com/FirefoxCSSThemers/Monterey-Fox) by [FirefoxCSSThemers](https://github.com/FirefoxCSSThemers)
- new theme in registry: [pro-fox](https://github.com/xmansyx/Pro-Fox) by [xmansyx](https://github.com/xmansyx)
- new theme in registry: [proton-connected-rounded-tabs](https://github.com/sagars007/Proton-UI-connected-rounded-tabs) by [sagars007](https://github.com/sagars007)
- new theme in registry: [technetium](https://github.com/edo0/Technetium) by [edo0](https://github.com/edo0)
- new theme in registry: [wavefox](https://github.com/QNetITQ/WaveFox) by [QNetITQ](https://github.com/QNetITQ)

### Fixed

- `ffcss init` was generating a file with a repository setting instead of download [[#36](https://github.com/ewen-lbh/ffcss/issues/36)]
- `ffcss init` added a .git suffix in the pre-filled value for `name`

## [0.1.1] - 2021-07-10

### Fixed

- crashes related to path handling on Windows

## [0.1.0] - 2021-07-09

### Added

#### For users
- Works on MacOS, GNU/Linux and Windows, tested on:
- Manjaro Linux Omara 21.0.7 (with kernel 5.12.9-1-MANJARO)
- MacOS Catalina 10.15.7
- Windows 10 20H2 (Build 19042.1083) (Please use the new Windows Terminal or something else that support ANSI escape sequences)
- a `use` commands to download & install themes
- works with any remote git repository
- works with any URL poiting to a .zip file
- shorthand for github repositories: OWNER/REPO
- a `init` command to add a `ffcss.yaml` manifest in your current repository
- basic for now, [a smarter version is planned](https://github.com/ewen-lbh/ffcss/issues/20)
- a `cache clear` command to clear the cache of downloaded repositories
- 8 themes available out-of-the-box (use them by typing their name only, it works)
- [chameleons-beauty](https://github.com/Godiesc/Chameleons-Beauty) by [Godiesc](https://github.com/Godiesc)
- [fxcompact](https://github.com/dannycolin/fx-compact-mode) by [dannycolin](https://github.com/dannycolin)
- [lepton](https://github.com/black7375/Firefox-UI-Fix) by [black7375](https://github.com/black7375)
- [materialfox](https://github.com/muckSponge/MaterialFox) by [muckSponge](https://github.com/muckSponge)
- [modoki](https://github.com/soup-bowl/Modoki-FirefoxCSS) by [soup-bowl](https://github.com/soup-bowl)
- [simplerentfox](https://github.com/MiguelRAvila/SimplerentFox) by [MiguelRAvila](https://github.com/MiguelRAvila)
- [verticaltabs](https://github.com/ranmaru22/firefox-vertical-tabs) by [ranmaru22](https://github.com/ranmaru22)
#### For theme makers
- a mechanism to handle theme variants:
- Variants can be declared in the same manifest file under the `variants` entry to override other entries
- per-OS paths: the {{ variant }} and {{ os }} placeholders get replaced with their values
- the value {{ os }} gets replaced with can be customized in the manifest file under the `os` entry, use `null` to mark an OS as incompatible
- Support for helper addons:
- Declare URLs to open after installation under the `addons` manifest entry (I plan to auto-install them in the future)
- Easy way to write about:config changes without writing a `user.js` file:
- Use the `config` manifest entry
- If you also use a `user.js`, you can combine both, they'll be written as a single `user.js` file
- Support for custom assets:
- Use the `assets` manifest entry to list out your assets
- Supports glob patterns
- If you store them under a `chrome` directory, you can use `copy from: chrome/` so that they don't get copied to `<profile directory>/chrome/chrome`

[Unreleased]: https://github.com/ewen-lbh/ffcss/compare/v0.1.2...HEAD
[0.1.2]: https://github.com/ewen-lbh/ffcss/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/ewen-lbh/ffcss/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/ewen-lbh/ffcss/releases/tag/v0.1.0

[//]: # (C3-2-DKAC:GGH:Rewen-lbh/ffcss:Tv{t})
