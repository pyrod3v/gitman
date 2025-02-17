# Gitman
Gitman is a TUI tool for creating and managing git repositories.
_Check [gitman-lite](https://github.com/pyrod3v/gitman-lite) if you want a faster, more lighweight CLI-only option._

## Features
- Repository initialization
- .gitignore selection
- License selection
- Custom gitignore and license templates

## Installing
To install the application, simply run `go install https://github.com/pyrod3v/gitman/cmd/gitman@latest` or download a release from the [Releases Tab](https://github.com/pyrod3v/gitman/releases).

## Configuration
To add custom .gitignore templates, put any `<name>.gitignore` file in `USER/.gitman/gitignores/`.  
To add custom .gitignore templates, put any file in `USER/.gitman/licenses/`.  
The configuration file is located at `USER/.gitman/config.yaml`. Currently, you can set the following configuration keys:
- CacheGitignores: whether to cache fetched gitignores. Defaults to false.
- CacheLicenses: whether to cache fetched licenses. Defaults to false.

## Contributing
All sorts of contributions are welcome. To contribute:
1. Fork this repository
2. Create your feature branch
3. Commit and push your changes
4. Submit a pull request

Please use meaningful commit messages
