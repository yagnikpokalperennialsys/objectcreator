## Architectural block for object creator
![image](https://github.com/yagnikpokalperennialsys/objectcreator/assets/148773637/4b1bd294-4a64-45d8-a9bf-ef6d34b19f78)
### To generate the arm 64 binary

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
```
### To make the zip file
```
zip main.zip ./main
```
### GET object

```
curl --location 'https://jp1pqmyhn9.execute-api.eu-central-1.amazonaws.com/dev/34ce89b2-92c6-47ce-8b71-d3f0f31bd2b0'
```
### Response
![image](https://github.com/yagnikpokalperennialsys/objectcreator/assets/148773637/18652b19-e35d-4f77-9f82-27e6914992a5)
### POST object
``` 
curl --location 'https://jp1pqmyhn9.execute-api.eu-central-1.amazonaws.com/dev/' \
--header 'Content-Type: text/plain' \
--data 'yagnikp'
```
### Response
![image](https://github.com/yagnikpokalperennialsys/objectcreator/assets/148773637/f5ce0bda-4939-402a-8a70-b62463f8a341)

### Objects in S3
![image](https://github.com/yagnikpokalperennialsys/objectcreator/assets/148773637/e6b767fe-9e57-4049-887a-a5beda97555e)

### API gateway configurations
![image](https://github.com/yagnikpokalperennialsys/objectcreator/assets/148773637/d72c536a-5333-4db6-afd5-6aac9f1a85ae)

#### Objectcreatorhandler
![image](https://github.com/yagnikpokalperennialsys/objectcreator/assets/148773637/7d435e6f-a0a7-41de-8dbc-6590daea2d9e)
Referances

- https://stackoverflow.com/questions/54353860/publish-message-to-sns-with-aws-go-sdk
- https://github.com/awslabs/aws-lambda-go-api-proxy/blob/master/gin/adapter.go

