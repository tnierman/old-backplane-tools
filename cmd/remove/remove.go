package remove

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tnierman/backplane-tools/pkg/tool"
	"github.com/tnierman/backplane-tools/pkg/utils"
)

func Cmd() *cobra.Command {
	removeCmd := &cobra.Command {
		Use: fmt.Sprintf("remove [all|%s]", strings.Join(tool.Names, "|")),
		Args: cobra.OnlyValidArgs,
		ValidArgs: append(tool.Names, "all"),
		Short: "Remove a tool",
		Long: "Removes one or more tools from the given list. It's valid to specify multiple tools: in this case, all tools provided will be removed. If 'all' is explicitly passed, then the entire tool directory will be removed, providing a clean slate for reinstall. If no specific tools are provided, no action is taken",
		RunE: run,
	}
	return removeCmd
}

// run removes the tool(s) specified by the provided positional args
func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		fmt.Println("No tools specified to be removed. In order to remove all tools, explicitly specify 'all'")
		return nil
	}
	if utils.Contains(args, "all") {
		return tool.RemoveInstallDir()
	}

	fmt.Println("Removing the following tools:")
	removeList := []tool.Tool{}
	for _, toolName := range args {
		fmt.Printf("- %s\n", toolName)
		removeList = append(removeList, tool.Map[toolName])
	}

	err := tool.Remove(removeList)
	if err != nil {
		return fmt.Errorf("failed to remove one or more tools: %w", err)
	}
	return nil
}
