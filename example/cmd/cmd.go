/**
	This file is MIT licensed.
**/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// These variables specify confuration flags.
var (
	// configFile specifies the location of the configuration file.
	configFile string

	// should this dummy program runs in background
	inBackground bool

	// dummy slice
	layers []int

	// dummy starting index
	begin int
)

func init() {
	// Link the commands together.
	Root.AddCommand(versionCmd)
	Root.AddCommand(runCmd)
	runCmd.AddCommand(steadyCmd)

	// Create the configuration flags.
	Root.PersistentFlags().StringVar(&configFile, "config", "./conf.toml", "configuration file location")

	runCmd.PersistentFlags().BoolVarP(&inBackground, "inBackground", "s", false, "Program will run in background if sent true")
	steadyCmd.Flags().IntSliceVar(&layers, "layers", []int{0, 2, 4, 6},	"Dummy slice of ints")
	steadyCmd.Flags().IntVar(&begin, "begin", 0, "Beginning row index.")

}

// Root is the main command.
var Root = &cobra.Command{
	Use:   "dummy",
	Short: "A dummy program",
	Long: `This is a longer description for a dummy program, which does not do anything and only exists for the purpose of being an example.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(`I'm ran.`)
	},
	DisableAutoGenTag: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "version prints the version number of this version of dummy.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 1")
	},
	DisableAutoGenTag: true,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the program.",
	Long: `run runs program and executes it.`,
	DisableAutoGenTag: true,
}

// steadyCmd is a command that runs a steady-state simulation.
var steadyCmd = &cobra.Command{
	Use:   "steady",
	Short: "Run dummy in steady-state mode.",
	Long: `steady runs program in steady-stated mode so nothing goes out of stability.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Yeah, I'm running steady")
		return nil
	},
	DisableAutoGenTag: true,
}