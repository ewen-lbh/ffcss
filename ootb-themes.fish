#!/usr/bin/env fish
# Out of the box themes
# generates markdown markup to be used in README.md/whatever.

for theme in themes/*.yaml
	echo - [(yq --raw-output .name < $theme)]\((yq --raw-output .download < $theme)\)
end
