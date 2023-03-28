# Zeth

A web based interface to manage your ethereum node(s).

## Run project

### Development

Open 3 shells and run the following:

```sh
# run persistent db (via docker)
make db
```

```sh
# run client
make client
```

```sh
# run server
make server
```

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
DEFINE FIELD user ON user TYPE string ASSERT is::email(\$value);
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

**Setup endpoint table:**

```sh
DATA="
BEGIN TRANSACTION;
DEFINE TABLE endpoint SCHEMAFULL
  PERMISSIONS 
    FOR select, create, update, delete WHERE user = \$auth.id;
DEFINE FIELD user ON endpoint TYPE record(user) VALUE \$auth.id;
DEFINE FIELD name ON endpoint TYPE string VALUE string::trim(\$value) ASSERT string::length(\$value) > 2 AND string::length(\$value) < 65;
DEFINE FIELD enabled ON endpoint TYPE bool VALUE \$value or true;
DEFINE FIELD date_added ON endpoint TYPE datetime VALUE time::now();
DEFINE FIELD rpc_url ON endpoint TYPE string ASSERT \$value = /^(http|https|ws|wss):\/\/.+/ AND string::length(\$value) < 129;
DEFINE FIELD proxy_url ON endpoint TYPE string VALUE string::concat(session::origin()', '/', \$this.id, '/rpc');
DEFINE FIELD symbol ON endpoint TYPE string ASSERT string::length(\$value) > 2 AND string::length(\$value) < 13;
DEFINE FIELD block_explorer_url ON endpoint TYPE string;
DEFINE INDEX idx_endpoint_name ON endpoint COLUMNS user, name UNIQUE;
DEFINE INDEX idx_endpoint_rpc_url ON endpoint COLUMNS user, rpc_url UNIQUE;
# DEFINE EVENT generate_proxy_url
#   ON TABLE endpoint
#   WHEN \$event='CREATE'
#   THEN (UPDATE endpoint SET proxy_url = string::concat(session::origin(), '/', \$after.id, '/rpc'));
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

```sh
curl -X POST \
  --header "Accept: application/json" \
  --header "NS: test" \
  --header "DB: test" \
  --user "admin:admin" \
  --data "remove table endpoint;" \
  http://localhost:8000/sql
```

```sh
curl -X POST \
  --header "Accept: application/json" \
  --header "NS: test" \
  --header "DB: test" \
	--header "SC: allusers" \
  --user "786zshan@gmail.com:admin" \
  --data "select * from endpoint;" \
  http://localhost:8000/sql
```

