package ansible // import "github.com/yushihui/ansible"

import (
	pb "github.com/yushihui/ansible-pb"
	"log"
)

func (j *Job) send(msg string, status pb.JobStatus) {
	//now := time.Now()
	log.Println("msg to send:" + msg)
	j.msgs <- &pb.AnsibleJobResponse{Message: msg, Status: status}

}
