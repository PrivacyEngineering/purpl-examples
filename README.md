# PEng Group 7 - Summer term 2023 - examples
### Topic: Hook in privacy capabilities for gRPC

## [/playground](/playground)
### [/interceptors](playground/interceptors)
Find first steps in the ./playground/interceptors directory.
As of right now it's a modified version of the [go-grpc-middleware Repo](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/v2.0.0-rc.5).
To to run, 
```
cd playground/interceptors/examples
go run server/main.go
go run client/main.go
```
Wait a few seconds and then stop the server (```ctrl + C```).

Changes to server/main.go:
- removed existing interceptors
- added own interceptor
- added own selector.MatchFunc

### [/pingpong](playground/pingpong)
Simple example for how to modify the server's response using a gRPC interceptor.
Implements data minimzation in forms of reduction, noising & generalization.
