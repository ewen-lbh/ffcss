#!/usr/bin/env ruby
require "pathname"
require "tty-prompt"

name = ARGV[0]
from = Pathname.new ARGV[1]
to = Pathname.new ARGV[2]

body = []
doc = []
start_at = end_at = 0

lines = File.read(from).split "\n"

signature_at = lines.find_index { |line| /^func\s+(\([^\)]+\)\s+)?\s*#{name}\s*\([^\)]*\).+\{$/ =~ line }

if signature_at.nil?
	puts "signature for #{name} not found"
	exit 1
end

# find documentation
lines[...signature_at].reverse.each_with_index do |line, i|
  start_at = signature_at - (i + 1)
  if line.start_with? "//"
    doc.prepend line
    if line.start_with? "// #{name}"
      break
    end
  else
    break
  end
end

# go thru body
lines[signature_at..].each_with_index do |line, at|
  body << line
  end_at = signature_at + at
  # closing brace w/o indetation means end of func.
  # deal with it, just gofmt your code.
  if line == "}"
    break
  end
end

to.write (([""] + doc + body).join "\n"), to.size, mode: "a"
from.write (lines[...start_at] + lines[end_at + 1..]).join "\n"
`make format`
