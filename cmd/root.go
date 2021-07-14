package cmd

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	. "helmup/pkg"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "helmup",
	Short: "Check for updates of your helm dependencies.",
	Long: `Helmup checks for updates of your helm dependencies
and lets you interactively choose which ones to apply in place.`,
	Example: "helmup /path/to/helm/directory",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(args); err != nil {
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

func run(args []string) error {
	path, err := GetProjectPath(args)
	if err != nil {
		return err
	}

	chartfile, err := LoadChartfile(path)
	if err != nil {
		return err
	}

	updatableDependencies, err := ResolveUpdates(chartfile)
	if err != nil {
		return err
	}
	if len(updatableDependencies) == 0 {
		fmt.Println(fmt.Sprintf("All dependencies are %s!", ansi.Color("up to date", "green")))
		return nil
	}

	chosenDependencies, err := Prompt(updatableDependencies)
	if err != nil {
		return err
	}
	if len(chosenDependencies) == 0 {
		return nil
	}

	err = UpdateChartfile(path, chartfile, chosenDependencies)
	if err != nil {
		return err
	}

	err = UpdateCharts(path)
	if err != nil {
		return err
	}

	return nil
}
