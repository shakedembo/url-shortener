package internal

import (
	"context"
	"url-shortener/pkg"
)

type UrlShortenerController struct {
	service UrlShortenerService
}

func NewUrlShortenerController(service UrlShortenerService) *UrlShortenerController {
	return &UrlShortenerController{
		service: service,
	}
}

func (u *UrlShortenerController) Create(ctx context.Context, req pkg.CreateRequest) (pkg.CreateResponse, error, int) {
	id, err := u.service.HandleCreate(ctx, req.Url)
	if err != nil {
		return pkg.CreateResponse{}, err, 406
	}

	return pkg.CreateResponse{
		Id: id,
	}, nil, 201
}

func (u *UrlShortenerController) Read(ctx context.Context, req pkg.ReadRequest) (pkg.ReadResponse, error, int) {
	url, err := u.service.HandleRead(ctx, req.Id)
	if err != nil {
		return pkg.ReadResponse{}, err, 406
	}

	return pkg.ReadResponse{
		Url: url,
	}, nil, 200
}

func (u *UrlShortenerController) Replace(ctx context.Context, req pkg.UpdateRequest) (pkg.UpdateResponse, error, int) {
	id, err := u.service.HandleReplace(ctx, req.Id, req.Url)
	if err != nil {
		return pkg.UpdateResponse{}, err, 406
	}

	return pkg.UpdateResponse{
		Id: id,
	}, nil, 200
}

func (u *UrlShortenerController) Delete(ctx context.Context, req pkg.DeleteRequest) (pkg.DeleteResponse, error, int) {
	err := u.service.HandleDelete(ctx, req.Id)
	if err != nil {
		return pkg.DeleteResponse{}, err, 406
	}

	return pkg.DeleteResponse{}, nil, 200
}
