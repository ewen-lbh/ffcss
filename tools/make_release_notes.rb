#!/usr/bin/env ruby
$; = $, = "\n"
arg = ARGV[0]

major, minor, patch = `ffcss version`.split '.'

case arg
when "major"
	new_major, new_minor, new_patch = major.to_i + 1, 0, 0
when "minor"
	new_major, new_minor, new_patch = major, minor.to_i + 1, 0
when "patch"
	new_major, new_minor, new_patch = major, minor, patch.to_i + 1
else; if arg =~ /\d+.\d+.\d+/
	new_major, new_minor, new_patch = arg.split '.'
else
	puts "invalid argument #{arg}"
	exit 1
end
end


ffcss_dot_go = File.open "ffcss.go"

File.write "ffcss.go", (
	ffcss_dot_go.read
		.sub(/VersionMajor = \d+/, "VersionMajor = #{new_major}")
		.sub(/VersionMinor = \d+/, "VersionMinor = #{new_minor}")
		.sub(/VersionPatch = \d+/, "VersionPatch = #{new_patch}")
)

if not ["major", "minor", "patch"].include? arg
	exit 0
end

lines = File.open("CHANGELOG.md").read.split
release_notes = []


until lines.length == 0 or lines.shift =~ /^## \[#{new_major}.#{new_minor}.#{new_patch}\] - \d{4}-\d{2}-\d{2}$/ do
	# nothing
end

until lines.length == 0 or lines[0] =~ /^## \[\d+\.\d+\.\d+\] - \d{4}-\d{2}-\d{2}$/ do
	release_notes << lines.shift
end




File.write 'release_notes.md', release_notes.join
