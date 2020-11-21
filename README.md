# Customer Profile System

Customer Profile System API

## Installation

Download [Go](https://golang.org/doc/install) and check GO version

```bash
go version
```

## Details
Customer Profile System contains customer information which are 
- id
- name
- email
- status

To run with your own Database environment, please use the following command
```go
set "DATABASE_URL=<YOUR_DATABASE_URL>" go run server.go
```
And for testing the API integration, run the following command parallel with server.go
```go
go test -v server_test.go
```

This file contains following API
To get all customer information or individual customer information
```api
GET http:localhost:2009/customers 
GET http:localhost:2009/customers/<customer id>
```
To create new customer
```api
POST http:localhost:2009/customers
```
with the example body
```json
{
    "name": "ekkaffwat",
    "email": "ekkawat@gmail.com",
    "status": "inactive"
}
```
To update specific customer information
```api
PUT http:localhost:2009/customers/<customer id>
```
with the example body
```json
{
    "name": "newname",
    "email": "newEmail@gmail.com",
    "status": "inactive"
}
```

To delete customer according to their ID
```python
DELETE http:localhost:2009/customers/<customer id>
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

Worranitta Kraisittipong 21/11/2020
