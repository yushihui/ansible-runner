package ansible

import (
	"os/exec"
	"strconv"
	"strings"

	"encoding/json"
)

type job struct {
	job        Job
	inventory   Inventory
	hosts       []string
}

func (j *job) runGalxy() error{
	return nil
}


func (j *job) runPlaybook() error {
	args, err := j.getPlaybookArgs()
	if err != nil {
		return err
	}
	cmd := exec.Command("ansible-playbook", args...) //nolint: gas
	cmd.Dir = util.Config.TmpPath + "/repository_" + j.job.Id
	cmd.Env = j.envVars(util.Config.TmpPath, cmd.Dir, nil)

	j.logCmd(cmd)
	cmd.Stdin = strings.NewReader("")
	return cmd.Run()
}

func (j *job) getPlaybookArgs() ([]string, error) {
	playbookName := j.job.Playbook
	var inventory string
	args := []string{
		"-i", inventory,
	}

	if j.job.Debug {
		args = append(args, "-vvvv")
	}

	if j.job.DryRun {
		args = append(args, "--check")
	}

	args = append(args, playbookName)


	return args, nil
}


