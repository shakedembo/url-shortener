# Url Shortener
---

This project is a very simple implementation of a URL shortener service.
The service is a RESTful service that uses Cassandra as db.

---

The hashing algorithm used to shorten the given URL is fnv32.
It was picked mainly for the lack of security constraints,
it's uniform distribution and fast encoding technique.
After hashing the URL to an unsigned integer of 32 bit, 
the hashing is then transposed into a base62 string which is being used as the short-code.

---

The REST API exposed the pattern `/urls/` on port `8080` with the HTTP methods:
* `POST` (create)
* `GET` (read)
* `PUT` (update)
* `DELETE` (delete)

---

Tested with dockerized Cassandra instance `docker run --rm --name cassandra-db -d -p 9042:9042 cassandra`
