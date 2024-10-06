# Rofi Prefixer

Rofi Prefixer is a Go-based utility that enhances the functionality of Rofi by adding customizable prefixes for quick actions and searches.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Features

- Customizable command prefixes
- Integration with various search engines and websites
- Workspace switching capability

## Installation

### Prerequisites

- Go 1.23.1 or higher
- Rofi
- Hyprland (for workspace switching)

### Building from source

1. Clone the repository:
   ```bash
   git clone https://github.com/vivalchemy/rofi-prefixer.git
   cd rofi-prefixer
   ```

2. Build the project:
   ```bash
   make build
   # if not using make, you can also run `go build -o ./build/rofi-prefixer ./main.go`
   ```

3. The binary will be available in the `./build/` directory.

## Usage

Run the `rofi-prefixer` binary:

```bash
./build/rofi-prefixer
```

This will open a Rofi menu with all available prefixes. Select a prefix to execute the associated command.

## Configuration

Edit the `ROFI_CMDS` slice in `main.go` to customize your prefixes and commands. Each command is defined as a `Cmd` struct with the following fields:

- `Command`: The command to execute
- `Name`: A friendly name for the command
- `Browser`: Set to `true` if the command should be opened in a browser
- `Prefix`: The prefix to trigger this command
- `Workspace`: The workspace number to switch to after executing the command (0 for no switch)

Example:

```go
{
    Browser:   true,
    Command:   "https://www.google.com/search?q=%s",
    Name:      `Google`,
    Prefix:    `g`,
    Workspace: 2,
},
```

> [!NOTE]
> The `Command` field supports placeholders. You can use `%s` to represent the search query. For example, the `Command` field for the Google search could be `https://www.google.com/search?q=%s`. To use a literal `%s` in the `Command` field, you can escape it with a backslash `\%s`. Using %s in the Command field prompts for input if no query is provided after the prefix.

## Development

### Prerequisites

- Go 1.23.1 or higher
- Make (optional)

### Available Make commands

- `make help`: Print available commands
- `make build`: Build the application
- `make run`: Run the application
- `make run/live`: Run the application with live reloading
- `make test`: Run all tests
- `make audit`: Run quality control checks
- `make tidy`: Tidy modfiles and format Go files

For a complete list of available commands, run `make help`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
