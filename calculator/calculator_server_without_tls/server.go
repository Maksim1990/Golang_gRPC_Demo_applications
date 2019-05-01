package main

import (
	"context"
	"fmt"
	"github.com/Maksim1990/grpcLearnExample/calculator/calculatorpb"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
)

type server struct {
}

func (*server)FindMaximum(stream caclulatorpb.CalculatorService_FindMaximumServer) error{
	fmt.Printf("Received FindMaximum RPC Client's streaming \n ")


	maximum:=int32(0)
	for{
		req,err:=stream.Recv()
		if err ==io.EOF{
			return nil
		}
		if err !=nil{
			log.Fatalf("Error while reading client stream: %v \n",err)
			return err
		}

		number:=req.GetNumber()
		if number>maximum{
			maximum=number

			sendErr:=stream.Send(&caclulatorpb.FindMaximumResponse{
				Maximum:maximum,
			})
			if sendErr !=nil{
				log.Fatalf("Error while sending response to the  client: %v \n",err)
				return  err
			}
		}
	}
}
func (*server)ComputeAverage(stream caclulatorpb.CalculatorService_ComputeAverageServer) error  {
	fmt.Printf("Received CalculateAverage RPC \n ")
		sum:=int32(0)
		count:=0
		for{
			req,err:=stream.Recv()
			if err ==io.EOF{
				average:=float64(sum)/float64(count)
				return stream.SendAndClose(&caclulatorpb.ComputeAverageResponse{
					Result: average,
				})
			}
			if err !=nil{
				log.Fatalf("Error while reading client stream: %v \n",err)
			}
			sum+=req.GetNumber()
			count++
		}
}

func (*server)PrimeNumberDecomposition(req *caclulatorpb.PrimeNumberDecompositionRequest,stream caclulatorpb.CalculatorService_PrimeNumberDecompositionServer) error{
	fmt.Printf("Received number: %v \n", req)
	number := req.GetNumber()
	divisor:=int64(2)

	for number >1{
		if number % divisor == 0{
			stream.Send(&caclulatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor:divisor,
			})
			number=number/divisor
		}else {
			divisor++
			fmt.Printf("Increased divisor %v \n",divisor)
		}
	}
	return nil
}

func (*server) Sum(ctx context.Context, req *caclulatorpb.SumRequest) (*caclulatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with: %v \n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	result:=firstNumber+secondNumber
	res:=&caclulatorpb.SumResponse{
		SumResult:result,
	}

	return  res,nil
}

func main() {
	fmt.Println("Calculator GRPC Server")
	fmt.Println("Server started on 127.0.0.1:50051")

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	caclulatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
