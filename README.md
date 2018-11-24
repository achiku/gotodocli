## gotodocli

Go Conference JP 2018 Sample App

##### PostgreSQL

```sql
-- psql -U admin -d template1
CREATE database gotodo;
CREATE USER gotodo_root;
ALTER USER gotodo_root WITH SUPERUSER;
```

```sql
-- psql -U gotodo_root -d gotodoit
CREATE USER gotodo_api;
CREATE USER gotodo_api_test;
CREATE SCHEMA gotodo_api AUTHORIZATION gotodo_api;
```
