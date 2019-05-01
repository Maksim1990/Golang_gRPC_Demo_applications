package main

import (
	"context"
	"fmt"
	"github.com/Maksim1990/grpcLearnExample/calculator/calculatorpb"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"io"
	"time"
)

func main() {
	fmt.Println("Calculator GRPC Client")
	fmt.Println("Client started on localhost:50051")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := caclulatorpb.NewCalculatorServiceClient(cc)

	//Call gRPC Unary func
	doUnary(c)

	//Call server streaming func
	doServerStreaming(c)

	//Call client streaming func
	doClientStreaming(c)

	//Call Bi-Di streaming func
	doBiDiStreaming(c)
}

func doBiDiStreaming(c caclulatorpb.CalculatorServiceClient) {
	fmt.Println("Starting FindMaximum BiDi Client streaming gRPC")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening BiDi FindMaximum RPC Streaming func: %v", err)
	}

	wait := make(chan struct{})

	go func() {
		numbers := []int32{4, 12, 6, 2, 40, 56, 78, 67}
		for _, numb := range numbers {
			fmt.Printf("Sending number %v in order to find MAX value \n", numb)
			stream.Send(&caclulatorpb.FindMaximumRequest{
				Number: numb,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		//Finish sending client's stream
		stream.CloseSend()
	}()

	go func() {

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				fmt.Print("Close channel & stop receiving responses \n")
				wait <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving BiDi FindMaximum RPC Streaming response: %v", err)
				break
			}

			fmt.Printf("Current max value is: %v \n", resp.GetMaximum())
		}
	}()

	<-wait
}
func doClientStreaming(c caclulatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Client streaming gRPC")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening RPC Streaming func: %v", err)
	}

	numbers := []int32{5, 6, 7}
	for _, numb := range numbers {
		fmt.Printf("Sending number %v in order to count average \n", numb)
		stream.Send(&caclulatorpb.ComputeAverageRequest{
			Number: numb,
		})
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet RPC func: %v", err)
	}
	fmt.Printf("Average value is: %v", resp.GetResult())
}

func doServerStreaming(c caclulatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Server streaming gRPC")

	request := &caclulatorpb.PrimeNumberDecompositionRequest{
		Number: 127597663,
	}
	resp, err := c.PrimeNumberDecomposition(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling server Streaming RPC func: %v", err)
	}

	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			//Reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from stream %v \n", msg.GetPrimeFactor())
	}
}
func doUnary(c caclulatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Sum Unary gRPC")

	request := &caclulatorpb.SumRequest{
		FirstNumber:  2,
		SecondNumber: 33,
	}

	resp, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Sum ROC func: %v", err)
	}

	log.Printf("The result of Sum gRPC calculation is: %v \n", resp.SumResult)
}
