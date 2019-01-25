#### 1. TCP packet design

| description | version | type | length | payload |
| ----------- | :-----: | :--: | :----: | :-----: |
| size(byte)  |    1    |  1   |   4    | length  |

version: the packet structure version.

type: the packet type. 

- 0 is BeartHeart
- 1 is Command
- 2 is Data

length: the length of payload. the value between 0-32KB(for 1Mbit/s network, 32KB data need about 250ms, 8KB data need about 60ms, so suggest the suitable size is 8-16KB). if the payload is too large, suggest split.

payload: the data you want to transfer

- type is BeartHeart: only version and type, without length and payload

- type is Command: 
  payload has two parts. first 4 bytes is message type(int32), the last is a JSON string which contains the real information of the message.

  | description | message type | data |
  | ----------- | ------------ | ---- |
  | size(bytes) | 4(int32)     |      |

  the follow is a base structure.

  {

  	id: string, 			//message id, a guid/uuid string to identify the message
		
  	version: "1.0.0.0"		//message version

  }

suggest command messages:

| command type            | description                                                  |
| :---------------------- | ------------------------------------------------------------ |
| ConnectMessage          | contains the information of client.                          |
| AcceptMessage           | response of ConnectMessage. can contains the information of server. |
| RejectMessage           | response of ConnectMessage. can contains the reason of reject. |
| DisconnectMessage       | close connection in normal cases. mainly send by client side. |
| PingMessage             | after server accept the client, test the business connection.data is "ping" and response is "pong". |
| TransferMessage         | the information of the transfer. should contains a guid/uuid(128bit) identify of the transfer process. |
| TransferCancelMessage   | cancel the transfer process.                                 |
| TransferCompleteMessage | after all the file data received, send this message as a ACK to let other know. |
| ErrorMessage            | information of errors. 0 means received data.                |
| CommandMessage          | can send some simple command to server side. may cause the security issue, if necessary should implement this  carefully. |
|                         |                                                              |



- type is Data: payload is byte stream. this packet type mainly use in file transferring

in the payload of the tcp packet

| description | transferid(GUID/UUID) | offset    | data  |
| ----------- | --------------------- | --------- | ----- |
| size(bytes) | 16                    | 8(uint64) | 0-16K |

when implement the file transfer, can split many litter packets.

**file transfer sequence**

```sequence
target->source:connect
source->target: accept
target->source:getfilelist
source->target:filelist(30-50items in once transfer.)
target->source:getfile(with uuid)
source->target:filedata(with uuid,target check hash)
target->source:transfercomplete
target->source:disconnect(then target close connection)
```



#### 2. REST API design

http://host/projectname/api/module/v1/xxxx?xxxx=xxxxx



#### 3. Security design

Key Agreement process

```sequence
    participant client1
    participant client2
    participant cloud
    client1->>cloud:upload PubKey1(https)
    client2->>cloud:upload PubKey2(https)
    cloud->>client1:download PubKey2(https)
    cloud->>client2:download PubKey1(https)
    client1->>client2:SessionKey signed by PriKey1,\nencrypt by PubKey2
    note over client2:SessionKey decrypt by PriKey2,\n verify signature by PubKey1
    client2->>client1:transfer data encrypted by SessionKey
    note over client1:decrypt data by SessionKey
```




#### 4. Data Store design

##### 4.1 User Config



##### 4.2 Cache

Redis, memcache 

##### 4.3 Database

mongoDB, InfluxDB, LevelDB, MySQL

#### 5. Client architecture design(ignore)

#### 6. Operation Maintenance(ignore)