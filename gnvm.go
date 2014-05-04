// TestFlag project main.go
package main

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
)

func flagMode() {
	// golang的flag包的一些基本使用方法

	// 待使用的变量
	var id int
	var name string
	var male bool

	// 是否已经解析
	fmt.Println("parsed? = ", flag.Parsed())

	// 设置flag参数 (变量指针，参数名，默认值，帮助信息)
	// 也可以用以下带返回值的方法代替，不过他们返回的是指针，比较麻烦点
	// Int(name string, value int, usage string) *int
	// String(name string, value string, usage string) *string
	// Bool(name string, value bool, usage string) *bool
	flag.IntVar(&id, "id", 123, "help msg for id")
	flag.StringVar(&name, "name", "default name", "help msg for name")
	flag.BoolVar(&male, "male", false, "help msg for male")

	// 解析
	flag.Parse()

	// 是否已经解析
	fmt.Println("parsed? = ", flag.Parsed())

	// 获取非flag参数
	fmt.Println("------ Args start ------")
	for i, v := range flag.Args() {
		fmt.Printf("arg[%d] = (%s).\n", i, v)
	}
	fmt.Println("------ Args end ------")

	// visit只包含已经设置了的flag
	fmt.Println("------ visit flag start ------")
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f.Name, f.Value, f.Usage, f.DefValue)
	})
	fmt.Println("------ visit flag end ------")

	// visitAll只包含所有的flag(包括未设置的)
	fmt.Println("------ visitAll flag start ------")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Println(f.Name, f.Value, f.Usage, f.DefValue)

	})
	fmt.Println("------ visitAll flag end ------")

	// flag参数
	fmt.Printf("id = %d\n", id)
	fmt.Printf("name = %s\n", name)
	fmt.Printf("male = %t\n", male)

	// flag参数默认值
	fmt.Println("------ PrintDefaults start ------")
	flag.PrintDefaults()
	fmt.Println("------ PrintDefaults end ------")

	// 非flag参数个数
	fmt.Printf("NArg = %d\n", flag.NArg())
	// 已设置的flag参数个数
	fmt.Printf("NFlag = %d\n", flag.NFlag())
}

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

func main() {
	//flagMode()
	cobraMode()
}
