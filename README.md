# FHIRpit

FHIRpit is a FHIR terminology server that implements CodeSystem $loopup and $find-matches operations. It is meant to be low latency and high throughput. It is written in Go and uses Google FHIR proto definitions.

# Usage

Import the data into the database. Ensure that SNOMED CT data files are in the data directory.

```sh
$ ./server --import
```

Run the server.

```sh
$ ./server
```

# Supported Operations

## $lookup

```sh
$ curl -X POST -H "Content-Type: application/fhir+json" -d '{"resourceType": "Parameters", "parameter": [{"name": "code", "valueCode": "123"}]}' http://localhost:8080/CodeSystem/$lookup
```

## $find-matches

```sh
$ curl -X POST -H "Content-Type: application/fhir+json" -d '{"resourceType": "Parameters", "parameter": [{"name": "system", "valueUri": "http://snomed.info/sct"}, {"name": "code", "valueCode": "123"}]}' http://localhost:8080/CodeSystem/$find-matches
```

# Performance

## $lookup

```sh
$ ab -n 100000 -c 100 -T "application/fhir+json" "localhost:8080/CodeSystem/\$lookup?code=100922009&system=http://snomed.info/sct"
```

## $find-matches

```sh
$ ab -n 100000 -c 100 -T "application/fhir+json" "localhost:8080/CodeSystem/\$find-matches?code=100922009&system=http://snomed.info/sct"
```
