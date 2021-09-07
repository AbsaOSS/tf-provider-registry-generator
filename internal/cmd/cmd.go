package cmd

import (
	"log"

	"github.com/AbsaOSS/gopkg/shell"
)

func Run(cmd string, args []string, dir string) (err error) {
	command := shell.Command{
		Command:    cmd,
		Args:       args,
		WorkingDir: dir,
	}
	cmdLog, err := shell.Execute(command)
	if err != nil {
		return
	}
	log.Printf("run %s %+v: %s", cmd, args, cmdLog)
	return
}
