package rpc //import "github.com/yushihui/rpc"

import (
	"log"
	pb "github.com/yushihui/ansible-pb"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
	"github.com/yushihui/ansible"
	"github.com/yushihui/util"
)

//const (
//	port = ":50051"
//)

type server struct{}

func (s *server) StartAnsibleJob(in *pb.AnsibleJob, jobServer pb.AnsibleService_StartAnsibleJobServer) (error) {

	log.Printf("received a job id: %s", in.Id)
	msgs := make(chan *pb.AnsibleJobResponse, 10)
	go ansible.Run(in, msgs)

	//go func() {
	msgs <- &pb.AnsibleJobResponse{Message: "running", Status: pb.JobStatus_Running}
	//	time.Sleep(2 * time.Second)
	//	msgs <- &pb.AnsibleJobResponse{Message: "Success", Status: pb.JobStatus_Success, JobPriority: pb.JobPriority_High}
	//	close(msgs)
	//}()

	for {
		select {
		case msg, ok := <-msgs:
			if ok {
				log.Printf(msg.Message)
				jobServer.SendMsg(msg)
			} else {
				log.Printf("Channel closed!") //
				return nil
			}
		default:
			log.Printf("No value ready, moving on.")
		}
		time.Sleep(500 * time.Millisecond)

	}

	log.Printf("job %s done", in.Id)
	return nil
}

func Start() {

	lis, err := net.Listen("tcp", util.GetAddr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("")
	log.Println("rpc server started...")
	// var opts []grpc.ServerOption
	//cre,er := credentials.NewClientTLSFromFile("/Users/dev/ssl/fsc.pem", "/Users/dev/ssl/fsc.key")
	//if er != nil {
	//	log.Fatalf("Failed to generate credentials %v", er)
	//}
	//opts = []grpc.ServerOption{grpc.Creds(cre)}
	//s := grpc.NewServer(grpc.Creds(cre))
	s := grpc.NewServer()

	pb.RegisterAnsibleServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Fatalf("server started to listen")
}
