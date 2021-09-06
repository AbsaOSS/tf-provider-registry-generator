package git

import (
	"github.com/AbsaOSS/gopkg/shell"
	"log"
)

func RunGit(args []string, dir string) (err error) {
	cmd := shell.Command{
		Command: "git",
		Args: args,
		WorkingDir: dir,
	}
	gitLog, err := shell.Execute(cmd)
	if err != nil {
		return
	}
	log.Printf("run git %+v: %s", args, gitLog)
	return
}
