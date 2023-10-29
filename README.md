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

### To run using docker
- Pull the postgres docker image ```docker pull postgres```
- create a docker newtork ```docker network create golangdemo-network```
- Run the postgres docker image ```docker run --name pgsql --network golangdemo-network  -e POSTGRES_PASSWORD=golangDemo -p 5432:5432 -d postgres```
- Create the users and invalidTokens as given above
- Pull the application docker image ```docker pull aaketk/golang-demo```
- Run the application docker image 
```
docker run --network golangdemo-network --name golangdemo -e PORT=8080 -e DB_URL=pgsql -e DB_USERNAME=postgres -e DB_PASSWORD=golangDemo -e DB_NAME=golangpg -e DB_PORT=5432 -e JWT_SECRET=V3S2AymxpNjn3qUQ3EMnEDc4JnutxT3aPnwNKEPfoppnhNucGjzdKJ67A45zW4HTmrfxq7dcHdhPu8okz74ap2mk8gGdRfmSX9HHtZunoPuFzF38nNQts83koC3qKRUn -p 8080:8080 -d aaketk/golang-demo:latest
```
Modify the env variables as you required.
