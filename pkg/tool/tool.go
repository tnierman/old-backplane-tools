package tool

import (
	"github.com/tnierman/backplane-tools/pkg/tool/oc"
	"github.com/tnierman/backplane-tools/pkg/tool/ocm"
)

type Tool interface {
	Name() string

	Install(to string) error

	Configure() error

	Remove() error
}

// Compile the list of available tools to manage
func List() []Tool {
	var tools []Tool
	tools = append(tools, ocm.NewTool())
	tools = append(tools, oc.NewTool())
	return tools
}
