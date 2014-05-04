// TestFlag project main.go
package main

import (
	// local
	"gnvm/command"
)

/*
func cobraMode() {

	var (
		Global bool
	)

	// defind root cmd
	var GnvmCmd = &cobra.Command{
		Use:   "gnvm",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
	            love by spf13 and friends in Go.
	            Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	// sub cmd
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}

	var installCmd = &cobra.Command{
		Use:   "install",
		Short: " install latest node.js version",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Global)
			fmt.Println(args)
			fmt.Println(" install latest node.js version")
		},
	}

	// add sub cmd to root
	GnvmCmd.AddCommand(versionCmd)
	GnvmCmd.AddCommand(installCmd)

	// set flag to root
	GnvmCmd.PersistentFlags().BoolVarP(&Global, "global", "g", false, "golbal output")

	// exec
	GnvmCmd.Execute()
}
*/

func main() {
	command.Exec()
}
