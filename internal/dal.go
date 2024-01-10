package internal

import (
	"context"
	"log"

	"url-shortener/config"

	"github.com/gocql/gocql"
)

type DAO[T, TE any] interface {
	Create(ctx context.Context, key T, val TE) error
	Read(ctx context.Context, key T) (TE, error)
	Update(ctx context.Context, oldKey, newKey T, newVal TE) error
	Delete(ctx context.Context, key T) error
}

type CassandraDAO struct {
	logger  *log.Logger
	cluster *gocql.ClusterConfig
}

func NewCassandraDAO(config *config.DbConfig, logger *log.Logger) *CassandraDAO {
	return &CassandraDAO{
		logger:  logger,
		cluster: gocql.NewCluster(config.Host),
	}
}

func (c *CassandraDAO) Create(ctx context.Context, key, val string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return err
	}
	defer session.Close()

	if err := session.Query(`INSERT INTO "urlShortener".urls (short_code, url) VALUES (?, ?)`,
		key, val).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to insert new entry to cassandra key: `%s`, val: `%s`. Error: `%v`", key, val, err)
		return err
	}

	return nil
}

func (c *CassandraDAO) Read(ctx context.Context, key string) (string, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return "", err
	}
	defer session.Close()

	url := ""
	if err := session.Query(`SELECT url FROM "urlShortener".urls WHERE short_code = (?)`,
		key).WithContext(ctx).Consistency(gocql.One).Scan(&url); err != nil {
		c.logger.Printf("Error occurred trying to read entry from cassandra key: `%s. Error: `%v`", key, err)
		return "", err
	}

	return url, nil
}

func (c *CassandraDAO) Update(ctx context.Context, oldKey, newKey, val string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return err
	}
	defer session.Close()

	if err := session.Query(`INSERT INTO "urlShortener".urls (short_code, url) VALUES (?, ?)`,
		newKey, val).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to create the entry id: "+
			"`%s`. Error: `%v`", newKey, err)
		return err
	}

	if err := session.Query(`DELETE FROM "urlShortener".urls WHERE short_code = (?) IF EXISTS`,
		oldKey).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to delete entry from cassandra id: `%s. Error: `%v`", newKey, err)
		return err
	}

	return nil
}

func (c *CassandraDAO) Delete(ctx context.Context, key string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return err
	}
	defer session.Close()

	if err := session.Query(`DELETE FROM "urlShortener".urls WHERE short_code = (?) IF EXISTS`,
		key).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to delete entry from cassandra key: `%s. Error: `%v`", key, err)
		return err
	}
	//TODO: select applied flag from cql return in order to determine if the row exists

	return nil
}
