package ansible

import (
	"os/exec"
	"strconv"
	"strings"

	"encoding/json"
)

type task struct {
	job        Job
	inventory   Inventory
	hosts       []string
}

func (t *task) runGalxy() error{
	return nil
}

func (t *task) getPlaybookArgs() error{
	return nil
}

func (t *task) runPlaybook() error {
	args, err := t.getPlaybookArgs()
	if err != nil {
		return err
	}
	cmd := exec.Command("ansible-playbook", args...) //nolint: gas
	cmd.Dir = util.Config.TmpPath + "/repository_" + strconv.Itoa(t.repository.ID)
	cmd.Env = t.envVars(util.Config.TmpPath, cmd.Dir, nil)

	t.logCmd(cmd)
	cmd.Stdin = strings.NewReader("")
	return cmd.Run()
}

func (t *task) getPlaybookArgs() ([]string, error) {
	playbookName := t.job.Playbook
	var inventory string
	args := []string{
		"-i", inventory,
	}

	if t.job.Debug {
		args = append(args, "-vvvv")
	}

	if t.job.DryRun {
		args = append(args, "--check")
	}

	if len(t.environment.JSON) > 0 {
		var js map[string]interface{}
		err := json.Unmarshal([]byte(t.environment.JSON), &js)
		if err != nil {
			t.log("JSON is not valid")
			return nil, err
		}

		extraVar, err := removeCommandEnvironment(t.environment.JSON, js)
		if err != nil {
			t.log("Could not remove command environment, if existant it will be passed to --extra-vars. This is not fatal but be aware of side effects")
		}

		args = append(args, "--extra-vars", extraVar)
	}
	args = append(args, playbookName)


	return args, nil
}


