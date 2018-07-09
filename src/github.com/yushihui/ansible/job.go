package ansible // import "github.com/yushihui/job"

import "time"

type Job struct {
	ID         string `db:"id" json:"id"`
	Status string `db:"status" json:"status"`
	Debug  bool   `db:"debug" json:"debug"`
	DryRun bool `db:"dry_run" json:"dry_run"`
	Playbook    string `db:"playbook" json:"playbook"`
	Environment string `db:"environment" json:"environment"`
	Created time.Time  `db:"created" json:"created"`
	Start   *time.Time `db:"start" json:"start"`
	End     *time.Time `db:"end" json:"end"`
}

// ansible log output
type JobOutput struct {
	JobId string       `db:"job_id" json:"job_id"`
	Job   string    `db:"job" json:"job"`
	Time   time.Time `db:"time" json:"time"`
	Output string    `db:"output" json:"output"`
}

type Inventory struct {

}
