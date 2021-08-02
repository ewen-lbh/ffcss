#!/usr/bin/env ruby
require "pathname"

name = ARGV[0]
files = Dir.glob "*.go"

files.each do |file|
	lines = Pathname.new(file).read.split"\n"
	lines.each do |line|
		signature_at = lines.find_index { |line| /^func\s+(\([^\)]+\)\s+)?\s*#{name}\s*\([^\)]*\).+\{$/ =~ line }
		if signature_at != nil
			puts "#{file}:#{signature_at+1}"
			exit 0
		end
	end
end

exit 1
