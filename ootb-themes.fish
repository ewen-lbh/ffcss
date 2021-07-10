#!/usr/bin/env fish
# Out of the box themes
# generates markdown markup to be used in README.md/whatever.

for theme in themes/*.yaml
	set username (yq --raw-output .download < $theme | rev | cut -d/ -f2 | rev)
	echo - [(yq --raw-output .name < $theme)]\((yq --raw-output .download < $theme)\) by [$username]\(https://github.com/$username\)
end
