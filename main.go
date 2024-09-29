package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Cmd represents a command to be executed, with options for workspace and prefix
type Cmd struct {
	Name      string // Friendly name of the command
	Command   string // Actual command to execute
	Prefix    string // Prefix to trigger the command
	Workspace int    // Workspace to switch to after executing the command
}

// CommandGenerator generates the final command to be executed.
// If `takeRofiInput` is true, the command uses rofi for user input; otherwise, it uses the provided `query`.
func (c *Cmd) CommandGenerator(takeRofiInput bool, query string) string {
	// Temporarily replace escaped \%s to handle literals
	const placeholder = "__ESCAPED_PERCENT_S__"
	command := strings.ReplaceAll(c.Command, `\%s`, placeholder)

	// Replace %s depending on whether rofi input is needed
	if takeRofiInput {
		command = strings.ReplaceAll(command, "%s", `rofi -dmenu -p 'Enter Search Query:'`)
	} else {
		command = strings.ReplaceAll(command, "%s", fmt.Sprintf(`echo "%s"`, query))
	}

	// Restore escaped %s
	command = strings.ReplaceAll(command, placeholder, `%s`)
	return command
}

// GenerateTerminalCommand returns the full terminal command for execution.
func (c *Cmd) GenerateTerminalCommand(takeRofiInput bool, query string) []string {
	return []string{"bash", "-c", c.CommandGenerator(takeRofiInput, query)}
}

// ExecutableTerminalCommand prepares the exec.Cmd structure for execution.
func (c *Cmd) ExecutableTerminalCommand(takeRofiInput bool, query string) *exec.Cmd {
	cmdArgs := c.GenerateTerminalCommand(takeRofiInput, query)
	return exec.Command(cmdArgs[0], cmdArgs[1:]...)
}

func (c *Cmd) SwitchWorkspace() error {
	if c.Workspace == 0 {
		return nil
	}
	cmd := exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", c.Workspace))
	return cmd.Run()
}

// ExecuteTerminalCommand executes the command and switches to the workspace if specified.
func (c *Cmd) ExecuteTerminalCommand(takeRofiInput bool, query string) (string, error) {
	// Execute the command
	cmd := c.ExecutableTerminalCommand(takeRofiInput, query)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	// Switch workspace if specified
	if err := c.SwitchWorkspace(); err != nil {
		return "", fmt.Errorf("workspace switch failed: %w", err)
	}

	return string(output), nil
}

// Cmds manages a set of commands that can be executed by prefix.
type Cmds struct {
	commands map[string]*Cmd
}

// NewCmds initializes a new Cmds instance from a list of Cmd objects.
func NewCmds(cmdList []Cmd) *Cmds {
	cmds := &Cmds{
		commands: make(map[string]*Cmd),
	}
	for i := range cmdList {
		cmds.commands[cmdList[i].Prefix] = &cmdList[i]
	}
	return cmds
}

// FindAndExecuteTerminalCommand finds a command by prefix and executes it.
func (cmds *Cmds) FindAndExecuteTerminalCommand(input string) {
	// Split the prefix and query
	prefix, query, found := strings.Cut(input, " ")

	cmd, exists := cmds.commands[prefix]
	if !exists {
		fmt.Println("No matching command for prefix:", prefix)
		return
	}

	// Execute the command; if no query is found, rofi input will be used.
	_, err := cmd.ExecuteTerminalCommand(!found, query)
	if err != nil {
		fmt.Printf("Error executing command '%s': %v\n", cmd.Name, err)
		return
	}
	// fmt.Println("Output:", output)
}

// Main logic to display the Rofi menu and execute selected command.
func main() {
	rofiCmds := NewCmds([]Cmd{
		{
			Name:      `Applications`,
			Command:   `rofi -show drun`,
			Prefix:    `a`,
			Workspace: 0,
		},
		{
			Name:      `Google`,
			Command:   `zen-browser "https://www.google.com/search?q=$(%s | tr " " "+")"`,
			Prefix:    `g`,
			Workspace: 2,
		},
		{
			Name:      `Calculator`,
			Command:   `rofi -show calc`,
			Prefix:    `=`,
			Workspace: 0,
		},
		{
			Name:      `Chatgpt.com`,
			Command:   `zen-browser "https://chat.openai.com/?q=$(%s | tr " " "+")"`,
			Prefix:    `gpt`,
			Workspace: 2,
		},
		{
			Name:      `Claude ai`,
			Command:   `zen-browser "https://claude.ai/new/?q=$(%s | tr " " "+")"`,
			Prefix:    `claude`,
			Workspace: 2,
		},
		{
			Name:      `Perplexity.ai`,
			Command:   `zen-browser "https://www.perplexity.ai/search?q=$(%s | tr " " "+")"`,
			Prefix:    `ai`,
			Workspace: 2,
		},
		{
			Name:      `window`,
			Command:   `rofi -show window`,
			Prefix:    `w`,
			Workspace: 0,
		},
		// Commented out as in the original
		// {
		// 	Name:      `Googel Gemini Ai`,
		// 	Command:   `zen-browser "https://gemini.google.com/app?q=$(%s | tr " " "+")"`,
		// 	Prefix:    `gem`,
		// 	Workspace: 2,
		// },
	})

	var rofiMenu strings.Builder
	for prefix, cmd := range rofiCmds.commands {
		rofiMenu.WriteString(fmt.Sprintf("%s --> %s\n", prefix, cmd.Name))
	}
	rofiPrompt := strings.TrimSuffix(rofiMenu.String(), "\n")

	// Display Rofi prompt for user to select a command
	rofiCommand := fmt.Sprintf(`echo "%s" | sort | rofi -sep "\n" -dmenu -p 'Prefix:' -i -mesg 'Select a command'`, rofiPrompt)
	cmd := exec.Command("bash", "-c", rofiCommand)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run rofi command: %v\n", err)
		os.Exit(1)
	}

	rawPrefix := strings.TrimSpace(string(stdout))
	prefix, _, _ := strings.Cut(rawPrefix, "-->")
	prefix = strings.TrimSpace(prefix)
	// fmt.Printf("prefix: '%s'\n", prefix)
	rofiCmds.FindAndExecuteTerminalCommand(prefix)
	// fmt.Println(rofiPrompt)
}
