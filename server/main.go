package main

import (
	"assignment/pb"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Report struct {
	UserID    string
	ReportID  string
	CreatedAt time.Time
}

var userIDs = []string{"user1", "user2", "user3"}

type server struct {
	pb.UnimplementedAssignmentServiceServer
	reports map[string]Report
}


func(s *server) GetHealth(ctx context.Context, req *emptypb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: "ok"}, nil
}


func(s *server) GenerateReport(ctx context.Context, req *pb.UserRequest) (*pb.ReportResponse, error) {
	
	report := Report{
		UserID:    req.UserId,
		ReportID:  generateRand(req.UserId),
		CreatedAt: time.Now(),
	}

	fmt.Println(report)

	s.reports[report.ReportID] = report

	// log.Printf("[generated report for %s: %s", report.UserID, report.ReportID)
	
	return &pb.ReportResponse{
		UserId:    report.UserID,
		ReportId:  report.ReportID,
		CreatedAt: report.CreatedAt.Format(time.RFC3339),

	}, nil
}

func generateRand(userID string) string {
	rand.Seed(time.Now().UnixNano())
	return userID + "-" + strconv.Itoa(rand.Intn(10000000))
}


func startCronJob(s *server) {
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		for i, userID := range userIDs {
			fmt.Println("<=========="+ strconv.Itoa(i+1) +"==========>")
			_, err := s.GenerateReport(context.Background(), &pb.UserRequest{UserId: userID})
			if err != nil {
				log.Printf("[Cron] Error generating report for %s: %v", userID, err)
			}
		}
	})
	c.Start()
}

func main(){
	listener , err := net.Listen("tcp", ":8080")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	s:= grpc.NewServer()
	reflection.Register(s)

	srv := &server{reports: make(map[string]Report)}

	pb.RegisterAssignmentServiceServer(s, srv)

	//cron job
	startCronJob(srv)
	
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)  
	}
}