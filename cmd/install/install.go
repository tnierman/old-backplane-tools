package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tnierman/backplane-tools/pkg/tool"
	"github.com/tnierman/backplane-tools/pkg/utils"
)

var (
	toolMap = map[string]tool.Tool{}
	toolNames = []string{}
)

func init() {
	for _, tool := range tool.List() {
		name := tool.Name()
		toolMap[name] = tool
		toolNames = append(toolNames, name)
	}
}

func Cmd() *cobra.Command {
	installCmd := &cobra.Command{
		Use: fmt.Sprintf("install [all|%s]", strings.Join(toolNames, "|")),
		Args: cobra.OnlyValidArgs,
		ValidArgs: append(toolNames, "all"),
		Short: "Install or upgrade a new tool",
		RunE: run,
	}
	return installCmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || utils.Contains(args, "all") {
		// If user doesn't specify, or explicitly passes 'all', give them all the things
		args = toolNames
	}

	fmt.Printf("Installing the following tools: %s\n", strings.Join(args, ", "))
	installList := []tool.Tool{}
	for _, specifiedTool := range args {
		installList =  append(installList, toolMap[specifiedTool])
	}

	err := installTools(installList)
	if err != nil {
		return fmt.Errorf("failed to install tools: %w", err)
	}
	return nil
}

func installTools(tools []tool.Tool) error {
	dir, err := createInstallDir()
	if err != nil {
		return fmt.Errorf("failed to create installation directory: %w", err)
	}

	for _, tool := range tools {
		fmt.Println("")
		fmt.Printf("Installing %s\n", tool.Name())
		err = tool.Install(dir)
		if err != nil {
			fmt.Printf("Encountered error while installing %s: %v\n", tool.Name(), err)
			fmt.Println("Skipping...")
		} else {
			fmt.Printf("Successfully installed %s\n", tool.Name())
		}
	}
	return nil
}

func createInstallDir() (installDir string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve $HOME dir: %w", err)
	}
	installDir = filepath.Join(homeDir, ".local", "bin", "backplane")
	err = os.MkdirAll(installDir, os.FileMode(0755))
	return installDir, err
}
