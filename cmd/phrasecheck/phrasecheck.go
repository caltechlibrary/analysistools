package main

import (
	"fmt"
	"flag"
	"os"
	"path"

	// Project package
	"github.com/caltechlibrary/analysistools"
)

func main() {
	appName := path.Base(os.Args[0])
	helpText := analysistools.HelpText
	version := analysistools.Version
	releaseDate, releaseHash := analysistools.ReleaseDate, analysistools.ReleaseHash
	fmtHelp := analysistools.FmtHelp

	showHelp, showVersion, showLicense := false, false, false
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.Parse()
	args := flag.Args()

	if showHelp {
			fmt.Fprintf(os.Stdout, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
			os.Exit(0)
	}
	if showVersion {
			fmt.Fprintf(os.Stdout, "%s %s %s\n", appName, version, releaseHash)
			os.Exit(0)
	}
	if showLicense {
			fmt.Fprintf(os.Stdout, "%s\n", analysistools.LicenseText)
			os.Exit(0)
	}
	// Create phrasecheck app
	app := &analysistools.PhraseCheckApp{}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "missing action, don't know what to do\n")
		os.Exit(1)
	}
	action, params := args[0], args[1:]
	if err := app.Run(appName, action, params); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)	
	}
	os.Exit(0)
}