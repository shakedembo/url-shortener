package internal

import (
	"context"
	"log"
	"url-shortener/internal/utils"
)

type UrlShortenerService interface {
	HandleCreate(ctx context.Context, in string) (id string, err error)
	HandleRead(ctx context.Context, id string) (url string, err error)
	HandleReplace(ctx context.Context, id, url string) (newId string, err error)
	HandleDelete(ctx context.Context, id string) (err error)
}

type SimpleUrlShortenerService struct {
	hasher utils.HashProvider[string]
	dao    DAO
	logger *log.Logger
}

func NewUrlShorteningService(
	hasher utils.HashProvider[string],
	dao DAO,
	logger *log.Logger,
) *SimpleUrlShortenerService {
	return &SimpleUrlShortenerService{
		hasher: hasher,
		dao:    dao,
		logger: logger,
	}
}

func (u *SimpleUrlShortenerService) HandleCreate(ctx context.Context, in string) (id string, err error) {
	ctx = context.WithValue(ctx, utils.Url, in)

	id = u.hasher.Get(in)
	ctx = context.WithValue(ctx, utils.Code, id)

	if err := u.dao.Create(ctx, id, in); err != nil {
		u.logger.Printf("Error occurred trying to create an entry in the db. id: %s, url: `%s`. Error: `%v`", id, in, err)
		return "", err
	}

	return id, nil
}

func (u *SimpleUrlShortenerService) HandleRead(ctx context.Context, id string) (url string, err error) {
	ctx = context.WithValue(ctx, utils.Code, id)

	if url, err = u.dao.Read(ctx, id); err != nil {
		u.logger.Printf("Error occurred trying to read an entry in the db. id: `%s`, url: `%s`. Error: `%v`", id, url, err)
		return "", err
	}

	ctx = context.WithValue(ctx, utils.Url, url)

	return url, err
}

func (u *SimpleUrlShortenerService) HandleReplace(ctx context.Context, id, url string) (newId string, err error) {
	ctx = context.WithValue(ctx, utils.Url, url)

	newId = u.hasher.Get(id)
	u.logger.Printf("Updating context code `%s` to `%s`", id, newId)
	ctx = context.WithValue(ctx, utils.Code, newId)

	if err := u.dao.Update(ctx, id, newId, url); err != nil {
		u.logger.Printf("Error occurred trying to update an entry in the db. id: `%s`, newId: `%s` url: `%s`",
			id,
			newId,
			url,
		)
		return "", err
	}

	ctx = context.WithValue(ctx, utils.Url, url)

	return newId, err
}

func (u *SimpleUrlShortenerService) HandleDelete(ctx context.Context, id string) (err error) {
	ctx = context.WithValue(ctx, utils.Code, id)

	if err := u.dao.Delete(ctx, id); err != nil {
		u.logger.Printf("Error occurred trying to delete an entry in the db. id: `%s`. Error: `%v`", id, err)
		return err
	}

	return nil
}
