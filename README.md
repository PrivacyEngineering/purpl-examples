<p align="center">
	<img src="purpl.png" width=50" />
</p>

# purpl - examples
**We used this repository for early stage evaluations and testing (mainly in the performance branch).**
## [/playground](/playground)
In the [/playground](/playground) directory you'll find two examples showcasing the purpl-interceptor's functionality.

### Example 1: [/pingpong](playground/pingpong)
Simple example for how to modify the server's response using a gRPC interceptor.
Implements data minimzation in forms of reduction, noising & generalization.

#### 🏓 two clients, one server

| who? | what? |
| ----------- | ----------- |
| goodclient | sends request to server |
| badclient | sends request to server |
| server_two | sends a name, phone number, street, age and sex as response |
| interceptor | minimizes the response depending on the client JWT |

To find out more and try it out, head over to [/pingpong](playground/pingpong).

### Example 2: [/interceptors](playground/interceptors)
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

