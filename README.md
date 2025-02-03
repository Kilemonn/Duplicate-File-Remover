# Duplicate-File-Remover

[![Go Coverage](https://github.com/Kilemonn/Duplicate-File-Remover/wiki/coverage.svg)](https://raw.githack.com/wiki/Kilemonn/Duplicate-File-Remover/coverage.html)

A command-line tool that takes input directories and create an output directory containing only unique files from the provided input directories. The files are determined as being unique based on its content hash.

This program will only create a new specified output directory. Or if none is provided a directory in the path of the running application called `output`.

## Installation

This application can be installed onto the command line with the following command:
> go install github.com/Kilemonn/Duplicate-File-Remover@latest

## Usage

The only required command line flag values into this program is the input directories, provided by multiple `-i` command line argument values.
Optionally the output directory can be provided with `-o`.

You can get application usage by using `-h` or `--help`.

E.g.
> Duplicate-File-Remover.exe -i input/directory/1 -i input/directory/2 -o output/directory

## Example

Given the following input directories:
```
input1/
    |--> text.txt (content equal to input2/text.txt)
    |--> image1.jpg (content is unique across both directories)
    |--> test1.png (content equal to input2/test2.txt)

input2/
    |--> text.txt (content equal to input1/text.txt)
    |--> image2.jpg (content is unique across both directories)
    |--> test2.png (content equal to input1/test1.txt)
```

After running the following command:
> Duplicate-File-Remover.exe -i input1/ -i input2/ -o output/

The output directory will contain the following files:
```
output/
    |--> text.txt (content is equal so only 1 is copied over)
    |--> image1.jpg (since content is unique across both directories)
    |--> image2.jpg (since content is unique across both directories)
    |--> test1.png (because the content is equal. Since "input1" was specified first in the command line arguments its file name will take precedence)
```

**The contents of the two input directories (`input1` and `input2`) remain untouched at the end of the file copying.**
