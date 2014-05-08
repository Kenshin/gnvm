package command

import (
	// lib
	"github.com/spf13/cobra"

	// go
	"fmt"
	"strconv"
	"strings"

	// local
	"gnvm/config"
	"gnvm/nodehandle"
)

var (
	global bool
	remote bool
)

// defind root cmd
var gnvmCmd = &cobra.Command{
	Use:   "gnvm",
	Short: "Gnvm is Node.js version manage by GO on Win",
	Long: `Gnvm is Node.js version manage by GO on Win
           like nvm. Complete documentation is available at https://github.com/kenshin/gnvm`,
	Run: func(cmd *cobra.Command, args []string) {
		// TO DO
	},
}

// sub cmd
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Gnvm",
	Long:  `Print the version number of Gnvm`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v" + config.VERSION)
	},
}

// sub cmd
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install any node.js version",
	Long: `install any node.js version like :
'gnvm install latest'
'gnvm install x.xx.xx'
'gnvm install x.xx.xx --global'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gnvm install args include " + strings.Join(args, " "))
		fmt.Println("global flag is " + strconv.FormatBool(global))
		//TO DO
	},
}

// sub cmd
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall local node.js version",
	Long: `uninstall local node.js version like :
'gnvm uninstall x.xx.xx'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gnvm uninstall args include " + strings.Join(args, " "))
		//TO DO
	},
}

// sub cmd
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use the specific version by current cmd( temp )",
	Long: `use the specific version by current cmd( temp ) like :
'gnvm use x.xx.xx'`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("gnvm use args include " + strings.Join(args, " "))
		//fmt.Println("global flag is " + strconv.FormatBool(global))

		if len(args) == 1 {

			if args[0] != "latest" && nodehandle.VerifyNodeVersion(args[0]) != true {
				fmt.Println("Use parameter support 'latest' or 'x.xx.xx', please check your input. See 'gnvm help use'.")
				return
			}

			// set use
			if ok := nodehandle.Use(args[0]); ok == true {
				// set global version
				config.SetConfig(config.GLOBAL_VERSION, nodehandle.GetTrueVersion(args[0]))
			}
		} else {
			fmt.Println("Use parameter maximum is 1, please check your input. See 'gnvm help use'.")
		}
	},
}

// sub cmd
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update global node.js",
	Long: `update global node.js like :
'gnvm update'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gnvm update args include " + strings.Join(args, " "))
		//TO DO
	},
}

// sub cmd
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list show all local | remote node.js version",
	Long: `list show all local | remote node.js version like :
'gnvm ls'
'gnvm ls --remote'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gnvm ls args include " + strings.Join(args, " "))
		fmt.Println("remote flag is " + strconv.FormatBool(remote))
		//TO DO
	},
}

// sub cmd
var nodeVersionCmd = &cobra.Command{
	Use:   "node-version",
	Short: "show global | current | latest node.js version",
	Long: `show global | current | latest node.js version like :
'gnvm node-version'
'laest version is x.xx.xx'
'global version is x.xx.xx'
'current version is x.xx.xx'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gnvm node-version args include " + strings.Join(args, " "))
		//TO DO
	},
}

// sub cmd
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "set | get registry and noderoot value",
	Long: `set | get registry and noderoot value like :
'gnvm config registry'
'registry is http://nodejs.org/dist/'
'gnvm config registry http://dist.u.qiniudn.com/'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			fmt.Println("gnvm config [" + args[0] + "] is " + config.GetConfig(args[0]))
		} else if len(args) == 2 {
			switch args[0] {
			case "registry", "noderoot":
				newValue := config.SetConfig(args[0], args[1])
				fmt.Println("Set success, [" + args[0] + "] new value is " + newValue)
			default:
				fmt.Println("Config parameter include 'registry' | 'noderoot', your input unknown, please check your input. See 'gnvm help config'.")
			}
		} else if len(args) > 2 {
			fmt.Println("Config parameter maximum is 2, please check your input. See 'gnvm help config'.")
		}
	},
}

func init() {
	// get node.exe root
	noderoot := nodehandle.GetGlobalNodePath()
	// set node.exe root to .gnvmrc
	config.SetConfig(config.NODEROOT, noderoot)
}

func Exec() {

	// add sub cmd to root
	gnvmCmd.AddCommand(versionCmd)
	gnvmCmd.AddCommand(installCmd)
	gnvmCmd.AddCommand(uninstallCmd)
	gnvmCmd.AddCommand(useCmd)
	gnvmCmd.AddCommand(updateCmd)
	gnvmCmd.AddCommand(lsCmd)
	gnvmCmd.AddCommand(nodeVersionCmd)
	gnvmCmd.AddCommand(configCmd)

	// flag
	installCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "get this version global version")
	//useCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "get this version global version")
	lsCmd.PersistentFlags().BoolVarP(&remote, "remote", "r", false, "get remote all node.js version list")

	// exec
	gnvmCmd.Execute()
}
