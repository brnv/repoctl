package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"
)

var version = "1.0"

var usage = `repoctl - repod client.

Description required.

Usage:
    repoctl [options] [-v] -L
    repoctl [options] [-v] -L <repo> <epoch> <db> <arch>
    repoctl [options] [-v] (-A|-S|-E|-R) <repo> <epoch> <db> <arch> [<package>] [--file=<path>]
    repoctl -h | --help
    repoctl --version

Options:
    --repod-address=<url>     Address to repod daemon.
                               [default: _repod]
    --repod-version=<string>  Address to repod daemon.
                               [default: v1]
    -L --list                 List repositories or repository packages.
    -A --add                  Add package to specified repository.
    -S --show                 Show package information.
    -E --edit                 Edit package in specified repository.
    -R --remove               Remove package from specified repository.
      <repo>                  Specify repository name.
      <epoch>                 Specify repository epoch.
      <db>                    Specify repository db.
      <arch>                  Specify repository architecture.
      <package>               Specify package to manipulate with.
      --file=<path>           Package file to add or edit in repository.
    --json                    Print output in json format.
    -h --help                 Show this help.
`

func main() {
	args, err := docopt.Parse(usage, nil, true, "repoctl "+version, false)
	if err != nil {
		panic(err)
	}

	var (
		repodAddress = args["--repod-address"].(string)
		repodVersion = args["--repod-version"].(string)

		modeList   = args["--list"].(bool)
		modeAdd    = args["--add"].(bool)
		modeShow   = args["--show"].(bool)
		modeEdit   = args["--edit"].(bool)
		modeRemove = args["--remove"].(bool)

		repo, _        = args["<repo>"].(string)
		epoch, _       = args["<epoch>"].(string)
		db, _          = args["<db>"].(string)
		arch, _        = args["<arch>"].(string)
		packageName, _ = args["<package>"].(string)
		packageFile, _ = args["--file"].(string)

		jsonOutput = args["--json"].(bool)
	)

	client := NewRepodClient(repodAddress, repodVersion)
	{
		if modeList || modeShow {
			client.method = "GET"
		}

		if modeAdd {
			client.method = "POST"
		}

		if modeRemove {
			client.method = "DELETE"
		}

		if modeEdit {
			client.method = "PATCH"
		}
	}

	if packageFile != "" {
		currentDirectory, err := os.Getwd()
		if err != nil {
			reportError(err)
		}

		err = client.LoadPackageFile(currentDirectory + "/" + packageFile)
		if err != nil {
			reportError(err)
		}

		if packageName == "" {
			packageName = packageFile
		}
	}

	client.appendURLParts([]string{
		repo,
		epoch,
		db,
		arch,
		packageName,
	})

	apiResponse, err := client.Do()
	if err != nil {
		reportError(err)
	}

	output := apiResponse.String()

	if jsonOutput {
		output, err = apiResponse.toJSON()
		if err != nil {
			reportError(err)
		}
	}

	if output != "" {
		fmt.Println(output)
	}
}

func reportError(err error) {
	fmt.Printf("%s", err.Error())
	os.Exit(1)
}
