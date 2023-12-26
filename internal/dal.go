package internal

import (
	"context"
	"log"
	"url-shortener/config"

	"github.com/gocql/gocql"
)

type DAO interface {
	Create(ctx context.Context, id, url string) error
	Read(ctx context.Context, id string) (string, error)
	Update(ctx context.Context, oldId, newId, url string) error
	Delete(ctx context.Context, id string) error
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

func (c *CassandraDAO) Create(ctx context.Context, id, url string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return err
	}
	defer session.Close()

	if err := session.Query(`INSERT INTO "urlShortener".urls (short_code, url) VALUES (?, ?)`,
		id, url).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to insert new entry to cassandra id: `%s`, url: `%s`. Error: `%v`", id, url, err)
		return err
	}

	return nil
}

func (c *CassandraDAO) Read(ctx context.Context, id string) (string, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return "", err
	}
	defer session.Close()

	url := ""
	if err := session.Query(`SELECT url FROM "urlShortener".urls WHERE short_code = (?)`,
		id).WithContext(ctx).Consistency(gocql.One).Scan(&url); err != nil {
		c.logger.Printf("Error occurred trying to read entry from cassandra id: `%s. Error: `%v`", id, err)
		return "", err
	}

	return url, nil
}

func (c *CassandraDAO) Update(ctx context.Context, oldId, newId, url string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return err
	}
	defer session.Close()

	if err := session.Query(`INSERT INTO "urlShortener".urls (short_code, url) VALUES (?, ?)`,
		newId, url).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to create the entry id: "+
			"`%s`. Error: `%v`", newId, err)
		return err
	}

	if err := session.Query(`DELETE FROM "urlShortener".urls WHERE short_code = (?) IF EXISTS`,
		oldId).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to delete entry from cassandra id: `%s. Error: `%v`", newId, err)
		return err
	}

	return nil
}

func (c *CassandraDAO) Delete(ctx context.Context, id string) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		c.logger.Printf("Error occurred trying to create a session to db. Error: `%v`", err)
		return err
	}
	defer session.Close()

	if err := session.Query(`DELETE FROM "urlShortener".urls WHERE short_code = (?) IF EXISTS`,
		id).WithContext(ctx).Exec(); err != nil {
		c.logger.Printf("Error occurred trying to delete entry from cassandra id: `%s. Error: `%v`", id, err)
		return err
	}
	//TODO: select applied flag from cql return in order to determine if the row exists

	return nil
}
