package k8s

import (
	"github.com/spf13/cobra"
)

var (
	k8sShortDesc = "Provides k8s specific tooling capability like custom template generation"
	k8sLongDesc  = `Provides added capability for k8s related stuff.
		For example it can generate k8s templates etc.`
	k8sExample = `
	### Available commands for k8s
	gotils k8s (generate)`
)

func InitK8s() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "k8s [sub]",
		Short:            k8sShortDesc,
		Long:             k8sLongDesc,
		Example:          k8sExample,
		TraverseChildren: true,
	}

	cmd.AddCommand(NewK8sGenerator())
	return cmd
}
