# Zeth

A web based interface to manage your ethereum node(s).

## Run project

### Development

Run frontend:

```sh
cd app
npm run dev
```

Run backend:

```sh
make run
```

View UI at [http://localhost:7000](http://localhost:7000)

---

### SurrealDB setup

Run local DB:

```sh
docker run --rm -it --name surrealdb -p 8000:8000 surrealdb/surrealdb:latest start --log trace --user admin --pass admin memory
```

#### Setup DB

source: https://gist.github.com/koakh/fbbc37cde630bedcf57acfd4d6a6956b

**INFO:**

```sh
curl -X POST \
  --header "Accept: application/json" \
  --header "NS: test" \
  --header "DB: test" \
  --user "admin:admin" \
  --data "INFO FOR DB;" \
  http://localhost:8000/sql | jq .
```

**Setup user table:**

```sh
DATA="
BEGIN TRANSACTION;
DEFINE TABLE user SCHEMAFULL
  PERMISSIONS 
    FOR select, update WHERE id = \$auth.id, 
    FOR create, delete NONE;
DEFINE FIELD user ON user TYPE string;
DEFINE FIELD pass ON user TYPE string;
DEFINE FIELD tags ON user TYPE array;
DEFINE INDEX idx_user ON user COLUMNS user UNIQUE;
COMMIT TRANSACTION;
"
curl -X POST \
  --header "Accept: application/json" \
  --header "NS: test" \
  --header "DB: test" \
  --user "admin:admin" \
  --data "${DATA}" \
  http://localhost:8000/sql
```

**Setup scope:**

```sh
DATA="
BEGIN TRANSACTION;
DEFINE SCOPE allusers
  -- the JWT session will be valid for 14 days
  SESSION 14d
  -- The optional SIGNUP clause will be run when calling the signup method for this scope
  -- It is designed to create or add a new record to the database.
  -- If set, it needs to return a record or a record id
  -- The variables can be passed in to the signin method
  SIGNUP ( CREATE user SET user = string::lowercase(string::trim(\$user)), pass = crypto::argon2::generate(\$pass), tags = \$tags )
  -- The optional SIGNIN clause will be run when calling the signin method for this scope
  -- It is designed to check if a record exists in the database.
  -- If set, it needs to return a record or a record id
  -- The variables can be passed in to the signin method
  SIGNIN ( SELECT * FROM user WHERE user = string::lowercase(string::trim(\$user)) AND crypto::argon2::compare(pass, \$pass) )
  -- this optional clause will be run when calling the signup method for this scope
COMMIT TRANSACTION;
"
curl -X POST \
  --header "Accept: application/json" \
  --header "NS: test" \
  --header "DB: test" \
  --user "admin:admin" \
  --data "${DATA}" \
  http://localhost:8000/sql
```

**Setup endpoints table:**

```sh
DATA="
BEGIN TRANSACTION;
DEFINE TABLE endpoints SCHEMAFULL
  PERMISSIONS 
    FOR select, update, delete WHERE id = \$auth.id,
    FOR create FULL;
DEFINE FIELD name ON endpoints TYPE string;
DEFINE FIELD enabled ON endpoints TYPE bool;
DEFINE FIELD date_added ON endpoints TYPE datetime;
DEFINE FIELD rpc_url ON endpoints TYPE string;
DEFINE FIELD type ON endpoints TYPE string;
DEFINE INDEX idx_endpoints_name ON endpoints COLUMNS name UNIQUE;
DEFINE INDEX idx_endpoints_rpc_url ON endpoints COLUMNS rpc_url UNIQUE;
-- TODO: define event for creation
-- TODO: define event for update
-- TODO: define event for deletion
COMMIT TRANSACTION;
"
curl -X POST \
  --header "Accept: application/json" \
  --header "NS: test" \
  --header "DB: test" \
  --user "admin:admin" \
  --data "${DATA}" \
  http://localhost:8000/sql
```

---
