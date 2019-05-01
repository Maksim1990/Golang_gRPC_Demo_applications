package main

import (
	"context"
	"fmt"
	"github.com/Maksim1990/grpcLearnExample/greet/greetpb"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"time"
)

func main() {
	fmt.Println("Hello GRPC Client")
	fmt.Println("Client started on localhost:50051")

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//Call gRPC Unary func
	doUnary(c)

	//Call server streaming func
	doServerStreaming(c)

	//Call server streaming func
	doClientStreaming(c)

	//Call Bi-Di streaming func
	doBiDiStreaming(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Bi-Di Client streaming gRPC")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC func: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Maksim",
				LastName:  "Test 1",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Helen",
				LastName:  "Test 2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
				LastName:  "Test 3",
			},
		},
	}

	wait := make(chan struct{})

	//Sending in parallel to the server requests
	go func() {
		for _, req := range requests {
			//Send each request to the server
			fmt.Printf("Sending request: %v \n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//Getting in parallel server responses
	go func() {
		for {
			resp, err := stream.Recv()

			if err == io.EOF {
				wait <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalf("Error while reading GreetEveryone RPC responce: %v", err)
				break
			}
			fmt.Printf("Received response: %v \n", resp.GetResult())
		}

	}()
	<-wait

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Client streaming gRPC")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC func: %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Maksim",
				LastName:  "Test 1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Helen",
				LastName:  "Test 2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
				LastName:  "Test 3",
			},
		},
	}

	//Iterate over slice and send each messages individually
	for _, req := range requests {
		fmt.Printf("Sending request: %v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving responce from LongGreet RPC func: %v", err)
	}

	fmt.Printf("LongGreet response: %v", resp)
}
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Server streaming gRPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "TestName",
			LastName:  "TestLast",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC func: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//Reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from stream %v \n", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary gRPC")

	greeting := &greetpb.Greeting{
		FirstName: "TestName",
		LastName:  "TestLast",
	}

	req := &greetpb.GreetRequest{
		Greeting: greeting,
	}

	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC func: %v", err)
	}

	log.Printf("Response from Greet gRPC: %v", resp.Result)
}
