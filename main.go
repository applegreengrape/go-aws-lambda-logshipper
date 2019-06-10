package main

import (
	"strings"
	"log"
	"net"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func SocketClient(m []byte) {
	conn, err := net.Dial("tcp", "your_tcp_endpoint:your_port")

	defer conn.Close()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("sending logs to syslog-ng:", string(m))
	conn.Write(m)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func Handler(request events.CloudwatchLogsEvent) error {

	cloudwatchLogsData, err := request.AWSLogs.Parse()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("cloudwatchLogsData.LogEvents:", cloudwatchLogsData.LogEvents)

	// Remove the prefix to get to the name of the Lambda.
	logGroup := strings.Replace(cloudwatchLogsData.LogGroup, "/aws/lambda/", "", 1)

	// What you want to capture is up to you. BUT, for example:
	type LogEntry struct {
		LogGroup  string `json:"log_group"`
		Timestamp int64  `json:"timestamp"`
		Message   string `json:"message"`
	}

	// Stuff the incoming log lines into the datastructure to serialize to Log Entries.
	for _, event := range cloudwatchLogsData.LogEvents {
		logEntry := LogEntry{LogGroup: logGroup, Timestamp: event.Timestamp, Message: event.Message}
	//	logEntry := LogEntry{Timestamp: event.Timestamp, Message: event.Message}
		j, err := json.Marshal(logEntry)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		
		// Send to syslog-ng over tcp, basically.
		SocketClient(j)
		fmt.Println(j)
	}

	return nil
}

func main() {
	// Entrypoint that the Lambda will execute.
	lambda.Start(Handler)
}
