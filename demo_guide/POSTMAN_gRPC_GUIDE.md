# POSTMAN + gRPC DEMO GUIDE

### Requirements for testing gRPC server with Postman
- Server should not use activated TLS mode
- Insatall and run **grpc-json-proxy**
```
go get -u github.com/jnewmano/grpc-json-proxy 
grpc-json-proxy
```
- Start Calculator demo app **without TLS mode** 
```
go run calculator/calculator_server_without_tls/server.go
```

### Example of Postman app settings to axtivate global proxy
![Mockup for feature A](https://github.com/Maksim1990/Golang_gRPC_Demo_applications/blob/master/demo_guide/1.PNG?raw=true)

### Example of request body for testing Sum method of Caculator service
![Mockup for feature A](https://github.com/Maksim1990/Golang_gRPC_Demo_applications/blob/master/demo_guide/2.PNG?raw=true)

### Example of headers for testing Sum method of Caculator service
![Mockup for feature A](https://github.com/Maksim1990/Golang_gRPC_Demo_applications/blob/master/demo_guide/3.PNG?raw=true)
