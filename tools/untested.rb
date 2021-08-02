#!/usr/bin/env ruby
require "pathname"
tested = []
untested = []

files = Dir.glob "*.go"
testfiles = Dir.glob "*_test.go"

testfiles.each do |file|
  lines = Pathname.new(file).read.split "\n"
  lines.each do |line|
    match = line.match /^func\s+(\([^\)]+\)\s+)?\s*(?<name>\w+)\s*\([^\)]*\).+\{$/
    if match != nil and match[:name].start_with? "Test"
      tested << match[:name].delete_prefix("Test")
    end
  end
end

files.each do |file|
  lines = Pathname.new(file).read.split "\n"
  lines.each_with_index do |line, i|
    match = line.match /^func\s+(\s*(?<receiver>\w+)\s+.+\)\s+)?\s*(?<name>\w+)\s*\([^\)]*\).+\{$/
    if match != nil and not match[:name].start_with? "Test" and match[:name] !~ /main|init/
      if not tested.find { |func| func.downcase == match[:name].downcase } and not file.end_with? "_test.go"
        untested << {
          name: (match[:receiver] || "") + match[:name],
          location: "#{file}:#{i + 1}",
        }
      end
    end
  end
end

def align(untested, func, key, right: false)
  longest_length = untested.map { |f| f[key].length }.max
  if right
    func[key].rjust longest_length
  else
    func[key].ljust longest_length
  end
end

if untested.size == 0
  puts "everything's tested"
  exit 0
end

untested.each do |func|
  puts "untested: #{align untested, func, :name, right: true} @ #{align untested, func, :location}"
end

puts "#{untested.size} function(s)"

exit 1
