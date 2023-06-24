# Benchmark

## Goal

We want to find out which database is best for a terminology server implementation.

Things to benchmark:

- Different architecture and indexes
- Fastest reads for code lookups (ex. : CodeSystem/$lookup)
- Fastest reads when searching for text matches (ex. : CodeSystem/$find-matches)
- Could we run this fully in-memory? How much memory do we need?

# NoSQL Databases

Which is best for our use case? Which architecture? Which indexes?
How does NoSQL perform? Could we store FHIR documents and avoid marshalling/unmarshalling?

## Aerospike

## Scyladb

## Couchbase

## Mongo

## Cassandra

## Redis

# SQL Databases

How does relation SQL database peform? Some code system use relations (ex. SNOMED CT)

## Postgres

## MySQL

## sqlite

## Aurora

## mssql

## Oracle

# Further analysis

Could we have multiple database types? For example, use SQL database for SNOMED CT and NoSQL for other code systems?

# Other things to consider

Can we use a graph database? How does it perform?
Can we use indexed search engine (ex. Solr/Lucene/ElasticSearch)? How does it perform?
