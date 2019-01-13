# REST_game
This is a RESTful web service that supports GRUD operations on
a postgreSQL database storing basic video game information. It
will respond to GET requests with one or more video game titles along
with their developer and rating in a JSON format.

## Database Set Up
The postgreSQL database that this application expects to interact
with should have a table that stores a video games title, developer,
and rating. 

```
CREATE TABLE games (title VARCHAR(100), developer VARCHAR(100), rating CHAR(1));
```

## Logging Into Database
The first and only command line argument should be a JSON file that will store
the database credentials. The application will read this JSON file and attempt
to connect to the specified database. 

### JSON File Format Example
```
{
	"host":"...",
	"port":1234,
	"user":"...",
	"dbname":"..."
}
```

### Starting the Server
```
$ ./REST_game database_login.json
```

## Endpoints
* gameAPI/add             (POST)
* gameAPI/{title}         (GET)
* gameAPI/{title}         (PUT)
* gameAPI/{title}         (DELETE)
* gameAPI/developer/{dev} (GET)
* gameAPI/rating/{rating} (GET)

## JSON Response Format
```
{
	"title":"...",
	"developer":"...",
	"rating":"..."
}
```

## Dependencies
* [Gorilla Mux ](https://github.com/gorilla/mux)
* [lib/pq](https://github.com/lib/pq)










