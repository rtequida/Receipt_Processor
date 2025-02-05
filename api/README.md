
# OpenAPI codegen Server Creation

**Install oapi-codegen** - Run this command to install the library. I used this library to auto-generate the server-side code of the API from the openapi.yml file and is only necessary if you want to update the API by updating the openapi.yml and running the command below to generate new server-side code.
```
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```
**Generate the server-side code from the openapi.yml file** - Run this command so codegen can auto-generate the server-side of the API using the defined API in the .yml file.
```
oapi-codegen -package api -generate types,server -o api/generatedserver.go api/openapi.yml
```
- oapi-codegen -package api: Tells it to use the codegen library and the api package
- \-generate types,server: Determines which types to generate. We want codegen to generate the server-side code so we can route requests to the handlers.
- \-o api/generatedserver.go: Specifies where to put the output of this command. In this case I am storing the output in a file called generatedserver.go which is located in a directory called api.
- api/openapi.yml: Determines the input file for codegen. In this case we want the openapi.yml file located in the api directory.
