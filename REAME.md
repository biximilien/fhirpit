# FHIRpit

FHIRpit is a FHIR termilogy server that implements CodeSystem $loopup and $find-matches operations. It is meant to be low latency and high throughput. It is written in Go and uses Google FHIR proto definitions.

# Usage

Import the data into the database. Ensure that SNOMED CT data files are in the data directory.

```sh
$ ./server --import
```

Run the server.

```sh
$ ./server
```
