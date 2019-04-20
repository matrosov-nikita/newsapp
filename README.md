# News App

This app is used to create news and receive it by ID.

### Run

Run mongo and nats services:
```
$ docker-compose -f docker-compose.db.yml up -d
```

Run client and storage services:
```
$ docker-compose -f docker-compose.services.yml up -d
```

## API

By default server is running on http://localhost:8888

#### POST /news

Creates news. 
 
Request:  
```json
{
  "header": "news header"
}
```

Response:
```json
{
    "id":"5cbad50a9871a8d0f6390b2b"
}
```


#### GET /news/:ID

Returns news by id.  

Response:
```json
{
    "id": "5cbad50a9871a8d0f6390b2b",
    "header": "test news",
    "createdAt": "2019-04-20T08:15:06.653Z"
}
```
