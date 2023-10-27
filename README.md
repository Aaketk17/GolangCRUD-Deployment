### To run the application
- Create an PostgreSQL DB (locally or using AWS RDS)
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
