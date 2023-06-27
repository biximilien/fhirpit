# Aerospike

## Database Planning

### 1. create namespace terminology

This is doing by adding the namespace to the configuration file.

### 2. create set snomed_description

```sql
INSERT INTO terminology.snomed_description (PK, effectiveTime, active, moduleId, conceptId, languageCode, typeId, term, caseSignificanceId) VALUES (1, "20201212", "1", "asd", "asd", "asd", "asd", "asd", "asd")
```

Can this be done by the client code?

### 3. create bins for snomed_description

They were created by the above query

Can this be done by the client code?

### 4. create secondary index on snomed_description

```sql
manage sindex create string cid ns terminology set snomed_description bin cid
```

Can this be done by the client code?
