

#Shrtnr

A very basic Golang URL shortener backed by Redis.

There is a basic Angular frontend written but this is intended to be consumed as a restful API.




## API


### Create URL

`POST /`

data:

```javascript
{
  "LongURL":    "a long url for a webpage",
  "RequestURL": "the requested short url"   // optional
}

```


### URL Stats

`GET /stats/[short url]`


Response:

(more data points may be added)

```javascript
{
  "HitCount": X
}

```

### Forward Via Short URL

`GET /[short url]`

Any get request that is not to `/`, `/static/*` or `/favicon.ico` will be interpreted as a forward request.



