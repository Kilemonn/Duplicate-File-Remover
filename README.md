# Duplicate-File-Merger
A command-line tool that will take two input directories and create an output directory containing only unique files from the two input directories. The files are determined as being unique based on its content hash.

This program will only create a new specified output directory. Or if none is provided a directory in the path of the running application called `output`.

## Usage

The only required command line flag values into this program is the input directories, provided by multiple `-i` command line argument values.
Optionally the output directory can be provided with `-o`.

You can get application usage by using `-h` or `--help`.

E.g.

> ./duplicate-file-merger -i input/directory/1 -i input/directory/2 -o output/directory
