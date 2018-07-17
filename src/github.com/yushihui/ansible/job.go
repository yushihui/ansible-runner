package ansible // import "github.com/yushihui/ansible"

import "time"



// ansible log output
type JobOutput struct {
	JobId string       `db:"job_id" json:"job_id"`
	Job   string    `db:"job" json:"job"`
	Time   time.Time `db:"time" json:"time"`
	Output string    `db:"output" json:"output"`
}

type Inventory struct {

}
