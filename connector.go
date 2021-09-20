package main

/* Product main file*/
import (
	"log"
	"strconv"
)

func main() {

	/****This part will just read properties configuration file inside the same folder******/

	log.Printf("Starting the connector and reading properties in the properties.ini file")
	/* Reading properties from ./properties.ini */
	prop, _ := ReadPropertiesFile("./properties.ini")
	port, _ := strconv.Atoi(prop["GreenplumPort"])
	batch, _ := strconv.Atoi(prop["batch"])
	// 1 if batch is persistent, 0 otherwise
	
	log.Printf("Properties read: Connecting to the Grpc server specified")

	/* Connect to the grpc server specified */
	gpssClient := MakeGpssClient(prop["GpssAddress"], prop["GreenplumAddress"], int32(port), prop["GreenplumUser"], prop["GreenplumPassword"], prop["Database"], prop["SchemaName"], prop["TableName"])
	gpssClient.ConnectToGrpcServer()

	log.Printf("Connected to the grpc server")

	log.Printf("Connecting to rabbit and consuming")
	/* Generate teh rabbit connection */
	rabbit := makeRabbitClient(prop["rabbitstream"], prop["streamName"], int32(batch), prop["offset"], gpssClient)
	rabbit.connect()
	rabbit.consume()

}
