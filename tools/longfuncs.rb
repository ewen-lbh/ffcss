#!/usr/bin/env ruby
require "pathname"

THRESHOLD = ARGV[0].to_i

files = Dir.glob "*.go"
funcs = {}
locations = {}
current_func = ""

files.each do |file|
  lines = Pathname.new(file).read.split "\n"
  lines.each_with_index do |line, at|
    if current_func == ""
      match = line.match /^func\s+(\([^\)]+\)\s+)?\s*(?<name>\w+)\s*\([^\)]*\).+\{$/
      if match != nil
        current_func = match[:name]
        funcs[current_func] = []
        locations[current_func] = "#{file}:#{at + 1}"
      end
    else
      funcs[current_func] << line
      current_func = "" if line == "}"
    end
  end
end

funcs.filter! do |name, body|
  body.size > THRESHOLD and not name.start_with? "Test"
end

exit 0 if funcs.size == 0

funcs = funcs.sort_by { |name, body| body.size }

funcs.each do |func, body|
  puts "#{locations[func]}\t\t\t#{body.size}\t#{func}"
end

exit 1
