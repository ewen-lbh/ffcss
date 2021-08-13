#!/usr/bin/env ruby
require "pathname"
tested = []
untested = []

files = Dir.glob "*.go"
testfiles = Dir.glob "*_test.go"
function_pattern = /^func\s+(\(\s*\w+\s+(?<receiver>[^)\s]+)\s*\)\s+)?\s*(?<name>\w+).+\{$/

testfiles.each do |file|
  lines = Pathname.new(file).read.split "\n"
  lines.each do |line|
    match = line.match function_pattern
    if match != nil and match[:name].start_with? "Test"
      tested << match[:name].delete_prefix("Test")
    end
  end
end

files.each do |file|
  lines = Pathname.new(file).read.split "\n"
  lines.each_with_index do |line, i|
    match = line.match function_pattern
    if match != nil and not match[:name].start_with? "Test" and match[:name] !~ /main|init/
      if not tested.find { |func| func.downcase == match[:name].downcase } and not file.end_with? "_test.go"
        untested << {
          name: (if match[:receiver] then "(#{match[:receiver]})." else "" end) + match[:name],
          # name: match[:name],
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
  puts (align untested, func, :location) + " "  + (align untested, func, :name)
end

puts
puts (if untested.size > 1
  "#{untested.size} functions"
else
  "one function"
end) + ", " + (Float(untested.size) / (untested.size + tested.size) * 100).round.to_s + "%"

exit 1
