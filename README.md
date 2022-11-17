Run rabbitmq within docker container:

```dockerfile
docker run -d --name some-rabbit -p 5672:5672 rabbitmq:3.10.7-management-alpine
```

if you want to point any specific user and pass, then you can pass env variables which control this user/pass flow

```dockerfile
docker run -d --name some-rabbit -p 5672:5672 -e RABBITMQ_DEFAULT_USER=someuser -e RABBITMQ_DEFAULT_PASS=somepass rabbitmq:3.10.7-management-alpine
```

## todo:
* what is routing key
* what is exhange type
* try to stop container with rabbit and play with durability option
* what is exchange and what type it could accept?

*what is consumer parameter in channel.Consume:*

this is identifier for consumer which might be used for cancelling delivering messages to consumer

*what is autoAck parameter:*

it means that message will be acked immediately after consuming is done

*what is exclusive parameter:*

when exclusive is true, the server will ensure that this is the sole consumer from this queue. Also, good to note: Exclusive queues may only be accessed by the current connection, and are deleted when that connection closes

*what is nowait parameter:*

do not wait a response from server

-- what is consumer argument in .Consumer("queueName", "consumerName") function

this way you can cancel such consumer by channel.Cancel("consumer id")



