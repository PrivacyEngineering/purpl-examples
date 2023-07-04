# 🏓 pingpong: two clients, one server

| who? | what? |
| ----------- | ----------- |
| goodclient | sends request to server |
| badclient | sends request to server |
| server | sends an address as response |
| interceptor | minimizes the address depending on the client JWT |

## 🧪 Try it
```
go run server/server.go
go run client/clients.go
```
change the data minimization technique in the [```server/server.go```](server/server.go) interceptor function.

## 🥸 Data minimization opions 
- reduction
- noising
- generalization

A reduced result looks like this:
```
2023/06/23 18:39:08 Message from server: 
2023/06/23 18:39:08 Street: Straße des 17 Juni 
2023/06/23 18:39:08 Number: 135
2023/06/23 18:39:09 Message from server for badclient: 
2023/06/23 18:39:09 Street: Straße des 17 Juni 
2023/06/23 18:39:09 Number: -1
```
A noised result looks like this:
```
2023/06/23 18:49:53 Message from server: 
2023/06/23 18:49:53 Street: Straße des 17 Juni 
2023/06/23 18:49:53 Number: 135
2023/06/23 18:49:54 Message from server for badclient: 
2023/06/23 18:49:54 Street: Straße des 17 Juni 
2023/06/23 18:49:54 Number: 145
```
A generalized (floored to the lower end of it's 10s-interval, e.g. 135 -> 131 or 99 --> 91) result looks like this:
```
2023/06/23 18:50:17 Message from server: 
2023/06/23 18:50:17 Street: Straße des 17 Juni 
2023/06/23 18:50:17 Number: 135
2023/06/23 18:50:18 Message from server for badclient: 
2023/06/23 18:50:18 Street: Straße des 17 Juni 
2023/06/23 18:50:18 Number: 131
```

## 🔑 Use of JSON Web Tokens

Check and generate them here: [jwt.io](https://jwt.io/).

Our token's secret: ```none```.

Right now our JWTs look like this:

```
payload for goodToken:
{
 	"policy": {
 	  "allowed": {
 		"name": "string",
 		"sex": "string"
 	  },
 	  "generalized": {
 		"phoneNumber": "string"
 	  },
 	  "noised": {
 		"age": "int"
 	  },
 	  "reduced": {
 		"street": "string"
 	  }
 	},
 	"exp": 1688843806,
 	"iss": "test"
}

payload for badToken:
{
 	"policy": {
 	"allowed": {},
 	"generalized": {
 		"age": "int",
 		"name": "string",
 		"phoneNumber": "string",
 		"sex": "string",
 		"street": "string"
 	},
 	"noised": {},
   	"reduced": {}
 	},
 	"exp": 1688483421,
 	"iss": "test"
}
```


- The clients append their respective JWTs to their request's context.
- The server's gRPC interceptor compares the outgoing response with the JWT's ```allowed``` and ```minimized```data fields. Allowed fields will be left untouched. Minimzed fields will be minimzed. Unmentioned fields will be reduced to 1 or nil

## 🧭 Rodamap
- ...
