# Rabbitmq streams - Greenplum GPSS connector

This project is a modification of my previous work: </br> 
Rabbitmq streams - Greenplum GPSS connector</br> 
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
      test=# create table companies(id varchar (200), city varchar (200), foundation timestamp, description text, data json);
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

