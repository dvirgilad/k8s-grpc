# k8s-grpc
This is a PoC of an application using gRpc to work with the k8s API using client-go


## Usage
* Clone the project
```bash
    git clone https://github.com/dvirgilad/k8s-grpc.git
    cd k8s-grpc
```
* Run the server
In a server with a valid `.kubecontext` file :
```bash
    go run server/main/main.go
```
The default port to access the server is 9999

* Run the client:
Change line 15 in client/main/main.go to point to the server
```bash
    go run client/main/main.go  
```
