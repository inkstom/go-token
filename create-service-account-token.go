package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Create ServiceAccount and get token for each namespace
	for _, ns := range namespaces.Items {
		serviceAccountName := "custom-service-account" // Change this to your desired service account name
		createServiceAccount(clientset, ns.Name, serviceAccountName)
		token, err := getServiceAccountToken(clientset, ns.Name, serviceAccountName)
		if err != nil {
			log.Printf("Error getting ServiceAccount token for namespace %s: %v", ns.Name, err)
		} else {
			fmt.Printf("Namespace: %s\nServiceAccount: %s\nToken: %s\n\n", ns.Name, serviceAccountName, token)
		}
	}
}

func createServiceAccount(clientset *kubernetes.Clientset, namespace, serviceAccountName string) {
	// Create ServiceAccount
	_, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.Background(), &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating ServiceAccount in namespace %s: %v", namespace, err)
	}
}

func getServiceAccountToken(clientset *kubernetes.Clientset, namespace, serviceAccountName string) (string, error) {
	// Get ServiceAccount Secret
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("serviceaccount=%s", serviceAccountName),
	})
	if err != nil {
		return "", fmt.Errorf("error getting ServiceAccount secrets: %v", err)
	}

	// Find the token in the Secret
	for _, secret := range secrets.Items {
		if _, ok := secret.Data["token"]; ok {
			return string(secret.Data["token"]), nil
		}
	}

	return "", fmt.Errorf("token not found in ServiceAccount secret")
}

func getKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	// Use the in-cluster configuration if kubeconfigPath is not provided
	if kubeconfigPath == "" {
		return rest.InClusterConfig()
	}

	// Use the provided kubeconfig file
	return rest.InClusterConfig()
}

