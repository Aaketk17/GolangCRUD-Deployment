### To run the application
- Create an PostgreSQL DB (locally or using AWS RDS)
- Create users table with columns ```user_id int pk```, ```user_type character varying```, ```name character varying```, ```email character varying```, ```password character varying```, ```phone character varying``` here user_id should be auto incremental

```
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    -- Other columns in users table
);

```
- Create invalidTokens table with columns ```id int pk```, ```token character varying``` here id should be auto incremental

```
  CREATE TABLE invalidtokens (
    id SERIAL PRIMARY KEY,
    -- Other columns in your table
  );
```
- Create .env file in the root directory
- ``` PORT=8080
  DB_URL=xxxxxxxxxx
  DB_USERNAME=xxxxxxxxxx
  DB_PASSWORD=xxxxxxxxxx
  DB_NAME=xxxxxxxxxx
  DB_PORT=xxxxxxxxxx
  JWT_SECRET=xxxxxxxxxx
  ```
  add these values to .env file and populate it with the respective values

- Run the commmand ```go run .```
