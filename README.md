## go-split-merge [![Go Reference][1]][2]
`go-split-merge` is an open-source tool for splitting/joining large files required to improve portability over email or GitHub commits.

## Features
- Split large files into small files of your desired filesize, real fast.
- Use the utility straight from your command line.
- Utility generates a text file with file hashes making it easy for users to verify that the end file isn't corrupted.

## Usage
1. Clone the repository:
```bash
git clone https://github.com/teocci/go-split-merge.git
```

2. Build the project as follows:
```bash
go build main.go
```

2.A. Run the split or merge command 
```bash
main split -f ./big_file.zip
```
or 
```bash
main merge -s ./big_file.zip.pt00
```

2.B. We can build and run the project like this:
```bash
go run main.go split -f big_file.zip
```

4. For more help, use the `-h` flag
```sh
go run main.go -h
```

Output:
```sh
go-split-merge v1.0.5
This application split/join large files to improve portability over email or github commits.

Usage:
  go-split-merge [flags]
  go-split-merge [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  merge       merge files into the original file
  split       Split files into small files

Flags:
  -h, --help      help for go-split-merge
  -v, --verbose   Run in debug mode
  -V, --version   Print version info and exit (default true)

Use "go-split-merge [command] --help" for more information about a command.
```

**NOTE:** Please take special care of the `GOROOT` and `GOPATH` variables. In case you get errors, simply create a new folder named `gosj` in the `GOPATH` and copy `main.go` file & the `sj` folder into the `gosj` directory. This should get the utility working fine.

----
#### In case you face trouble, please feel free to open an issue.

[1]: https://pkg.go.dev/badge/github.com/teocci/go-split-merge.svg
[2]: https://pkg.go.dev/github.com/teocci/go-split-merge
[3]: https://github.com/teocci/go-split-merge/releases/latest