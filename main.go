package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Cmd represents a command to be executed, with options for workspace and prefix
type Cmd struct {
	Browser   bool   // Browser to use for opening links
	Command   string // Actual command to execute
	Name      string // Friendly name of the command
	Prefix    string // Prefix to trigger the command
	Workspace int    // Workspace to switch to after executing the command
}

type Cmds struct {
	commands map[string]*Cmd
}

// NewCmds initializes a new Cmds instance from a list of Cmd objects.
func newCmds(cmdList []Cmd) *Cmds {
	cmds := &Cmds{
		commands: make(map[string]*Cmd),
	}
	for i := range cmdList {
		cmds.commands[cmdList[i].Prefix] = &cmdList[i]
	}
	return cmds
}

// Main logic of the program
func (cmds *Cmds) findCommand(input string) (cmd *Cmd, query string) {
	// Split the prefix and query
	prefix, query, _ := strings.Cut(input, " ")

	cmd, exists := cmds.commands[prefix]
	if !exists {
		// fmt.Println("No matching command for prefix:", prefix)
		return nil, ""
	}
	// fmt.Println("fC cmd:", cmd, "query:", query)
	return cmd, query
}

func (cmd *Cmd) needsQuery() bool {
	// First, check if the string contains %s
	if !strings.Contains(cmd.Command, "%s") {
		return false
	}

	// Use a regular expression to check for non-escaped %s
	re := regexp.MustCompile(`(^|[^\\])%s`)
	return re.MatchString(cmd.Command)
}

func getQuery(input string) string {
	if input != "" {
		return input
	}
	cmd := exec.Command("bash", "-c", "rofi -dmenu -p 'Enter Search Query:'")
	query, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run rofi command: %v\n", err)
		os.Exit(1)
	}
	// fmt.Println("gQ query:", string(query))
	return strings.TrimSpace(string(query))
}

func (cmd *Cmd) convertToBrowserQuery(query string) string {
	if cmd.Browser {
		query = strings.ReplaceAll(query, " ", "+")
	}
	// fmt.Println("cTBQ query:", query)
	return query
}

func (cmd *Cmd) makeBrowserCommand(query string) string {
	if !cmd.Browser {
		return query
	}
	// fmt.Println("mBC query: zen-browser", query)
	return fmt.Sprintf("zen-browser %s", query)
}

func (cmd *Cmd) replacePlaceholderInCommand(query string) string {
	// Temporarily replace escaped \%s to handle literals
	const placeholder = "__ESCAPED_PERCENT_S__"
	command := strings.ReplaceAll(cmd.Command, `\%s`, placeholder)

	// Replace %s with the query
	command = strings.ReplaceAll(command, "%s", query)

	// Restore escaped %s
	command = strings.ReplaceAll(command, placeholder, `%s`)

	// fmt.Println("rPC command:", command)
	return command
}

func (cmd *Cmd) executeCommand(command string) string {
	executableCommand := exec.Command("bash", "-c", command)
	output, err := executableCommand.Output()
	if err != nil {
		// fmt.Println("Error executing command:", err)
		return ""
	}
	// fmt.Println("eC output:", string(output))
	return string(output)
}

func (c *Cmd) switchWorkspace() error {
	if c.Workspace == 0 {
		return nil
	}
	cmd := exec.Command("hyprctl", "dispatch", "workspace", fmt.Sprintf("%d", c.Workspace))
	return cmd.Run()
}

// Main logic to display the Rofi menu and execute selected command.
func main() {
	rofiCmds := newCmds([]Cmd{
		{
			Command:   `rofi -show drun`,
			Name:      `Applications`,
			Browser:   false,
			Prefix:    `a`,
			Workspace: 0,
		},
		{
			Browser:   true,
			Command:   "https://www.google.com/search?q=%s",
			Name:      `Google`,
			Prefix:    `g`,
			Workspace: 2,
		},
		{
			Browser:   false,
			Command:   `rofi -show calc`,
			Name:      `Calculator`,
			Prefix:    `=`,
			Workspace: 0,
		},
		{
			Browser:   true,
			Command:   "https://chat.openai.com/?q=%s",
			Name:      `Chatgpt.com`,
			Prefix:    `gpt`,
			Workspace: 2,
		},
		{
			Browser:   true,
			Command:   "https://claude.ai/new/?q=%s",
			Name:      `Claude ai`,
			Prefix:    `claude`,
			Workspace: 2,
		},
		{
			Name:      `Perplexity.ai`,
			Command:   "https://www.perplexity.ai/search?q=%s",
			Browser:   true,
			Prefix:    `ai`,
			Workspace: 2,
		},
		{
			Name:      `window`,
			Command:   `rofi -show window`,
			Browser:   false,
			Prefix:    `w`,
			Workspace: 0,
		},
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

	// Command flow
	rawPrefix := strings.TrimSpace(string(stdout))
	prefix, _, _ := strings.Cut(rawPrefix, "-->")
	prefix = strings.TrimSpace(prefix)
	fmt.Printf("prefix: '%s'\n", prefix)
	command, query := rofiCmds.findCommand(prefix)
	if command.needsQuery() {
		query = getQuery(query)
	}
	query = command.convertToBrowserQuery(query)
	query = command.replacePlaceholderInCommand(query)
	query = command.makeBrowserCommand(query)
	finalOutput := command.executeCommand(query)
	if err := command.switchWorkspace(); err != nil {
		fmt.Printf("Error switching workspace: %v\n", err)
		return
	}
	fmt.Println(finalOutput)
}
