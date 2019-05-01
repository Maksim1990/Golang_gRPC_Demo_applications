# Golang_gRPC_Demo_applications
GO demo applications with all available gRPC types

### About Application
- Consist of 3 separate GO+gRPC demo apps
- In each demo app integrated examples of following gRPC calls:
  - Unary 
  - Server streaming
  - Client streaming
  - Bidirectional streaming

### How To Run each application (example for **Greet** demo app)

1) Start server 

```
go run greet/greet_server/server.go
```

2) Start client app and see out put of all gRPC call types

```
go run greet/greet_client/client.go 
```

### In order to test all other app (**Calculator** and **Blog CRUD + MongoDB**) run as well both server and client respectively.

# How to test gRPC calls in [Postman](https://www.getpostman.com/)
### CHECK  [INSTRUCTIONS](https://github.com/Maksim1990/Symfony_APP_Twitter/blob/Build_app_guide/public/github/APP_GUIDE.md)
