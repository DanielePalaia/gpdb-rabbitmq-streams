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
