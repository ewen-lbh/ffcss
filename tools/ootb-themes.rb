#!/usr/bin/env ruby
require "yaml"

major_version = `ffcss version major`.strip.to_i

Dir.glob('themes/*.yaml').each do |filename|
	theme = YAML.load_file filename, symbolize_names: true
	if theme[:ffcss] != major_version
		return
	end
	parts = theme[:download].split '/'
	username = parts.length >= 2 ? parts[-2] : nil
	puts "- [#{theme[:name]}](#{theme[:download]})" + (username ? " by [#{username}](https://github.com/#{username})" : "")
end
