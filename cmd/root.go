package cmd

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helmup/pkg"
	"os"
	"runtime"
)

var Version = "0.1.0-SNAPSHOT"
var GitCommit = "manual"

var rootCmd = &cobra.Command{
	Use:   "helmup [flags] [path]",
	Short: "Check for updates of your helm dependencies.",
	Long: `helmup checks for updates of your helm dependencies
and lets you interactively choose which ones to apply in place.`,
	Example: "helmup /path/to/helm/directory",
	Version: fmt.Sprintf("%s (%s)\n%s", Version, GitCommit, runtime.Version()),
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		notInteractive, err := cmd.Flags().GetBool("no-interactive")
		if err != nil {
			cobra.CheckErr(err)
		}

		path, err := pkg.GetProjectPath(args)
		if err != nil {
			cobra.CheckErr(err)
		}

		if err := run(path, !notInteractive); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("no-interactive", "n", false, "only print updates")
	rootCmd.AddCommand(completionCmd)
}

func run(path string, interactive bool) error {
	helmCli := cli.New()
	settings := &pkg.ResolverSettings{Env: helmCli, Getters: getter.All(helmCli)}

	chartfile, err := pkg.LoadChartfile(path)
	if err != nil {
		return err
	}

	updatableDependencies, err := pkg.ResolveUpdates(chartfile, settings)
	if err != nil {
		return err
	}
	if len(updatableDependencies) == 0 {
		fmt.Println(fmt.Sprintf("All dependencies are %s!", ansi.Color("up to date", "green")))
		return nil
	}

	if !interactive {
		for _, dependency := range updatableDependencies {
			fmt.Println(dependency.String())
		}

		return nil
	}

	chosenDependencies, err := pkg.Prompt(updatableDependencies)
	if err != nil {
		return err
	}
	if len(chosenDependencies) == 0 {
		return nil
	}

	err = pkg.UpdateChartfile(path, chartfile, chosenDependencies)
	if err != nil {
		return err
	}

	err = pkg.UpdateCharts(path, settings)
	if err != nil {
		return err
	}

	return nil
}
