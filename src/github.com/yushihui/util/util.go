package util

import (
	"io/ioutil"
	"os"
	"flag"
	"fmt"
)

const PLAYBOOK_EXAMPLE = `- name: Hello World Sample
  hosts: all
  tasks:
    - name: Hello Message
      debug:
        msg: "Hello World!"
    - name: ping google
      command: ping -c 3 google.com`
const INVENTORY = `[myos]
192.168.4.125
[myos:vars]
ansible_connection=local
ansible_ssh_user=alvin
ansible_ssh_pass=alvin`

const tmp = "/tmp/ansible/jobs/"

var port, host *string

func write2File(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0664)
}

func SavePlaybook(job string, playbook string, content []byte) error {

	jobPath := tmp + job
	return write2File(jobPath+"/"+playbook+".yml", content)
}

func SaveInventory(job string, content []byte) error {
	jobPath := tmp + job
	return write2File(jobPath+"/inventory", content)
}

func SaveLog(job string, content []byte) error {
	jobPath := tmp + job
	return write2File(jobPath+"/log.txt", content)
}

func CreateJobDir(job string) (path string, err error) {
	path = tmp + job
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0700)
		}
	}
	return path, err
}

func Init() {
	host = flag.String("host", "192.168.1.177", "rpc server")
	port = flag.String("port", "50051", "rpc port")
	flag.Parse()
}

func GetAddr() string {
	return fmt.Sprintf("%s:%s", *host, *port)
}
