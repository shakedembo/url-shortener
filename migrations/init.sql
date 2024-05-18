-- drop keyspace urlShortener;
create keyspace "urlShortener"
with replication = {'class': 'SimpleStrategy', 'replication_factor': 1 };

create table "urlShortener".urls (
    short_code text,
    url text,
    primary key ( short_code )
);