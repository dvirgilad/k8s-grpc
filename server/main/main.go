package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net"
	"path/filepath"

	pb "github.com/dvirgilad/grpcNode/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9999, "The server port")
)

// server is used to implement .
type server struct {
	pb.UnimplementedNodeServiceServer
	kubernetes.Clientset
}

func GetKubeConfig() *rest.Config {
	var kubeconfig string

	if flag.Lookup("kubeconfig") == nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
	} else {
		kubeconfig = flag.Lookup("kubeconfig").Value.String()

	}

	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil
	}

	return config
}

func (s *server) GetNodes(ctx context.Context, _ *pb.NodeRequest) (*pb.NodeResponse, error) {
	config := GetKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err

	}
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	log.Print("Received NodeRequest")

	var nodeList []*pb.Node
	for _, clusterNode := range nodes.Items {
		isReady := false
		for _, condition := range clusterNode.Status.Conditions {
			if condition.Type == "Ready" {
				switch condition.Status {
				case v1.ConditionTrue:
					isReady = true
				case v1.ConditionFalse:
					isReady = false
				default:
					isReady = false

				}
			}
		}
		log.Printf("found node: %s", clusterNode.Name)

		nodeList = append(nodeList, &pb.Node{
			Name:    clusterNode.Name,
			Version: clusterNode.Status.NodeInfo.KubeletVersion,
			Ready:   isReady,
		})

	}
	return &pb.NodeResponse{Nodes: nodeList}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
