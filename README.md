# Gitman
Gitman is a TUI tool for creating and managing git repositories.
_Check [gitman-lite](https://github.com/pyrod3v/gitman-lite) if you want a faster, more lighweight CLI-only option._

## Features
- Repository initialization
- .gitignore selection
- Custom gitignore templates

## Configuration
The application's configuration is located at `$USER/.config/gitman` on unix-like systems and at `%appdata%\Roaming\gitman` on Windows.
To add custom .gitignore templates, put any `<name>.gitignore` file in the gitignore directory in your config.

## Installing
To install the application, simply run `go install https://github.com/pyrod3v/gitman/cmd/gitman@latest` or clone this repository and run `go install`

## Contributing
All sorts of contributions are welcome. To, contribute:
1. Fork this repository
2. Create your feature branch
3. Commit and push your changes
4. Submit a pull request

Please use meaningful commit messages
