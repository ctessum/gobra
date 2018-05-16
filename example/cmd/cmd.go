package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

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

	// addition flags
	num1, num2 int

	// paths of file to print, measure
	path  string
	path2 string
	path3 []string
)

func init() {
	// Link the commands together.
	Root.AddCommand(versionCmd)
	Root.AddCommand(runCmd)
	runCmd.AddCommand(steadyCmd)
	runCmd.AddCommand(addition)
	runCmd.AddCommand(printCmd)
	runCmd.AddCommand(printMultipleCmd)

	// Create the configuration flags.
	Root.PersistentFlags().StringVar(&configFile, "config", "./conf.toml", "configuration file location")

	runCmd.PersistentFlags().BoolVarP(&inBackground, "inBackground", "s", false, "Program will run in background if sent true")
	steadyCmd.Flags().IntSliceVar(&layers, "layers", []int{0, 2, 4, 6}, "Dummy slice of ints")
	steadyCmd.Flags().IntVar(&begin, "begin", 0, "Beginning row index.")
	addition.Flags().IntVar(&num1, "num1", 1, "First number")
	addition.Flags().IntVar(&num2, "num2", 1, "Second number")
	printCmd.Flags().StringVar(&path, "path", "", "filepath to determine length")
	printCmd.Flags().StringVar(&path2, "path2", "", "file to print")
	printMultipleCmd.Flags().StringSliceVar(&path3, "path3", []string{""}, "files to print")

}

// Root is the main command.
var Root = &cobra.Command{
	Use:   "dummy",
	Short: "A dummy program",
	Long:  `This is a longer description for a dummy program, which does not do anything and only exists for the purpose of being an example.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("I'm always printed, because I'm in PersistentPreRun. (server)")
		cmd.Println("I'm always printed, because I'm in PersistentPreRun. (client)")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hey, run me before Run execute")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I'm dummy")
		cmd.Println("This is supposed to print to the other output")
	},
	DisableAutoGenTag: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "version prints the version number of this version of dummy.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I'm version!")
		cmd.Println("Version 1")
	},
	DisableAutoGenTag: true,
}

var runCmd = &cobra.Command{
	Use:               "run",
	Short:             "Run the program.",
	Long:              `run runs program and executes it. If there's no subcommand, I'll count from 0 to 9.`,
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Running program the program.")
		for i := range make([]int, 10) {
			cmd.Printf("Output: %d\n", i)
			time.Sleep(time.Duration(200) * time.Millisecond)
		}
		return nil
	},
}

var addition = &cobra.Command{
	Use:   "add",
	Short: "adds two number",
	Long:  "We perform the addition operation on two numerical operands. The operation yields the sum of two operands, which are inputted as flags.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(num1 + num2)
	},
}

// steadyCmd is a command that runs a steady-state simulation.
var steadyCmd = &cobra.Command{
	Use:   "steady",
	Short: "Run dummy in steady-state mode.",
	Long:  `steady runs program in steady-stated mode so nothing goes out of stability.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Yeah, I'm running steady")
		return errors.New("Oh no! An error! I'm not steady")
	},
	DisableAutoGenTag: true,
}

// printCmd takes a file path and prints it
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints file content",
	Long:  "Takes in a filepath and prints its content. Easy enough.",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(path)
		c, err := ioutil.ReadFile(path2)
		if err != nil {
			cmd.Println(err)
		}

		cmd.Println("Filesize of the first file: ", len(string(b)))
		cmd.Println("Content of ", path2, "is: ", string(c))
	},
}

// printCmd takes multiple files and prints them.
var printMultipleCmd = &cobra.Command{
	Use:   "multiprint",
	Short: "prints content of multiple files",
	Long:  "Takes in files and prints their contents.",
	Run: func(cmd *cobra.Command, args []string) {
		for i, path := range path3 {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				cmd.Println(err)
			}
			cmd.Printf("Content of file %d (%s) is:\n\n%s\n\n", i, path, string(data))
		}
	},
}
