#!/usr/bin/env ruby
require "yaml"

Dir.glob('themes/*.yaml').each do |filename|
	theme = YAML.load_file filename, symbolize_names: true
	parts = theme[:download].split '/'
	username = parts.length >= 2 ? parts[-2] : nil
	puts "- [#{theme[:name]}](#{theme[:download]})" + (username ? " by [#{username}](https://github.com/#{username})" : "")
end
