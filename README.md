# Anzer product

Service can be used with anzer project to make production concatenation of services.
It uses RabbitMQ.

## Usage

1. Set environments:

    ```
    RMQ=amqp://guest:guest@172.17.0.1:5672/
    IN1=in1
    IN2=in2
    OUT=out
    TYPE1=a
    TYPE2=b
    ```
    And start it by `go build && ./anzer-product` or with docker.
2. Create queue and bind it to routing key `out`.
3. Now send to routing key `in1` any json-valid message with header `pid=1`.
4. Send another json-valid message to `in2` with the same `pid` header.
5. Wait for message like

    ```json
    {
        "a": "data from the input 1",
        "b": "data from the input 2"
    }
    ```
    With header `pid=1`.
 
6. Profit!