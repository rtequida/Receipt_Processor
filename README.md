
# Receipt_Processor
This is a RESTful API developed in Go which takes in receipts, generates IDs for each receipt and stores them in-memory, then calculates points based on a specific criteria for any given receipt.  

Endpoints included in this API are:
- /receipts/process
- /receipts/{id}/points

## Getting Started
### Prerequisites
- [Go](https://go.dev/doc/install) - The language used
- [Postman](https://www.postman.com/downloads/) - optional if you want to use an API testing suite

### Additional Libraries Used
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - used to generate the server-side code for the API. This is only necessary if you want to make updates to the API through the openapi.yml file and need to run codegen again to regenerate the server-side code. See [api readme](api/README.md)

### Deploying
Run the project ```go run .``` in the root directory.

From here you can send payloads using an API testing suite such as Postman. I've included a Postman collection you can import to easily setup the calls and have a payload ready to go.
![postman screenshot](Images/postman_screenshot.png)
Or use a curl command
```
curl -X GET http://localhost:8080/receipts/{id}/points
```
Breakdown of API endpoints and how to hit them can be found in this [section](#api-endpoint-breakdown)

### Running Tests
I've provided unit tests that test the functions within the ```app``` and ```handlers``` packages.
To run all tests 
```
go test ./...
```
To run a specific test, simply by specifying the package you want to test
```
go test ./app

output:
?       github.com/rtequida/Receipt_Processor   [no test files]
?       github.com/rtequida/Receipt_Processor/api       [no test files]
ok      github.com/rtequida/Receipt_Processor/app       (cached)
ok      github.com/rtequida/Receipt_Processor/handlers  0.210s
```
You can also include the verbose tag ```-v``` to get a more detailed breakdown of the tests
```
go test ./app -v

partial output:
=== RUN   TestGetPoints
--- PASS: TestGetPoints (0.00s)
=== RUN   TestValidateReceipt
--- PASS: TestValidateReceipt (0.00s)
=== RUN   TestValidateID
--- PASS: TestValidateID (0.00s)
PASS
ok      github.com/rtequida/Receipt_Processor/app       0.163s
```
### API Endpoint Breakdown
#### Process Receipts
- Path: /receipts/process
- Method: POST
- Sample command:
```
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{"retailer": "M&M Corner Market", "purchaseDate": "2022-03-20", "purchaseTime": "14:33", "items":[{"shortDescription": "Gatorade", "price": "2.25"}, {"shortDescription": "Gatorade", "price": "2.25"}, {"shortDescription": "Gatorade", "price": "2.25"}, {"shortDescription": "Gatorade", "price": "2.25"}], "total": "9.00"}' http://localhost:8080/receipts/process
```
Status Code:  ```200 OK```
Response: ```{"id":"7f14946f-05ca-4302-ad42-817bb63e4df2"}```

#### Get Points
- Path: /receipts/{id}/points
- Method: GET
- Sample Command
```
curl -X GET http://localhost:8080/receipts/7f14946f-05ca-4302-ad42-817bb63e4df2/points
```
Status Code: ```200 OK```
Response: ```{"points": 109}```
