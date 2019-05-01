package main

import (
	"fmt"
	"github.com/Maksim1990/grpcLearnExample/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"net"
	"log"
	"context"
	"strconv"
	"time"
)

type server struct {}

func (*server)GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error{
	fmt.Printf("GreetEveryone function was invoked with streaming request \n")

	for{
		req,err:=stream.Recv()
		if err==io.EOF{
			return nil
		}

		if err!=nil{
			log.Fatalf("Error while reading client's stream: %v \n",err)
			return err
		}

		firstName:=req.GetGreeting().GetFirstName()
		result:="Hello "+ firstName+"!"
		err=stream.Send(&greetpb.GreetEveryoneResponse{
			Result:result,
		})

		if err!=nil{
			log.Fatalf("Error while sending data to client : %v \n",err)
			return err
		}
	}
}

func (*server)LongGreet(stream greetpb.GreetService_LongGreetServer) error{
	fmt.Printf("Client streaming with LongGreet function with streaming request \n")
	result:= "Hello "
	for{
		req,err:=stream.Recv()
		if err==io.EOF{
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err!=nil{
			log.Fatalf("Error while reading client stream: %v\n",err)
		}
		firstName:=req.GetGreeting().GetFirstName()
		result+=firstName+"! "
	}
}

func (*server)GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error{
	fmt.Printf("Streaming with GreetManyTimes function was invoked with: %v \n", req)
	firstName:=req.GetGreeting().GetFirstName()
	for i:=0;i<10;i++{
		result:="Hello "+firstName+" "+strconv.Itoa(i)
		resp:=&greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(resp)
		time.Sleep(1000 *time.Millisecond)
	}

	return nil
}

func (*server) Greet( ctx context.Context,req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	fmt.Printf("Greet function was invoked with: %v \n", req)
	firstName:=req.GetGreeting().GetFirstName()
	result:= "Hello "+ firstName
	resp:=&greetpb.GreetResponse{
		Result: result,
	}
	return resp,nil
}
func main(){
	fmt.Println("Hello GRPC Server")
	fmt.Println("Server started on 0.0.0.0:50051")

	lis, err:= net.Listen("tcp","0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err:=s.Serve(lis); err !=nil{
		log.Fatalf("Failed to serve: %v", err)
	}
}
