package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"time"
	"encoding/json"
)

type JSON struct {
	Items []FileEvent `json:"Items"`
}

type Event struct {
	Message   string `json:"Message"`
	Timestamp time.Time `json:"Timestamp"`
}

type FileEvent struct {
	Namespace string `json:"Namespace"`
	Events    []Event `json:"Event"`
}

func main() {

	fmt.Println("-- Syndeno collector is running --")
	fmt.Println("-- This app collects events every 14 min --")

	for range time.Tick(time.Minute * 14) {
		go func() {
			_, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST")

			var config *rest.Config
			var err error

			if inCluster {
				config, err = rest.InClusterConfig()
				if err != nil {
					panic(err.Error())
				}
			} else {
				fmt.Println("[ERROR]: This code requires to be in a kubernetes pod")
			}

			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				panic(err.Error())
			}

			ns := namespaces(clientset)

			eventCollector(clientset, ns)

			length := needsCompress()
			
			if length == 99 {
				fmt.Println("Compress files")
			}

		}()
	}

}

func namespaces(clientset *kubernetes.Clientset) []string {
	array := []string{}

	nsList, _ := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})

	for _, v := range nsList.Items {
		array = append(array, v.Name)
	}

	return array
}

func eventCollector(clientset *kubernetes.Clientset, ns []string) {

	JsonEvents := FileEvent{}
	finaljson := JSON{}
	Array := []FileEvent{}
	EventArray := []Event{}

	for _, v := range ns {
		events, _ := clientset.CoreV1().Events(v).List(context.TODO(), metav1.ListOptions{})
		for _, item := range events.Items {
			EventArray = append(EventArray, Event{Message: item.Message, Timestamp: item.CreationTimestamp.Time})
		}
		JsonEvents = FileEvent{Namespace: v, Events: EventArray}
		Array = append(Array, JsonEvents)
		EventArray = []Event{}
	}

	finaljson = JSON{Items: Array}

	JSON, _ := json.Marshal(finaljson)

	fileTime := fmt.Sprint(time.Now().Year()) + "-" + fmt.Sprint(time.Now().Month()) + "-" + fmt.Sprint(time.Now().Day()) + "-" + fmt.Sprint(time.Now().Hour()) + "-" + fmt.Sprint(time.Now().Minute()) + "-" + fmt.Sprint(time.Now().Second())

	fileName := "kubernetes-events-" + fileTime + ".log"

	file, _ := os.Create("/app/logs/" + fileName)

	defer file.Close()

	_, err := file.WriteString(string(JSON))
	if err != nil {
        fmt.Println(err)
    }

	fmt.Println("Kubernetes event log " + fileName + " created")

}

func needsCompress() int{

	dirLength := []int{}

	dir, err := os.Open("/app/logs")
	if err != nil {
		fmt.Println("[ERROR]: Can't open directory /app/logs")
	}

	files, _ := dir.ReadDir(0)

	for i := range files {
		dirLength = append(dirLength, i)
	}

	length := len(dirLength)

	return length
}