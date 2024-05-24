package main

import (
	_ "github.com/pquerna/ffjson/fflib/v1"
	"github.com/pquerna/ffjson/generator"
	_ "github.com/pquerna/ffjson/inception"

	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var outputPathFlag = flag.String("w", "", "Write generate code to this path instead of ${input}_ffjson.go.")
var goCmdFlag = flag.String("go-cmd", "", "Path to go command; Useful for `goapp` support.")
var importNameFlag = flag.String("import-name", "", "Override import name in case it cannot be detected.")
var forceRegenerateFlag = flag.Bool("force-regenerate", false, "Regenerate every input file, without checking modification date.")
var resetFields = flag.Bool("reset-fields", false, "When unmarshalling reset all fields missing in the JSON")

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [options] [input_file]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s generates Go code for optimized JSON serialization.\n\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

var extRe = regexp.MustCompile(`(.*)(\.go)$`)

func main() {
	flag.Parse()
	extra := flag.Args()

	if len(extra) != 1 {
		usage()
	}

	inputPath := filepath.ToSlash(extra[0])

	var outputPath string
	if outputPathFlag == nil || *outputPathFlag == "" {
		outputPath = extRe.ReplaceAllString(inputPath, "${1}_ffjson.go")
	} else {
		outputPath = *outputPathFlag
	}

	var goCmd string
	if goCmdFlag == nil || *goCmdFlag == "" {
		goCmd = "go"
	} else {
		goCmd = *goCmdFlag
	}

	var importName string
	if importNameFlag != nil && *importNameFlag != "" {
		importName = *importNameFlag
	}

	err := generator.GenerateFiles(goCmd, inputPath, outputPath, importName, *forceRegenerateFlag, *resetFields)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s:\n\n", err)
		os.Exit(1)
	}

	println(outputPath)
}
