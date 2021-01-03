# Analysis API

REST-powered API for analyzing media and text.

## Endpoints

- Media analysis `<bind-address>/analysis/{reqType}/{type}/{id}`
- Text analysis `<bind-address>/analysis/{reqType}`
- Health check `<bind-address>/healthy`

## Request types

- Complete media analysis: `complete`
- Media fingerprinting: `fingerprinting`
- Media NSFW recognition: `nsfw`
- Gibberish text analysis: `gibberish`

## Returned data

The data returned by the server when analyzing media is in the form:

```go
type Analysis struct {
    Fingerprint            FingerprintResponse
    NSFW                   NSFWResponse
    FingerprintErrorString string
    NSFWErrorString        string
}

```

The data returned by the server when analyzing text is in the form:

```go
type GibberishResponse struct {
	IsGibberish bool
}
```

## Environment options

- API bind address and port: `API_BIND_ADDRESS` (defaults to `localhost:9999`).
- Path to configuration file: `API_CFG_PATH` (defaults to `config.toml`).

## Configuration file structure

```toml
fingerprintendpoint = "http://localhost:10000/fingerprinting"
nsfwendpoint = "http://localhost:10001/nsfw"
gibberishendpoint   = "http://localhost:10002/gibberish"
testing = <testing-status>

[database]
  authsource = "admin"
  collectionname = "tokens"
  databasename = "fpserver"
  hosts = ["localhost:27017"]
  password = "fpserver"
  replicasetname = "shitposting"
  username = "fpserver"

[redis]
    Address = "localhost:6379"
    NSFWDatabase = 2
    FingerprintDatabase = 1
```
