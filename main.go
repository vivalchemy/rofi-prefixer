package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Cmd struct {
	Name      string
	Command   string
	Prefix    string
	Workspace int
}

func (c *Cmd) GenerateTerminalCommand() []string {
	return []string{"bash", "-c", c.Command}
}

func (c *Cmd) ExecutableTerminalCommand() *exec.Cmd {
	cmdArgs := c.GenerateTerminalCommand()
	return exec.Command(cmdArgs[0], cmdArgs[1:]...)
}

func (c *Cmd) SwitchWorkspace() error {
	if c.Workspace == 0 {
		return nil
	}
	cmd := exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", c.Workspace))
	return cmd.Run()
}

// ExecuteTerminalCommand executes the command and returns its output
func (c *Cmd) ExecuteTerminalCommand() (string, error) {
	// Execute the command
	cmd := c.ExecutableTerminalCommand()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Switch to the workspace if specified
	if err := c.SwitchWorkspace(); err != nil {
		return "", err
	}
	return string(output), nil
}

type Cmds []Cmd

func (commands *Cmds) ExecuteTerminalCommand(prefix string) {
	for _, cmd := range *commands {
		if cmd.Prefix == prefix {
			output, err := cmd.ExecuteTerminalCommand()
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", cmd.Name, err)
				return
			}
			fmt.Println("Output:", output)
			return
		}
	}
	fmt.Println("No matching command for prefix:", prefix)
}

func main() {
	ROFI_CMDS := Cmds{
		Cmd{
			Name:      `Calculator`,
			Command:   `rofi -show calc`,
			Prefix:    `=`,
			Workspace: 0,
		},
		Cmd{
			Name:      `window`,
			Command:   `rofi -show window`,
			Prefix:    `w`,
			Workspace: 0,
		},
		Cmd{
			Name:      `Applications`,
			Command:   `rofi -show drun`,
			Prefix:    `a`,
			Workspace: 0,
		},
		Cmd{
			Name:      `Browser`,
			Command:   `zen-browser $(echo "https://www.google.com/search?q=$(rofi -dmenu -p 'Enter Search Query:' | tr " " "+")")`,
			Prefix:    `g`,
			Workspace: 2,
		},
	}

	cmd := exec.Command("rofi", "-dmenu", "-p", "Prefix:")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	prefix := strings.TrimSpace(string(stdout))
	fmt.Printf("prefix: '%s'\n", prefix)

	ROFI_CMDS.ExecuteTerminalCommand(prefix)
}
