package cli

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"path/filepath"
)

type CLI struct {
	projectSlug, configFile, username, password *string
	rootDir                                     string
}

func NewCLI() CLI {
	cli := CLI{
		projectSlug: flag.String("project", "", "REQUIRED - the transifex project slug"),
		configFile: flag.String("config", "", "REQUIRED - The location of the configuration file"),
		username: flag.String("username", "", "The transifex username"),
		password: flag.String("password", "", "The transifex password")}
	flag.Parse()
	if *cli.configFile == "" {
		fmt.Printf("The 'config' flag is required.  \n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *cli.projectSlug == "" {
		fmt.Printf("The 'project' flag is required.  \n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cli.rootDir = filepath.Dir(*cli.configFile)

	if !strings.HasSuffix(cli.rootDir, "/") {
		cli.rootDir = cli.rootDir + "/"
	}
	return cli
}

func (cli CLI) RootDir() string {
	return cli.rootDir
}

func (cli CLI) ProjectSlug() string {
	return *cli.projectSlug
}

func (cli CLI) ConfigFile() string {
	return *cli.configFile
}

func (cli CLI) Username() string {
	readAuth(cli.username, "username")
	return *cli.username
}

func (cli CLI) Password() string {
	readAuth(cli.password, "password")

	return *cli.password
}

func readAuth(field *string, prompt string) {

	if *field == "" {
		var line string
		var readlineErr error
		in := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter your %s: ", prompt)
		if line, readlineErr = in.ReadString('\n'); readlineErr != nil {
			log.Fatalf("Failed to read %s", prompt)
		}

		*field = strings.TrimSpace(line)
	}
}
