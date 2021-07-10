---
name: Adding a theme to the registry
about: To add a theme's manifest to the built-in ones, so that users can do `ffcss
  use YOUR-THEME`
title: Add <your theme's name> to the registry
labels: registry
assignees: ''

---

Steps:

1. Upload your manifest as `themes/<your theme's name>.yaml`
1. If you want your username to appear in the "Built-in themes" section of the README, make sure that you either
  - ~~Included a `by` field in the manifest~~ _Will be available once #33 is closed_
  - Used a github repository in your `download` field
  - Told me your author name/URL in this PR, so that I can add it manually when I release the next version of `ffcss`
1. Ran `chachacha added new theme <your theme's name> in the registry` (if you don't have `chachacha`, run `pip install chachacha` (you need Python for this) _Note: if you have problems installing chachacha, don't worry, I'll modify the changelog myself_
