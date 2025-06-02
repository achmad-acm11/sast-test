package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sast-integration/app/helper"
	"sast-integration/app/shareVar"
)

type SemgrepCli struct {
	command string
	stdLog  *helper.StandartLog
}

func NewSemgrepCli() *SemgrepCli {
	return &SemgrepCli{
		stdLog: helper.NewStandardLog(shareVar.SemgrepCli, shareVar.Repository),
	}
}

func (t *SemgrepCli) Init() *SemgrepCli {
	t.command = "semgrep scan --config=p/default --json"
	return t
}

func (t *SemgrepCli) AddProjectPath(value string) *SemgrepCli {
	t.command += fmt.Sprintf(" %s", value)
	return t
}

func (t *SemgrepCli) AddOutput(value string) *SemgrepCli {
	t.command += fmt.Sprintf(" > %s", value)
	return t
}

func (t *SemgrepCli) Exec() ([]byte, error) {
	t.stdLog.InfoFunction("Executing command: " + t.command)

	cmd := exec.Command("bash", "-c", t.command)
	cmd.Stdin = os.Stdin
	return cmd.CombinedOutput()
}
