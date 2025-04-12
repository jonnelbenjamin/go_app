package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	uppercase bool
	name      string
)

var rootCmd = &cobra.Command{
	Use:   "greeter",
	Short: "A colorful greeting app",
	Run: func(cmd *cobra.Command, args []string) {
		// If no name provided, prompt interactively
		if name == "" {
			color.Yellow("What's your name? ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
		}

		// Print greeting
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Hello, %s!\n", green(">>"), name)

		// Uppercase flag
		if uppercase {
			red := color.New(color.FgRed, color.Bold).SprintFunc()
			fmt.Printf("%s YOUR NAME IN UPPERCASE: %s\n", red("!!"), strings.ToUpper(name))
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&uppercase, "uppercase", "u", false, "Shout your name")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Provide name directly")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}