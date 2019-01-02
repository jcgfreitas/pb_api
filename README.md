# pb_api

## API Calls

#### Get Coupon

##### GET /coupons/{id:[0-9]+}

This endpoint returns a coupon

#####Parameters

| Parameters | Required | Description      | Param type | Data type |
|------------|----------|------------------|------------|:---------:|
|     id     |    yes   | Coupon unique id |    path    |    uint   |

#####Http Status

|         Status        | Code |
|:---------------------:|------|
|           Ok          |  200 |
|        NotFound       |  404 |
| Internal Server Error |  500 |

#####Curl Example

curl -X GET http://localhost:8080/coupons/1 -i
```
HTTP/1.1 200 OK
Date: Wed, 02 Jan 2019 18:30:30 GMT
Content-Length: 194
Content-Type: text/plain; charset=utf-8

{
  "ID": 1,
  "CreatedAt": "2019-01-02T17:53:14.954953Z",
  "UpdatedAt": "2019-01-02T17:53:14.954953Z",
  "DeletedAt": null,
  "name": "CouponName",
  "brand": "CouponBrand",
  "value": 10,
  "expiry": "2020-01-01T23:59:59Z"
}
```

---

#### Create Coupon

##### POST /coupons

This endpoint creates an coupon

#####Parameters

| Parameters | Required | Description        | Param type | Data type |
|------------|----------|--------------------|------------|:---------:|
|    name    |    yes   | Coupon name        |    body    |   string  |
|    brand   |    yes   | Coupon brand       |    body    |   string  |
|    value   |    yes   | Coupon value       |    body    |   uint    |
|    expiry  |    yes   | Coupon expiry date |    body    |   string  |

expiry needs to be in `time.RFC3339` format

#####Http Status

|         Status        | Code |
|:---------------------:|------|
|        Created        |  201 |
|       BadRequest      |  400 |
| Internal Server Error |  500 |

#####Curl Example
curl -X POST --data '{"name" : "CouponName","brand" : "CouponBrand","value" : 10,"expiry" : "2020-01-01T23:59:59Z"}' http://localhost:8080/coupons -i

```
HTTP/1.1 201 Created
Date: Wed, 02 Jan 2019 18:26:48 GMT
Content-Length: 0
```
---

#### Delete Coupon

##### DELETE /coupons/{id:[0-9]+}

This endpoint deletes a coupon

#####Parameters

| Parameters | Required | Description      | Param type | Data type |
|------------|----------|------------------|------------|:---------:|
|     id     |    yes   | Coupon unique id |    path    |    uint   |

#####Http Status

|         Status        | Code |
|:---------------------:|------|
|           Ok          |  200 |
|        NotFound       |  404 |
| Internal Server Error |  500 |

#####Curl Example

curl -X DELETE http://localhost:8080/coupons/1 -i
```
HTTP/1.1 200 OK
Date: Wed, 02 Jan 2019 19:13:40 GMT
Content-Length: 0
```

#### Update Coupon

##### POST /coupons/{id:[0-9]+}

This endpoint updates a coupon

#####Parameters

| Parameters | Required | Description        | Param type | Data type |
|------------|----------|--------------------|------------|:---------:|
|     id     |    yes   | Coupon unique id   |    path    |    uint   |
|    name    |    no    | Coupon name        |    body    |   string  |
|    brand   |    no    | Coupon brand       |    body    |   string  |
|    value   |    no    | Coupon value       |    body    |   uint    |
|    expiry  |    no    | Coupon expiry date |    body    |   string  |

At least one of the body's elements is required

#####Http Status

|         Status        | Code |
|:---------------------:|------|
|           Ok          |  200 |
|       BadRequest      |  400 |
|        NotFound       |  404 |
| Internal Server Error |  500 |

#####Curl Example

curl -X POST --data '{"name" : "CouponName","brand" : "CouponBrand","value" : 10,"expiry" : "2020-01-01T23:59:59Z"}' http://localhost:8080/coupons/4 -i
```
HTTP/1.1 200 OK
Date: Wed, 02 Jan 2019 19:19:27 GMT
Content-Length: 0
```
---

#### Get Coupons

##### GET /coupons

This endpoint queries coupons and returns a slice of coupons

#####Parameters

| Parameters | Required |                Description/Query               | Param type | Data type |
|------------|:--------:|:----------------------------------------------:|------------|:---------:|
|    name    |    no    |                `WHERE name = ?`                |    query   |   string  |
|    brand   |    no    |                `WHERE brand = ?`               |    query   |   string  |
|    value   |    no    |                `WHERE value = ?`               |    query   |    uint   |
|    limit   |    no    |      limits the number of coupons received     |    query   |    uint   |
|    page    |    no    |  used to get the next batch of limited coupons |    query   |    uint   |
|     le     |    no    |      Lesser Than Expiry `WHERE expiry < ?`     |    query   |   string  |
|     ge     |    no    |     Greater Than Expiry `WHERE expiry > ?`     |    query   |   string  |
|     lc     |    no    |  Lesser Than Created_at `WHERE created_at < ?` |    query   |   string  |
|     gc     |    no    | Greater Than Created_at `WHERE created_at > ?` |    query   |   string  |
|     lv     |    no    |       Lesser Than Value `WHERE value < ?`      |    query   |    uint   |
|     gv     |    no    |       Greater Than Value `WHERE value > ?`     |    query   |    uint   |

#####Http Status

|         Status        | Code |
|:---------------------:|------|
|           Ok          |  200 |
|       BadRequest      |  400 |
| Internal Server Error |  500 |

#####Curl Example

curl -X GET http://localhost:8080/coupons?value=5 -i
```
HTTP/1.1 200 OK
Date: Wed, 02 Jan 2019 20:00:43 GMT
Content-Length: 787
Content-Type: text/plain; charset=utf-8
19:59 $ curl -X GET http://localhost:8080/coupons?value=5 -s | jq .
[
  {
    "ID": 23,
    "CreatedAt": "2019-01-02T19:57:57.397605Z",
    "UpdatedAt": "2019-01-02T19:57:57.397605Z",
    "DeletedAt": null,
    "name": "CouponName",
    "brand": "CouponBrand",
    "value": 5,
    "expiry": "2020-01-01T23:59:59Z"
  },
  {
    "ID": 24,
    "CreatedAt": "2019-01-02T19:58:16.986643Z",
    "UpdatedAt": "2019-01-02T19:58:16.986643Z",
    "DeletedAt": null,
    "name": "CouponName2",
    "brand": "CouponBrand2",
    "value": 5,
    "expiry": "2020-01-01T23:59:59Z"
  },
  {
    "ID": 25,
    "CreatedAt": "2019-01-02T19:58:17.617733Z",
    "UpdatedAt": "2019-01-02T19:58:17.617733Z",
    "DeletedAt": null,
    "name": "CouponName2",
    "brand": "CouponBrand2",
    "value": 5,
    "expiry": "2020-01-01T23:59:59Z"
  },
  {
    "ID": 26,
    "CreatedAt": "2019-01-02T19:58:18.695968Z",
    "UpdatedAt": "2019-01-02T19:58:18.695968Z",
    "DeletedAt": null,
    "name": "CouponName2",
    "brand": "CouponBrand2",
    "value": 5,
    "expiry": "2020-01-01T23:59:59Z"
  }
]
```
---