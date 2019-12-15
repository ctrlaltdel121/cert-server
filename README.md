# Certificate Server

This server grants self-signed certificates based on user input. It stores
the certs for later retrival, and they can also be deleted. File system storage is supported, with stubs for S3 support in the future

## Testing
You can test the code with `go test ./...`

## Building Docker
`docker build .` will build a docker image that can run the server.

`docker run -P <imgID>` will then run the server on a random port.


## Usage example

### Creating a cert

You can create a cert via a command like the following:

```
curl -XPOST http://localhost/certificates -d '{"names":["test.domain.com"]}'
```

Accepted parameters are:
- names (array)
- valid_from (time)
- valid_to (time)
- is_ca (bool)
- organization_name (string)
- country_name (string)
- state_name (string)
- locality_name (string)
- organizational_unit (string)
- common_name (string)
- email_address (string)

The response will be a JSON object with three parameters:

- The Base64-encoded certificate in the "Cert" field
- The Base64-encoded key in the "Key" field. THIS IS ONLY AVAILABLE WHEN CREATED. IT IS NOT STORED.
- The Serial number of the certificate, used to look it up or delete it from this API.

### Getting the certificate

You can get a cert via a command like the following:

```
curl http://localhost/certificates/<serial>
```

This will return the Cert and Serial fields, but not Key. Key is not stored.

### Deleting the cert

You can delete a cert via a command like the following:

```
curl -XDELETE http://localhost/certificates/<serial>
```

This will delete the Cert from the server.