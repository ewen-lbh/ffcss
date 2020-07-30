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

Where:
    KEY         a setting key (see firefox's about:config)
    THEME_NAME  a theme name or URL
```
