QR-DB
---

A simple project for leisure time to play around with GoLang. The application allows associating a generic string to string map with a random code that can be rendered as a QR code. From this code, it's possible to trace back to the initial map. The application provides a simple REST API interface.

This is a work in progress, and considering the fact I do not really have free time, it will probably not evolve more than that.

Run the application
===

```
go run main.go
```

Create a resource associated to a QR code
===

Request
```
curl -X POST localhost:8080/qr --data {"a": "1", "b": "B", "c": "false"}
```

Response
```
{"success":true,"error":"","code":"NFSHE1IZ"}
```

Get data associated to a QR code
===
Request
```
curl -X GET localhost:8080/qr/NFSHE1IZ
```

Response
```
{"success":true,"error":"","data":{"a":"1","b":"B","c":"false"}}
```

Render a QR code as PNG without page numbers
===

```
curl -X GET localhost:8080/qr/render/NFSHE1IZ
```

Render a QR code as PNG with page numbers
===

```
curl -X GET localhost:8080/qr/render/NFSHE1IZ/page/2/of/4
```
