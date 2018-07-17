package ansible // import "github.com/yushihui/ansible"

import (
	"os/exec"
	pb "github.com/yushihui/ansible-pb"
	"time"
	"log"
	"bufio"
	"bytes"
	"io"
	"github.com/yushihui/util"
)

func Run(ansibleJob *pb.AnsibleJob, msgs chan *pb.AnsibleJobResponse) error {
	defer func() {
		time.Sleep(3 * time.Second)
		close(msgs)
	}()
	dir, err := util.CreateJobDir(ansibleJob.Id)
	if err != nil {
		return err
	}

	j := &Job{Id: ansibleJob.Id, Playbook: ansibleJob.Name, msgs: msgs, Debug: true, DryRun: true, Workspace: dir}
	j.persistJob2Tmp()
	j.Start = time.Now()
	if err := j.runPlaybook(); err != nil {
		j.send("Running playbook failed: "+err.Error(), pb.JobStatus_Failed)
		j.End = time.Now()
		return err
	}

	j.End = time.Now()
	j.send("Running playbook finished: ", pb.JobStatus_Success)
	return nil
}

func (j *Job) persistJob2Tmp() { // todo get from rpc client
	util.SavePlaybook(j.Id, j.Playbook, []byte(util.PLAYBOOK_EXAMPLE))
	util.SaveInventory(j.Id, []byte(util.INVENTORY))
}

type Job struct {
	// private might be better
	Id          string    `db:"id" json:"id"`
	Status      string    `db:"status" json:"status"`
	Debug       bool      `db:"debug" json:"debug"`
	DryRun      bool      `db:"dry_run" json:"dry_run"`
	Playbook    string    `db:"playbook" json:"playbook"`
	Environment string    `db:"environment" json:"environment"`
	Created     time.Time `db:"created" json:"created"`
	Start       time.Time `db:"start" json:"start"`
	End         time.Time `db:"end" json:"end"`
	inventory   Inventory
	hosts       []string
	msgs        chan *pb.AnsibleJobResponse
	Workspace   string
}

func (j *Job) runGalxy() error {
	return nil
}

func (j *Job) runPlaybook() error {
	log.Println("start job .......")
	args, err := j.getPlaybookArgs()
	if err != nil {
		log.Println("start job failed.......")
		return err
	}
	cmd := exec.Command("ansible-playbook", args...)
	//cmd.Dir = util.Config.TmpPath + "/repository_" + j.job.Id
	cmd.Dir = j.Workspace
	//cmd.Env = j.envVars(util.Config.TmpPath, cmd.Dir, nil)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	var buffer, logBuffer bytes.Buffer
	ticker := time.NewTicker(1 * time.Second)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		buffer.WriteString(line)
		logBuffer.WriteString(line)
		select {
		case _ = <-ticker.C:
			j.send(buffer.String(), pb.JobStatus_Running)
			buffer.Reset()
		}
	}
	ticker.Stop()
	cmd.Wait()
	j.send(buffer.String(), pb.JobStatus_Running)
	util.SaveLog(j.Id, logBuffer.Bytes())
	//
	//j.logCmd(cmd)
	//cmd.Stdin = strings.NewReader("")
	//return cmd.Run()
	return nil
}

func (j *Job) getPlaybookArgs() ([]string, error) {
	playbookName := j.Playbook
	args := []string{
		"-i", "inventory",
	}

	//if j.Debug {
	//	args = append(args, "-vvvv")
	//}

	if j.DryRun {
		args = append(args, "--check")
	}

	args = append(args, playbookName+".yml")

	return args, nil
}
