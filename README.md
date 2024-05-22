# Chirpy an alternative to X

## About
Building this guided backend of a webserver served the purpose to learn how http-networking works in GOLANG.
This small server is structred as a standard REST-API.

## Installation
Only usable by cloning the repository and running the backend server on your localhost on port 8080.
You can play around with the GET requests in the browser, but for POST,PUT and DELETE requests you need an http-client.
I use the Thunder CLient extension on VSCode for that.


You also need to create a .env file where you store the key you sign your JWT tokens with and also the API-Key for authentication for the webhook between the API and the simulated 3rd-party payment provider. 

## Dependencies
This projects uses 3 external packages:
+ For encrypting the password that is stored in the database.
+ For creating the JWT-Tokens.
+ For reading environment variables out of a .env file.


You can install them with these commands on your machine:

```
 go get golang.org/x/crypto v0.23.0

 go get github.com/joho/godotenv v1.5.1

 go get github.com/golang-jwt/jwt/v5 v5.2.1
```

## Possible actions
+ Create a user
+ Login to your account
+ Post chirps (small text snippets)
+ Delete only the chirps you made.
+ Update the users e-mail and password.
+ Using access and refresh tokens
+ Using a webhook for simulated 3rd-party payment confirmation

## Endpoints
All data transfer is done through JSON format.

### /api/users
#### POST Request
```
{
    "email": "example@gmail.com"
    "password" : 12345
}
```
##### Response Body
```
{
  "email": "walt@breakingbad.com",
  "id": 1,
  "is_chirpy_red": false,
  
} 
```
#### PUT Request
Have to be logged in to change password or e-mail address of a user.
```
{
    "email": "example@gmail.com",
    "password" : 12345
}
```
##### Response Body
```
{
  "email": "example@gmail.com",
  "id": 1,
  "is_chirpy_red": false,
  
}
```
### /api/login
#### POST Request
```
{ 
    "email": "example@gmail.com",
    "password" : 12345
}
```
##### Response Body
```
{
  "email": "example@gmail.com",
  "id": 1,
  "is_chirpy_red": false,
  "refresh_token":"example refresh_tokentoken" ,
  "token": "example token"
}
```

### /api/refresh
Refreshes the access token with a longer lived refresh token.

#### POST Request
Put the refresh token `Bearer {refresh_token}`  in the authorization header of the request. 

### /api/revoke 
Revokes a current active refresh token.

#### POST Request
Put the refresh token `Bearer {refresh_token}`  in the authorization header of the request. 

### /api/chirps

#### POST Request
Posts a chirp when logged in

```
 {
    "body": "example_chirp"
 }
 ```
 ##### Response Body
 ```
 {
  "author_id": 1,
  "body": "example_chirp",
  "id": 1
}
 ```
#### GET Request
Returns a slice of all chirps present in the database.

##### Response Body
```
[
    {
    "author_id": 1,
    "body": "example_chirp1",
    "id": 1
    },
    {
    "author_id": 2,
    "body": "example_chirp2",
    "id": 2
    },

]
```
### api/chirps?author_id={example_authorId}&sort={example_sort_order}

#### GET Request
By adding the authors ID as a query parameter we can get a slice of the authors tweets. By setting the sort parameter to `asc` we sort them in ascending chirp-ID order by setting the parameter to `desc` we are sorting them in descending chirp-ID order.
Sorting defaults to ascending when no query parameter is provided.

##### Response Body
```
[
    {
    "author_id": 1,
    "body": "example_chirp1",
    "id": 1
    },
    {
    "author_id": 1,
    "body": "example_chirp2",
    "id": 2
    },

]
```
### /api/chirps/{example_ID_int}

#### GET Request
Returns a specific chirp with the provided ID.

##### Response Body
{
  "author_id": 1,
  "body": "example_chirp",
  "id": 1
}
 #### DELETE Request
 Only possible when logged in.
 Put the access token `Bearer {token}` in the authorization header of the request. 
 Deletes the chirp with the provided ID.







