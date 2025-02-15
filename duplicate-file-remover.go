package main

import (
	"flag"
	"fmt"

	argument_list "github.com/Kilemonn/Duplicate-File-Remover/argument-list"
	"github.com/Kilemonn/Duplicate-File-Remover/consts"
	"github.com/Kilemonn/Duplicate-File-Remover/files"
)

func main() {
	var argList argument_list.ArgumentList
	flag.Var(&argList, consts.SHORT_INPUT_DIR, "Input directories, please provide the flag multiple times for as many input directories as you require.")
	shortOutputDir := flag.String(consts.SHORT_OUTPUT_DIR, consts.DEFAULT_OUTPUT_DIR, "Output directory")
	retainModifiedTime := flag.Bool(consts.RETAIN_MODIFIED_DATE, false, "Retains the modified date of the original file in the copied destination files. This can be required if you wish the output files modified date to match that of the original file.")

	flag.Parse()

	if len(argList.Args) < 2 {
		fmt.Printf("Please provide multiple input directories as command-line flag [-%s], atleast 2 are required.\n", consts.SHORT_INPUT_DIR)
		return
	}

	fmt.Println("When duplicate files are detected, files will be taken from the directories in the order they were provided.")

	err := files.CreateOutputDirectory(*shortOutputDir)
	if err != nil {
		fmt.Printf("Failed to create output directory [%s]. Error: %s.\n", *shortOutputDir, err.Error())
		return
	}

	err = files.MergeFileDirs(argList.Args, *shortOutputDir, *retainModifiedTime)
	if err != nil {
		fmt.Printf("Error encountered. %s\n", err.Error())
	}
}
