# echo

echo request info

```
$ go run main.go 9001
Serving echo HTTPServer on :9001 ...

POST /1234?a=2&b=2 HTTP/1.1
User-Agent: curl/7.35.0
Accept: */*
Content-Length: 3
Content-Type: application/x-www-form-urlencoded

a=2
```

```
$ curl "http://127.0.0.1:9001/1234?a=2&b=2#3" -d a=2
POST /1234?a=2&b=2 HTTP/1.1
User-Agent: curl/7.35.0
Accept: */*
Content-Length: 3
Content-Type: application/x-www-form-urlencoded

a=2
```
