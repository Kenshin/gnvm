package command

import (
	// lib
	"github.com/spf13/cobra"

	// go
	"fmt"
	//"strconv"
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
	Long: `install any node.js version e.g.
'gnvm install latest'
'gnvm install x.xx.xx y.yy.yy'
'gnvm install x.xx.xx --global'`,
	Run: func(cmd *cobra.Command, args []string) {
		var newArgs []string

		if len(args) == 0 {
			fmt.Println("Error: 'gnvm install' need parameter, please check your input. See 'gnvm help install'.")
		} else {
			for _, v := range args {

				// check latest
				if v == config.LATEST {
					newArgs = append(newArgs, v)
					continue
				}

				// check version format
				if ok := nodehandle.VerifyNodeVersion(v); ok != true {
					fmt.Printf("Error: [%v] format error, the correct format is x.xx.xx. \n", v)
				} else {
					newArgs = append(newArgs, v)
				}
			}
			nodehandle.Install(newArgs, global)
		}
	},
}

// sub cmd
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall local node.js version",
	Long: `uninstall local node.js version e.g.
gnvm uninstall x.xx.xx
gnvm uninstall latest
gnvm uninstall 0.10.26 0.11.2 latest`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("gnvm uninstall args include " + strings.Join(args, " "))
		if len(args) == 0 {
			fmt.Println("Error: 'gnvm uninstall' need parameter, please check your input. See 'gnvm help uninstall'.")
		} else {
			for _, v := range args {

				// get true version
				v = nodehandle.GetTrueVersion(v, true)

				// check version format
				if ok := nodehandle.VerifyNodeVersion(v); ok != true {
					fmt.Printf("Error: [%v] format error, the correct format is x.xx.xx. \n", v)
				} else {
					nodehandle.Uninstall(v)
				}
			}
		}
	},
}

// sub cmd
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use the specific version by global",
	Long: `use the specific version by global e.g.
'gnvm use x.xx.xx'
'gnvm use latest'`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("gnvm use args include " + strings.Join(args, " "))
		//fmt.Println("global flag is " + strconv.FormatBool(global))

		if len(args) == 1 {

			if args[0] != "latest" && nodehandle.VerifyNodeVersion(args[0]) != true {
				fmt.Println("Use parameter support 'latest' or 'x.xx.xx', e.g. 0.10.28, please check your input. See 'gnvm help use'.")
				return
			}

			// set use
			if ok := nodehandle.Use(args[0]); ok == true {
				// set global version
				config.SetConfig(config.GLOBAL_VERSION, nodehandle.GetTrueVersion(args[0], false))
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
	Long: `update global node.js e.g.
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
	Long: `list show all local | remote node.js version e.g.
gnvm ls
gnvm ls --remote`,
	Run: func(cmd *cobra.Command, args []string) {

		// check args
		if len(args) > 0 {
			fmt.Println("Warning: gnvm ls no parameter, please check your input. See 'gnvm help ls'.")
		}

		if remote == true {
			nodehandle.LsRemote()
		} else {
			// check local ls
			nodehandle.LS()
		}
	},
}

// sub cmd
var nodeVersionCmd = &cobra.Command{
	Use:   "node-version",
	Short: "show global | latest node.js version",
	Long: `show global | latest node.js version e.g.
gnvm node-version
Node.exe global verson is [x.xx.xx]
Node.exe latest verson is [x.xx.xx]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Println("Warning: gnvm node-version no parameter, please check your input. See 'gnvm help node-version'.")
		}
		nodehandle.NodeVersion()
	},
}

// sub cmd
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "set | get registry and noderoot value",
	Long: `set | get registry and noderoot value e.g.
'gnvm config registry'
'registry is http://nodejs.org/dist/'
'gnvm config registry http://dist.u.qiniudn.com/'`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			fmt.Println("gnvm config [" + args[0] + "] is " + config.GetConfig(args[0]))
		} else if len(args) == 2 {

			if args[1] == "default" {
				fmt.Println("Waring: Please use capital letter 'DEFAULT'.")
				args[1] = "DEFAULT"
			}

			switch {
			case args[0] == "registry" && args[1] != "DEFAULT":
				if newValue := config.SetConfig(args[0], args[1]); newValue != "" {
					fmt.Println("Set success, [" + args[0] + "] new value is " + newValue)
				}
			case args[0] == "registry" && args[1] == "DEFAULT":
				if newValue := config.SetConfig(args[0], config.REGISTRY_VAL); newValue != "" {
					fmt.Println("Registry reset success, [" + args[0] + "] new value is " + newValue)
				}
			case args[0] == "noderoot":
				fmt.Printf("Waring: [%v] Temporarily does not support.\n", args[0])
			default:
				fmt.Println("Config parameter include <registry>, your input unknown, please check your input. See 'gnvm help config'.")
			}
		} else if len(args) > 2 {
			fmt.Println("Config parameter maximum is 2, please check your input. See 'gnvm help config'.")
		}
	},
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
