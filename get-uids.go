package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Set the kubeconfig file path (optional)
	kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig file")
	flag.Parse()

	// Create a Kubernetes client
	config, err := getKubeConfig(*kubeconfig)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Get namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting namespaces: %v", err)
	}

	// Print namespaces and UIDs
	for _, ns := range namespaces.Items {
		fmt.Printf("Namespace: %s\nUID: %s\n\n", ns.Name, ns.UID)
	}
}

func getKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	// Use the in-cluster configuration if kubeconfigPath is not provided
	if kubeconfigPath == "" {
		return rest.InClusterConfig()
	}

	// Use the provided kubeconfig file
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

