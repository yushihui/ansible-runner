package rpc

import (
	"log"
	"time"
	pb "github.com/yushihui/ansible-pb"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}


func (s *server) StartAutoMationJob(in *pb.AnsibleJob, jobServer pb.AnsibleExecutor_StartAutoMationJobServer) (error) {

	log.Printf("received a job id: %s", in.Id)
	msg := pb.AnsibleJobResponse{Message: "running", Status: pb.JobStatus_Running, JobPriority: pb.JobPriority_High}
	jobServer.SendMsg(&msg)
	time.Sleep(2 * time.Second)
	msg2 := pb.AnsibleJobResponse{Message: "Success", Status: pb.JobStatus_Success, JobPriority: pb.JobPriority_High}
	jobServer.SendMsg(&msg2)
	log.Printf("job %s done", in.Id)
	return nil
}

func start(){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("rpc server started...")
	s := grpc.NewServer()
	pb.RegisterAnsibleExecutorServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Fatalf("server started to listen")
}
