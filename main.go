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

type Cmds struct {
	commands map[string]*Cmd
}

func NewCmds(cmdList []Cmd) *Cmds {
	cmds := &Cmds{
		commands: make(map[string]*Cmd),
	}
	for i := range cmdList {
		cmds.commands[cmdList[i].Prefix] = &cmdList[i]
	}
	return cmds
}

func (commands *Cmds) ExecuteTerminalCommand(prefix string) {
	cmd, exists := commands.commands[prefix]
	if !exists {
		fmt.Println("No matching command for prefix:", prefix)
		return
	}

	output, err := cmd.ExecuteTerminalCommand()
	if err != nil {
		fmt.Printf("Error executing command '%s': %v\n", cmd.Name, err)
		return
	}
	fmt.Println("Output:", output)
}

func main() {
	ROFI_CMDS := NewCmds([]Cmd{
		{
			Name:      `Applications`,
			Command:   `rofi -show drun`,
			Prefix:    `a`,
			Workspace: 0,
		},
		{
			Name:      `Browser`,
			Command:   `zen-browser $(echo "https://www.google.com/search?q=$(rofi -dmenu -p 'Enter Search Query:' | tr " " "+")")`,
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
			Command:   `zen-browser $(echo "https://chat.openai.com/?q=$(rofi -dmenu -p 'Enter Search Query:' | tr " " "+")")`,
			Prefix:    `gpt`,
			Workspace: 2,
		},
		{
			Name:      `Claude ai`,
			Command:   `zen-browser $(echo "https://claude.ai/new/?q=$(rofi -dmenu -p 'Enter Search Query:' | tr " " "+")")`,
			Prefix:    `claude`,
			Workspace: 2,
		},
		{
			Name:      `Perplexity.ai`,
			Command:   `zen-browser $(echo "https://www.perplexity.ai/search?q=$(rofi -dmenu -p 'Enter Search Query:' | tr " " "+")")`,
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
		// 	Command:   `zen-browser $(echo "https://gemini.google.com/app?q=$(rofi -dmenu -p 'Enter Search Query:' | tr " " "+")")`,
		// 	Prefix:    `gem`,
		// 	Workspace: 2,
		// },
	})
	cmd := exec.Command("rofi", "-dmenu", "-p", "Prefix:")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	prefix := strings.TrimSpace(string(stdout))
	fmt.Printf("prefix: '%s'\n", prefix)

	ROFI_CMDS.ExecuteTerminalCommand(prefix)
}
