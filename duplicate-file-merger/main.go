package main

import (
	argument_list "duplicate-file-merger/argument-list"
	"duplicate-file-merger/consts"
	"duplicate-file-merger/files"
	"flag"
	"fmt"
)

func main() {
	var argList argument_list.ArgumentList
	flag.Var(&argList, consts.SHORT_INPUT_DIR, "Input directories, please provide the flag multiple times for as many input directories as you require.")
	shortOutputDir := flag.String(consts.SHORT_OUTPUT_DIR, consts.DEFAULT_OUTPUT_DIR, "Output directory")

	flag.Parse()

	if len(argList.Args) < 2 {
		fmt.Printf("Please provide multiple input directories as command-line flag [-%s], atleast 2 are required.", consts.SHORT_INPUT_DIR)
		return
	}

	fmt.Println("When duplicate files are detected, files will be taken from the directories in the order they were provided.")

	fmt.Println(consts.INPUT_DIR, ":", argList.Args)
	fmt.Println(consts.SHORT_OUTPUT_DIR, ":", *shortOutputDir)

	err := files.CreateOutputDirectory(*shortOutputDir)
	if err != nil {
		fmt.Printf("Failed to create output directory [%s]. Error: %s.\n", *shortOutputDir, err.Error())
		return
	}

	err = files.MergeFileDirs(argList.Args, *shortOutputDir)
	if err != nil {
		fmt.Printf("Error encountered. %s\n", err.Error())
	}
}
