# CRATE

## 表结构 data structure

```sql
id int8 NOT NULL
state jsonb NOT NULL
...
```

## API

```text
GET /cyclone-api/db-schema
```

```text
GET /cyclone-api/db-table?schema=
```

```text
GET /cyclone-api/:schema/:table
```

```text
GET /cyclone-api/:schema/:table/:uuid/:id
```

```text
POST /cyclone-api/:schema/:table
```

```text
PUT /cyclone-api/:schema/:table/:uuid/:id
```

```text
DELETE /cyclone-api/:schema/:table/:uuid/:id
```
