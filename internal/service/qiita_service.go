package service

import (
	"context"

	"github.com/hiro-env/grpcaggregator/pkg/qiita"
)

type QiitaService struct {
	qiita.UnimplementedQiitaServiceServer
}

func (s *QiitaService) SearchArticles(ctx context.Context, req *qiita.SearchRequest) (*qiita.SearchResponse, error) {
	return &qiita.SearchResponse{
		Articles: []*qiita.Article{
			{
				Title:       "サンプル記事",
				Author:      "サンプル著者",
				Url:         "https://qiita.com",
				PublishedAt: "2023-01-01T00:00:00Z",
			},
		},
	}, nil
}
