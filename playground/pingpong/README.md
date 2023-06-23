# 🏓 pingpong: two clients, one server

| who? | what? |
| ----------- | ----------- |
| goodclient | sends request to server |
| badclient | sends request to server |
| server | sends an address as response |
| interceptor | minimizes the address depending on the client name |

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



## 🧭 Rodamap
- use proper attribute for minimzation decision (probably JWT)
- ...
