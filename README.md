# sport-articles-keeper

## Assumptions made while developing
As long as no additional incrowd API params where provided, the application assumes that continuous calls of `https://www.htafc.com/api/incrowd/getnewlistinformation?count=50` might return new articles as they were added to incrowd.

So, the application constantly calls the same `GET /getnewlistinformation` in order to find an article that wasn't stored in DB yet.

Once that article(s) found the app calls `GET /getnewsarticleinformation` to get full article text and save it all together in the DB.

If no new articles were found in `GET /getnewlistinformation` response no article is saved.

If new API params will be introduced, it would be easily to add into the app.

## How to run
The only one things you need to run ths app is Docker(and docker-compose) installed.

### Run 
```shell
docker-compose up -d
```

### Stop
```shell
docker-compose down
```

## The API introduced by the app

### Get multiple articles
Parameters - standard pagination stuff:
* `limit` how many documents to return
* `skip` how many documents to skip
```shell
curl --location '0.0.0.0:8080/articles?limit=2&skip=5'
```

### Get single article by ID
```shell
curl --location '0.0.0.0:8080/articles/6ded1955-6338-47eb-b8d4-0e16f9026141'
```