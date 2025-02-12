# Digital Shelf

Digital Shelf is a RESTful API backend that is used to inventory items on shelves. The idea being that it'll help you know what books, movies, music, and TV shows that you have on your shelves. It will also tell you where each item is located. This allows for keeping track of the items in your personal library/collection. I am also planning on adding a wishlist feature, so that you'll know what items you would like to get for your collection. 

This project is still a work in progress and currently only supports movies, but support for other items is in the works.

This project was designed for personal use, but has a good foundation that supports many users and locations. The same backend server could be used to track library items across a great many locations with a large number of users. The backend is built in Go and uses a PostgreSQL database, and is very fast.

![code tests badge](https://github.com/rodabaugh/digitalshelf/actions/workflows/ci.yml/badge.svg)

# Frontend/Client

This is just the backend of the application, but I have also written a CLI frontend for the application. I would like to create a webapp interface or a mobile app for this at some point, but haven't gotten to that point yet. If you would like to create your own frontend, please feel free to do so. Documentation for the REST API can be found below.

The CLI interface I wrote for this application can be found at: <https://github.com/Rodabaugh/digitalshelf-cli/>

# Manual Setup

## Prerequisites

To run this project, you will need:
1. A PostgreSQL database (this can be on the same server as the backend, if you would like)
2. A backend server with Go installed.

## Get the repo

1. Clone the repo `git clone https://github.com/Rodabaugh/digitalshelf/`
2. Navigate to the program dir `cd digitalshelf`

## Configuration

Environment variables are used for configuration. Create a .env file in the root of the project dir. You need to specify your DB_URL, PLATFORM, and JWT_SECRET. DB_URL is the url for your database. Platform can either be "prod" or "dev". You can generate your JWT_SECRET with `openssl rand -base64 64`. Your `.env` file should look something like the one below. Please be sure to create your own JWT_SECRET and use your own DB_URL.
```
PLATFORM=prod
DB_URL="postgres://postgresUser:postgresPass@localhost:5432/digitalshelf?sslmode=disable"
JWT_SECRET="ALItvAPa64TLZ4wjqWsaiVW3ZrQ7ZT209sAkIsos8K3p6ldeMb+K5Ji5j90kI4cQ
k0I6WY6KgXALHP7EjeLXOw=="
```

A port may also be specified using ```PORT=1234```. If a port is not specified, it will default to 8080.

## Setting up the database

Goose is used to manage the database migrations. Install goose with `go install github.com/pressly/goose/v3/cmd/goose@latest`

Navigate to the sql/schema dir `cd sql/schema`

Setup the database using goose `goose postgres <connection_string> up` e.g `goose postgres postgres://postgresUser:postgresPass@localhost:5432/digitalshelf?sslmode=disable up`

## Compile and run the backend

Once your .env has been configured, and your database is setup, it is time to build and run the backend.

Build the application with `go build`

Run the backend application with `./digitalshelf`

Once the backend server is running, you can setup your server to run the backend application as a service. 

## Success

At this point, the DigitalShelf backend should be running on your server. From here, you can use the CLI frontend (<https://github.com/Rodabaugh/digitalshelf-cli/>) to interact with the backend, or write your own frontend.

# Docker

There is also a Docker image for the backend. You will still need to setup the database using Goose, and the same environment variables should be specified. Environment variables can either be specified inline with the docker run command, or you can import them using the `--env-file` argument.

Example: `docker run -p 8080:8080 --env-file ./.env digitalshelf`

The docker image can be found here at <https://hub.docker.com/repository/docker/rodabaugh/digitalshelf/general> and can be pulled using `docker pull rodabaugh/digitalshelf:latest`

# API Endpoints

## Design Overview

Items (such as movies) are stored on shelves. Those shelves are in cases. Those cases are located at a location. The location is owned by a user, and has users as members. This allows you to know what shelf an item is on, which case that shelf is in, and where the case is located.

## Users

### POST /api/users
Request body:
```json
{
    "name": "John Smith",
    "email": "jsmith@example.com",
    "password": "Password123"
}
```

Response body:
```json
{
    "id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "name": "John Smith",
    "email": "jsmith@example.com"
}
```

### PUT /api/users
Used to update email address or password. Uses auth token to identify the user that is being updated.

Requires token in auth headers.

Request body:
```json
{
    "email": "jsmith@example.com",
    "password": "Password123"
}
```

Response body:
```json
{
    "id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "name": "John Smith",
    "email": "jsmith@example.com"
}
```

### GET /api/users

Not documented yet.

### GET /api/users/{user_id}
Used to get user details using the user's ID.

No auth required.

Request body: None

Response body:
```json
{
  "id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
  "name": "John Smith",
  "email": "jsmith@example.com",
  "created_at": "2025-01-18T17:19:06.09633Z",
  "updated_at": "2025-01-19T16:42:08.06933Z"
}
```

### GET /api/search/users?email=
Uses the email URL query to search for users by email address.

No auth required.

Request body: None

Response body:
```json
{
  "id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
  "name": "John Smith",
  "email": "jsmith@example.com",
  "created_at": "2025-01-18T17:19:06.09633Z",
  "updated_at": "2025-01-18T17:19:06.09633Z"
}
```

## Locations

### POST /api/locations

Used to create new locations for the authenticated user.

Auth token is required.

Request body:
```json
{
  "name": "John's House",
  "owner_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3"
}
```

Response body:
```json
{
  "id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
  "name": "John's House",
  "owner_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
  "created_at": "2025-01-26T13:44:17.361433Z",
  "updated_at": "2025-01-26T13:44:17.361433Z"
}
```

### GET /api/locations

Not documented yet.

### GET /api/locations/{location_id}

Used to get information about the location. 

Auth token is required. The user must be a member of the location.

Request body: None

Response body:
```json
{
  "id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
  "name": "John's House",
  "owner_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
  "created_at": "2025-01-26T13:44:17.361433Z",
  "updated_at": "2025-01-26T13:44:17.361433Z"
}
```

### GET /api/users/{user_id}/locations

Used to get the locations the user is a member of.

Auth token is required.

Request body: None

Response body:
```json
[
  {
    "userID": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "location_id": "5722d862-97d8-409c-91e1-3281ff7882aa",
    "location_name": "home",
    "owner_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "joined_at": "2025-01-18T18:15:10.172788Z"
  },
  {
    "userID": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "location_id": "d62960f6-0c11-4b00-a635-d54755002d02",
    "location_name": "work",
    "owner_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "joined_at": "2025-01-26T13:33:22.442817Z"
  }
]
```

### GET /api/search/locations?owner_id=
Uses the owner_id query to search for locations by the owner's user ID.

No auth required.

Request body: None

Response body:
```json
[
  {
    "id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
    "name": "John's House",
    "owner_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "created_at": "2025-01-26T13:44:17.361433Z",
    "updated_at": "2025-01-26T13:44:17.361433Z"
  }
]
```

## Invites

### POST /api/locations/{location_id}/invites

Used to invite a user to a location.

Auth token required. User must be the owner of the location.

Request body:
```json
{
  "user_id":"4b2c6a66-cfc6-4a2b-aad6-6cec9507debe"
}
```

Response body:
```json
{
  "location_id": "5722d862-97d8-409c-91e1-3281ff7882aa",
  "user_id": "4b2c6a66-cfc6-4a2b-aad6-6cec9507debe",
  "invited_at": "2025-01-26T14:15:27.029995988-07:00"
}
```

### GET /api/locations/{location_id}/invites

Returns the invites for the location.

Auth token is required. The user must be the owner of the location.

Request body: None

Response body:
```json
[
  {
    "location_id": "5722d862-97d8-409c-91e1-3281ff7882aa",
    "userID": "4b2c6a66-cfc6-4a2b-aad6-6cec9507debe",
    "user_name": "Bill Smith",
    "user_email": "bsmith@example.com",
    "invited_at": "2025-01-26T14:01:04.140827Z"
  }
]
```

### GET /api/users/{user_id}/invites

Used to get the authenticated user's invites. These are locations that the user can join.

Auth token is required.

Request body: None

Response body:
```json
[
  {
    "userID": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "location_id": "5c94a1cc-8127-469d-bb7c-d891167872d8",
    "location_name": "bills_house",
    "owner_id": "4b2c6a66-cfc6-4a2b-aad6-6cec9507debe",
    "invited_at": "2025-01-26T13:38:46.425538Z"
  },
  {
    "userID": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "location_id": "82ad63af-e615-42b2-8042-7688c88294cb",
    "location_name": "bills_storage",
    "owner_id": "4b2c6a66-cfc6-4a2b-aad6-6cec9507debe",
    "invited_at": "2025-01-26T13:39:12.962544Z"
  }
]
```

### DELETE /api/locations/{location_id}/invites/{user_id}

Removes the invite for location.

Auth token is required. The user must be the invited user or the owner of the location.

Request body: None

Response body: None

## Auth

### POST /api/login

Used to login and get a token and refresh token.

No auth required.

Request body:
```json
{
    "email": "jsmith@example.com",
    "password": "Password123"
}
```

Response body:
```json
{
  "id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
  "name": "John Smith",
  "email": "jsmith@example.com",
  "created_at": "2025-01-18T17:19:06.09633Z",
  "updated_at": "2025-01-19T16:42:08.06933Z",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJkaWdpdGFsc2hlbGYtYWNjZXNzIiwic3ViIjoiZDJkYjc1OGMtYmQ4NC00YzljLTk1YTEtOTNlNjBjNzRjOWMzIiwiZXhwIjoxNzM3OTI2ODg0LCJpYXQiOjE3Mzc5MjMyODR9.w-R67hzUyXxtcpRH8jVrC9KpmgP8rYinhbHspJ_Mmqk",
  "refresh_token": "ae1f9ea8039ef8be825738d9678036ca8bd01d9cc01779577aa5c78a6eeb4312"
}
```

### POST /api/refresh
Used to get a new token, using the refresh token.

Auth token is required. The refresh token should be provided.

Request body: None

Response Body:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJkaWdpdGFsc2hlbGYtYWNjZXNzIiwic3ViIjoiZDJkYjc1OGMtYmQ4NC00YzljLTk1YTEtOTNlNjBjNzRjOWMzIiwiZXhwIjoxNzM3OTMyNzM5LCJpYXQiOjE3Mzc5MjkxMzl9.LZMyMIAkNTmW2jL56lGYUfUWz8lvCGjvY4TIJssyImQ"
}
```

### POST /api/revoke

Not documented yet.

### POST /api/revoke-all

Not documented yet.

## Members

### POST /api/locations/{location_id}/members

Used to add a member to a location.

Auth token required. User must have an invite for the location, or is the owner of the location.

Request body:
```json
{
  "user_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3"
}
```

Response body:
```json
{
  "location_id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
  "user_id": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
  "joined_at": "2025-01-26T13:52:47.281237Z"
}
```

### GET /api/locations/{location_id}/members

Returns the members for the location.

Auth token is required. The user must be the owner of the location.

Request body: None

Response body:
```json
[
  {
    "location_id": "5722d862-97d8-409c-91e1-3281ff7882aa",
    "userID": "d2db758c-bd84-4c9c-95a1-93e60c74c9c3",
    "user_name": "John Smith",
    "user_email": "jsmith@example.com",
    "joined_at": "2025-01-18T18:15:10.172788Z"
  }
]
```

### DELETE /api/locations/{location_id}/members/{user_id}

Removes the member from a location.

Auth token is required. The user must be a member being removed or the owner of the location.

Request body: None

Response body: None

## Cases

### POST /api/cases

Create a case at a location.

Auth token is required. User must be a member of the location.

Request body:
```json
{
  "location_id":"970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
  "name":"New Case"
}
```

Response body:
```json
{
  "id": "205bb035-d6b5-4b8d-9ea9-6b755343a92e",
  "name": "New Case",
  "location_id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
  "created_at": "2025-01-26T15:15:12.895479Z",
  "updated_at": "2025-01-26T15:15:12.895479Z"
}
```

### GET /api/cases
Not documented yet.

### GET /api/cases/{case_id}
Get details about a case.

Auth token is required. User must be a member of the case's location.

Request body: None

Response body:
```json
{
  "id": "205bb035-d6b5-4b8d-9ea9-6b755343a92e",
  "name": "New Case",
  "location_id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
  "created_at": "2025-01-26T15:15:12.895479Z",
  "updated_at": "2025-01-26T15:15:12.895479Z"
}
```

### GET /api/locations/{location_id}/cases

Get the cases at a location.

Auth token is required. User must be a member of the location.

Request body: None

Response body:
```json
[
  {
    "id": "205bb035-d6b5-4b8d-9ea9-6b755343a92e",
    "name": "New Case",
    "location_id": "970ea0e1-9fe2-4b71-a756-3e733f96b6b5",
    "created_at": "2025-01-26T15:15:12.895479Z",
    "updated_at": "2025-01-26T15:15:12.895479Z"
  }
]
```

## Shelves

### POST /api/shelves
Create a shelf in a case.

Auth token is required. User must be a member of the case's location.

Request body:
```json
{
  "name":"New Shelf",
  "case_id":"205bb035-d6b5-4b8d-9ea9-6b755343a92e"
}
```

Response body:
```json
{
  "id": "d41dc884-7f5b-4c3e-b05c-b7f5cbad8d22",
  "name": "New Shelf",
  "case_id": "205bb035-d6b5-4b8d-9ea9-6b755343a92e",
  "created_at": "2025-01-26T15:28:39.873399Z",
  "updated_at": "2025-01-26T15:28:39.873399Z"
}
```

### GET /api/shelves/{shelf_id}
Get details about a shelf.

Auth token is required. User must be a member of the shelf's location.

Request body: None

Response body:
```json
{
  "id": "d41dc884-7f5b-4c3e-b05c-b7f5cbad8d22",
  "name": "New Shelf",
  "case_id": "205bb035-d6b5-4b8d-9ea9-6b755343a92e",
  "created_at": "2025-01-26T15:28:39.873399Z",
  "updated_at": "2025-01-26T15:28:39.873399Z"
}
```

### GET /api/shelves

Not documented yet.

### GET /api/cases/{case_id}/shelves
Get the shelves in a case.

Auth token is required. User must be a member of the case's the location.

Request body: None

Response body:

```json
[
  {
    "id": "d41dc884-7f5b-4c3e-b05c-b7f5cbad8d22",
    "name": "New Shelf",
    "case_id": "205bb035-d6b5-4b8d-9ea9-6b755343a92e",
    "created_at": "2025-01-26T15:28:39.873399Z",
    "updated_at": "2025-01-26T15:28:39.873399Z"
  }
]
```

## Movies

### POST /api/movies
Add a movie to the database. A shelf ID must be provided, as the shelf is where the movie is located.

Auth token is required. The requesting user must be a member of the shelf's location.

Request body:
```json
{
  "title": "Dune: Part Two",
  "genre": "Sci-fi",
  "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
  "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
  "director": "Denis Villeneuve",
  "barcode": "883929802357",
  "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
  "release_date": "2024-03-01T00:00:00Z"
}
```

Response body:
```json
{
  "id": "7b43a93f-34eb-49b4-9396-e48b21697a5f",
  "title": "Dune: Part Two",
  "genre": "Sci-fi",
  "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
  "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
  "director": "Denis Villeneuve",
  "barcode": "883929802357",
  "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
  "release_date": "2024-03-01T00:00:00Z",
  "created_at": "2025-01-26T15:10:22.03059Z",
  "updated_at": "2025-01-26T15:10:22.03059Z"
}
```

### GET /api/movies

Not documented yet.

### GET /api/shelves/{shelf_id}/movies
Gets a list of movies on the shelf.

Auth token is required. The user must be a member of the shelf's location.

Request body: None

Response body:
```json
[
  {
    "id": "97a940ab-bc47-4cd9-861b-f9f9d7e2e333",
    "title": "Dune: Part Two",
    "genre": "Sci-fi",
    "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
    "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
    "director": "Denis Villeneuve",
    "barcode": "883929802357",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2024-03-01T00:00:00Z",
    "created_at": "2025-01-18T17:27:56.484798Z",
    "updated_at": "2025-01-18T17:27:56.484798Z"
  }
]
```

### GET /api/movies/{movie_id}
Provides details about the movie object.

Auth token is required. The user must be a member of the movie's location.

Request body: None

Response body:
```json
{
    "id": "97a940ab-bc47-4cd9-861b-f9f9d7e2e333",
    "title": "Dune: Part Two",
    "genre": "Sci-fi",
    "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
    "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
    "director": "Denis Villeneuve",
    "barcode": "883929802357",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2024-03-01T00:00:00Z",
    "created_at": "2025-01-18T17:27:56.484798Z",
    "updated_at": "2025-01-18T17:27:56.484798Z"
}
```

### GET /api/locations/{location_id}/movies
Gets a list of movies at the location ID.

Auth token is required. The user must be a member or the owner of the location.

Request body: None

Response body:
```json
[
  {
    "id": "97a940ab-bc47-4cd9-861b-f9f9d7e2e333",
    "title": "Dune: Part Two",
    "genre": "Sci-fi",
    "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
    "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
    "director": "Denis Villeneuve",
    "barcode": "883929802357",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2024-03-01T00:00:00Z",
    "created_at": "2025-01-18T17:27:56.484798Z",
    "updated_at": "2025-01-18T17:27:56.484798Z"
  }
]
```

### GET /api/search/movie_barcodes/{barcode}

Searches movies for items matching the barcode.

Request body: None

Response body:
```json
{
    "id": "97a940ab-bc47-4cd9-861b-f9f9d7e2e333",
    "title": "Dune: Part Two",
    "genre": "Sci-fi",
    "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
    "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
    "director": "Denis Villeneuve",
    "barcode": "883929802357",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2024-03-01T00:00:00Z",
    "created_at": "2025-01-18T17:27:56.484798Z",
    "updated_at": "2025-01-18T17:27:56.484798Z"
}
```

### GET /api/search/movie

Searches movies for a search term. Database searches title, actors, genre, writer, and director.

Auth token is required. The user must be a member or the owner of the location.

Request body:
```json
{
  "query":"Dune",
  "location_id":"5722d862-97d8-409c-91e1-3281ff7882aa"
}
```

Response body:
```json
[
  {
      "id": "97a940ab-bc47-4cd9-861b-f9f9d7e2e333",
      "title": "Dune: Part Two",
      "genre": "Sci-fi",
      "actors": "Timothée Chalamet, Zendayam, Rebecca Ferguson",
      "writer": "Denis Villeneuve, Jon Spaihts, Frank Herbert",
      "director": "Denis Villeneuve",
      "barcode": "883929802357",
      "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
      "release_date": "2024-03-01T00:00:00Z",
      "created_at": "2025-01-18T17:27:56.484798Z",
      "updated_at": "2025-01-18T17:27:56.484798Z"
  }
]
```

### PUT /api/movies/{movie_id}
Not documented yet. 

## Shows

### POST /api/shows
Add a show to the database. A shelf ID must be provided, as the shelf is where the show is located.

Auth token is required. The requesting user must be a member of the shelf's location.

Request body:
```json
{
  "title": "Person of Interest",
  "season": 2,
  "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
  "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
  "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
  "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
  "barcode": "883929278596",
  "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
  "release_date": "2013-05-09T00:00:00Z"
}
```

Response body:
```json
{
  "id": "fc3bece2-5810-4176-ac4f-b5ecbb50d1f0",
  "title": "Person of Interest",
  "season": 2,
  "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
  "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
  "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
  "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
  "barcode": "883929278596",
  "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
  "release_date": "2013-05-09T00:00:00Z",
  "created_at": "2025-01-26T15:10:22.03059Z",
  "updated_at": "2025-01-26T15:10:22.03059Z"
}
```

### GET /api/shows

Not documented yet.

### GET /api/shelves/{shelf_id}/shows
Gets a list of shows on the shelf.

Auth token is required. The user must be a member of the shelf's location.

Request body: None

Response body:
```json
[
  {
    "id": "fc3bece2-5810-4176-ac4f-b5ecbb50d1f0",
    "title": "Person of Interest",
    "season": 2,
    "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
    "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
    "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
    "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
    "barcode": "883929278596",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2013-05-09T00:00:00Z",
    "created_at": "2025-01-26T15:10:22.03059Z",
    "updated_at": "2025-01-26T15:10:22.03059Z"
  }
]
```

### GET /api/shows/{show_id}
Provides details about the show object.

Auth token is required. The user must be a member of the show's location.

Request body: None

Response body:
```json
{
  "id": "fc3bece2-5810-4176-ac4f-b5ecbb50d1f0",
  "title": "Person of Interest",
  "season": 2,
  "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
  "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
  "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
  "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
  "barcode": "883929278596",
  "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
  "release_date": "2013-05-09T00:00:00Z",
  "created_at": "2025-01-26T15:10:22.03059Z",
  "updated_at": "2025-01-26T15:10:22.03059Z"
}
```

### GET /api/locations/{location_id}/shows
Gets a list of shows at the location ID.

Auth token is required. The user must be a member or the owner of the location.

Request body: None

Response body:
```json
[
  {
    "id": "fc3bece2-5810-4176-ac4f-b5ecbb50d1f0",
    "title": "Person of Interest",
    "season": 2,
    "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
    "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
    "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
    "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
    "barcode": "883929278596",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2013-05-09T00:00:00Z",
    "created_at": "2025-01-26T15:10:22.03059Z",
    "updated_at": "2025-01-26T15:10:22.03059Z"
  }
]
```

### GET /api/search/show_barcodes/{barcode}

Searches shows for items matching the barcode.

Request body: None

Response body:
```json
{
  "id": "fc3bece2-5810-4176-ac4f-b5ecbb50d1f0",
  "title": "Person of Interest",
  "season": 2,
  "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
  "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
  "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
  "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
  "barcode": "883929278596",
  "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
  "release_date": "2013-05-09T00:00:00Z",
  "created_at": "2025-01-26T15:10:22.03059Z",
  "updated_at": "2025-01-26T15:10:22.03059Z"
}
```

### GET /api/search/show

Searches shows for a search term. Database searches title, actors, genre, writer, and director.

Auth token is required. The user must be a member or the owner of the location.

Request body:
```json
{
  "query":"Action",
  "location_id":"5722d862-97d8-409c-91e1-3281ff7882aa"
}
```

Response body:
```json
[
  {
    "id": "fc3bece2-5810-4176-ac4f-b5ecbb50d1f0",
    "title": "Person of Interest",
    "season": 2,
    "genre": "Action, Crime, Drama, Mystery, Sci-Fi, Thriller",
    "actors": "Jim Caviezel, Taraji P. Henson, Kevin Chapman, Michael Emerson",
    "writer": "Jonathan Nolan, Denise Thé, Sean Hennen, Erik Mountain",
    "director": "Richard J. Lewis, Jon Cassar, Jeffrey Hunt, James Whitmore Jr., Félix Alcalá, Frederick E. O. Toye, Helen Shaver, Clark Johnson, Stephen Surjik, Chris Fisher, John Dahl, Jonathan Nolan, Kenneth Fink, Tricia Brock",
    "barcode": "883929278596",
    "shelf_id": "86a210c7-2c90-4c64-b481-9059b4b376db",
    "release_date": "2013-05-09T00:00:00Z",
    "created_at": "2025-01-26T15:10:22.03059Z",
    "updated_at": "2025-01-26T15:10:22.03059Z"
  }
]
```