# Rabbitmq streams - Greenplum GPSS connector

This project is a modification of my previous work: </br> 
https://github.com/DanielePalaia/gpss-rabbit-greenplum-connector </br> 
to let the connector work also with the new rabbitmq streams functionality: </br>
https://blog.rabbitmq.com/posts/2021/07/rabbitmq-streams-overview</br>
This project should be considered experimental, still under development and non production ready  </br>
It is using this rabbitmq streams golang client </br>
https://github.com/rabbitmq/rabbitmq-stream-go-client </br>

The following reading can help you to better understand the software:

**RabbitMQ:** </br>
https://www.rabbitmq.com/ </br>
**RabbitMQ streams:** </br>
https://blog.rabbitmq.com/posts/2021/07/rabbitmq-streams-overview</br>
**GRPC:**  </br>
https://grpc.io/ </br>
**Greenplum GPSS:**</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/overview.html</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/api/dev_client.html</br>

GPSS is able to receive requests from different clients (Kafka, Greenplum-Informatica connector) as shown in the pic and proceed with ingestion process. We are adding support for RabbitMQ
![Screenshot](./pics/image2.png)

The connector will attach to a rabbitmq stream specified at configuration time will then batch a certain amount of elements specified and finally will ask the gpss server to push them on a greenplum table. </br>

## Prerequisites

1. **Activate the gpss extension on the greenplum database you want to use (for example test)**
   
      ```
      test=# CREATE EXTENSION gpss;
      ```
   
2. **Create the Greenplum table to be ingested**

      The table can be generic with any number of coloumns and data types. The important thing is that the input to ingest         will be coherent with the table definition. Let's try this table as example:
   
      ```
      test=# create table people(id int, name text, surname text, email text, profession text);
      ```

   ![Screenshot](./pics/definition.png)
   
3. **Run a gpss server with the right configuration (ex):**
  
      gpss ./gpsscfg1.json --log-dir ./gpsslogs
      where gpsscfg1.json 
  
      ```
      {
         "ListenAddress": {
            "Host": "",
            "Port": 50007,
            "SSL": false
         },
         "Gpfdist": {
            "Host": "",
            "Port": 8086
         }
      }
      ```

4. **download, install and run a rabbitmq broker**

      https://www.rabbitmq.com/download.html </br>
      then run the broker </br>
      ./rabbitmq-server </br>
      Then go with your browser to the rabbit web UI in: </br>
      http://localhost:15672/ </br></br>
      ![Screenshot](./pics/connection.png)<br/>
      and log with guest/guest (default)
      
## Running the application

1. **Find binaries** </br>
The application is written in GO. If you are using MacOs or Linux64 then you can directly use the binary version inside ./bin/osx and ./bin/linux of this project called: gpss-rabbit-greenplum-connect otherwise you must compile it with the GO compiler<br/>

2. **Setting property file**    
Use the file properties.ini (that should be place in the same directory of the binary in order to instruct the program        with this properties

```
GpssAddress=192.168.12.145:50009
GreenplumAddress=192.168.12.145
GreenplumPort=5432
GreenplumUser=gpadmin
GreenplumPasswd=
Database=dashboard
SchemaName=public
TableName=companies 
rabbitstream=rabbitmq-stream://daniele:daniele@localhost:5552/
streamName=gpss
batch=3
offset=lastconsumed
```

Batch is the amount of items the client will batch before sending the ingest request to gpss, and also it will commit the last read offset to the server when reaching this batch amount (see the tracking offset section)

## Tracking offset

Client offset are kept inside the server. Rabbitmq streams offers this functionality to store the last reading offset of a client in a server.
In the project this is tracked by the offset offset=lastconsumed property.
You can specify 4 values: </br>
**first**: the client will start reading from the first offset </br>
**last**: The client will read from the last offset (will basically wait for new messages to arrive on the stream) </br>
**lastconsumed**: The client will start from the last read and registered offset </br>
**an integer**: The client will start from the offset specified

## How to execute

You can just run the binary connector inside /bin/osx or /bin/linux in order to start the connector and create the rabbitmq stream

```
dpalaia@dpalaia-a02 osx % ./connector
2021/09/20 12:52:20 Starting the connector and reading properties in the properties.ini file
2021/09/20 12:52:20 Properties read: Connecting to the Grpc server specified
2021/09/20 12:52:20 Connected to the grpc server
2021/09/20 12:52:20 Connecting to rabbit and consuming
2021/09/20 12:52:20 streamName: rabbitmq-stream://daniele:daniele@localhost:5552/
2021/09/20 12:52:20 connecting to rabbit
2021/09/20 12:52:20 declaring stream
```

You should then see a queue of type stream with name you specified in the properties.ini created:

 ![Screenshot](./pics/stream.png)<br/>
 
 you can then publish the message using rabbitmq UI
 
  ![Screenshot](./pics/publish.png)<br/>
 
 Every batch items inserted on the stream you should find something like this on the server side:
 
2021/09/20 13:08:35 Batch reached: I'm sending request to write to gpss/gprc server
2021/09/20 13:08:35 connecting to a greenplum database
2021/09/20 13:08:35 Beginning to write to greenplum
2021/09/20 13:08:35 table informations
2021/09/20 13:08:35 prepare for writing
Result:  SuccessCount:3 

and then in the Greenplum table the items inserted:

```
dashboard=# select * from people;
 id |  name   | surname |   email    |  profession   
----+---------+---------+------------+---------------
  2 | daniele | palaia  | aaa@aaa.it | data engineer
  3 | daniele | palaia  | aaa@aaa.it | data engineer
  1 | daniele | palaia  | aaa@aaa.it | data engineer

```
 
 Play with the offset property in order to test the offset management from start, to the lastconsumed or to a random offset.
