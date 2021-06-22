// NOTE this command line tool only works at PROJECT ROOT directory.
// This command line tool is used to initialize project of different
// program languages in vscode.

package main

import (
	"flag"
	"fmt"
	"os"

	"local/src/golang"
	"local/src/js"
	"local/src/python"
	"local/src/ts"
	"local/src/util"
)

const languages = "go/py/ts/js"

func helpMsg() {
	fmt.Println("please specify language -", languages)
	fmt.Println("usage: vsinit <language> [<args>]")
	fmt.Println("eg: vsinit go")
}

func main() {
	if err := util.CheckCMDInstall("code"); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if len(os.Args) < 2 {
		helpMsg()
		os.Exit(2)
	}

	// flag.ExitOnError will os.Exit(2) if subcommand Parse() error.
	tsjsSet := flag.NewFlagSet("tsjs", flag.ExitOnError)
	jestflag := tsjsSet.Bool("jest", false, "add 'jest' - unit test components")

	var err error
	switch os.Args[1] {
	case "go":
		err = golang.InitProject()
	case "py":
		err = python.InitProject()
	case "ts":
		err = ts.InitProject(tsjsSet, jestflag)
	case "js":
		err = js.InitProject(tsjsSet, jestflag)
	default:
		helpMsg()
		os.Exit(2)
	}

	// 统一打印 error
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
