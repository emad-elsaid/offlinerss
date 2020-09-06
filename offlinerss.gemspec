# frozen_string_literal: true

Gem::Specification.new do |spec|
  spec.name          = 'offlinerss'
  spec.version       = '0.1.0'
  spec.authors       = ['Emad Elsaid']
  spec.email         = ['emad.elsaid.hamed@gmail.com']

  spec.summary       = 'Download RSS entries to local directory'
  spec.description   = 'Download RSS feed and split them to local directory'
  spec.homepage      = 'https://github.com/emad-elsaid/offlinerss'
  spec.license       = 'MIT'
  spec.required_ruby_version = Gem::Requirement.new('>= 2.3.0')

  spec.metadata['homepage_uri'] = spec.homepage
  spec.metadata['source_code_uri'] = 'https://github.com/emad-elsaid/offlinerss'

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`.split("\x0").reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  spec.bindir        = 'exe'
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }

  spec.add_runtime_dependency 'rss'
end
