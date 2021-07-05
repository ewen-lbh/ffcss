1. [`utils.go:ProfileDirsPaths`] Get profile directories' paths
2. [`download.go`] Download the data into the config
   1. [`ResolveURL`] Resolve the given theme name into a URL to download the data from/a repo to clone into (and also determine whether to clone or copy). If the passed theme is not a valid URL or github shorthand, try `.config/ffcss/themes/<themeName>.yaml:repository`.
   2. [`Download`] Download or clone to `.cache/ffcss/<themeName>`
3. [`manifest.go`] Search for a manifest
    1. [`LoadManifest`] Try locally at `.config/ffcss/themes/<themeName>.yaml`
    2. [`LoadManifest`] Use `.cache/ffcss/<themeName>/ffcss.yaml`
1. [`use.go`] Let user choose a variant, if needed
    1. [`AskForVariant`] prompt user
2. [`utils.go:GOOStoOS`] Resolve which OS to use
2. [`manifest.go`] Select which fields from manifest to use
    1. [`(Manifest).Resolve`] Resolve a `Manifest` to a `Theme`
3. [`userprefs.go`] Generate a `user.js` from `config`
4. [`copy.go`] Resolve manifest into "copy over" instructions
    1. Concatenate userChrome, userContent, user.js and assets
    1. [`RenderFileTemplate`] Replace `{{ os }}` and `{{ variant }}`
    1. [`os.Glob`] Glob every item of the list
    1. [`DestinationPathOf`] Take their path relative to `.cache/ffcss/<themeName>` (not including `<themeName>`), and prepend firefox session directory to the relative path
