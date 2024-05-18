# Url Shortener

---
This project is a very simple implementation of a URL shortener service.
The service is a RESTful service the uses cassandra as db.
---
The hashing algorithm used to shorten the given url is fnv32. Been picked mainly for the lack of security constrains,
it's uniform distribution and fast encoding technique.
After hashing the url to an unsigned integer of 32 bit, 
the hashing is then transposed into a base62 string that is used as the short-code.
---
The REST api exposed the pattern `/urls/` on port `8080` with the HTTP methods:
* `POST` (create)
* `GET` (read)
* `PUT` (update)
* `DELETE` (delete)

---
Tested with dockerized cassandra instance `docker run --rm --name cassandra-db -d -p 9042:9042 cassandra`