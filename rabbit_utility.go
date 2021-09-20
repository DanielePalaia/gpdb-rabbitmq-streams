package main

import (
	
	"fmt"
	"log"
	"os"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream" // Main package
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp" // amqp 1.0 package to encode messages
	"sync/atomic"
	"strconv"
)

type rabbitClient struct {
	// This is the connection string to all the stream servers
	connString  string
	streamName   string
	// this one count the number of item processed
	count       int32
	// Total number of batches
	size        int32
	// Keep the batched items
	buffer      []string
	gpssclient  *gpssClient
	env         *stream.Environment
	// Type of offset we want to read from
	offset      string
}

func makeRabbitClient(connString string, streamName string, size int32, offset string, gpssclient *gpssClient) *rabbitClient {
	client := new(rabbitClient)
	client.connString = connString
	client.streamName = streamName
	client.count = 0
	client.size = size
	client.buffer = make([]string, size)
	client.gpssclient = gpssclient
	client.offset = offset

	return client
}

func (client *rabbitClient) failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (client *rabbitClient) connect() () {

	log.Printf("streamName: %v", client.connString)

	addresses := []string{
		client.connString}

	var err error
	log.Printf("connecting to rabbit")
	client.env, err = stream.NewEnvironment(
		stream.NewEnvironmentOptions().SetUris(addresses))

	if err != nil  {
		log.Fatalf("Error connecting to the rabbitmq stream engine %s", err)
	}

	log.Printf("declaring stream")

	err = client.env.DeclareStream(client.streamName,
		stream.NewStreamOptions().
		SetMaxLengthBytes(stream.ByteCapacity{}.GB(2)))

	/*if err != nil  {
		log.Fatalf("Error declaring the stream error: %v stream name: %s", err, client.streamName)
	}*/

}

func (client *rabbitClient) consume() {

	handleMessages := func(consumerContext stream.ConsumerContext, message *amqp.Message) {

		
		client.buffer[client.count] = string(message.Data[0])
		fmt.Printf("consumer name: %s, text: %s \n ", consumerContext.Consumer.GetName(), client.buffer[client.count])

		if atomic.AddInt32(&client.count, 1)%client.size == 0 {
			
				log.Printf("Batch reached: I'm sending request to write to gpss/gprc server")
				client.gpssclient.ConnectToGreenplumDatabase()
				client.gpssclient.WriteToGreenplum(client.buffer)
				client.gpssclient.DisconnectToGreenplumDatabase()
				client.count = 0

				// AVOID to store for each single message, it will reduce the performances
				// The server keeps the consume tracking using the consumer name
				err := consumerContext.Consumer.StoreOffset()
				if err != nil {
					CheckErr(err)
				}
				
		}
	
}

	startOffset := stream.OffsetSpecification{}.First()

	if client.offset == "first"   {
		startOffset = stream.OffsetSpecification{}.First()
	} else if client.offset == "last" {
		startOffset = stream.OffsetSpecification{}.Last()
	}  else if client.offset == "lastconsumed" {
		startOffset = stream.OffsetSpecification{}.LastConsumed()
	} else  {
		offsetInt, _ := strconv.Atoi(client.offset)
		offsetInt64 := int64(offsetInt)
		startOffset = stream.OffsetSpecification{}.Offset(offsetInt64)
	}	

	consumer, err := client.env.NewConsumer(
			client.streamName,
			handleMessages,
			stream.NewConsumerOptions().
			SetConsumerName("gpss").                 // set a consumer name
			SetOffset(startOffset))

			CheckErr(err)

	channelClose := consumer.NotifyClose()
	// channelClose receives all the closing events, here you can handle the
	// client reconnection or just log
	defer consumerClose(channelClose)

}


func consumerClose(channelClose stream.ChannelClose) {
	event := <-channelClose
	fmt.Printf("Consumer: %s closed on the stream: %s, reason: %s \n", event.Name, event.StreamName, event.Reason)
}



func CheckErr(err error) {
	if err != nil {
		fmt.Printf("%s ", err)
		os.Exit(1)
	}
}