package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/hiro-env/grpcaggregator/pkg/qiita"
)

type QiitaService struct {
	qiita.UnimplementedQiitaServiceServer
}

type qiitaRSS struct {
	Entries []qiitaEntry `xml:"entry"`
}

type qiitaEntry struct {
	Title     string `xml:"title"`
	Author    string `xml:"author>name"`
	Url       string `xml:"url"`
	Published string `xml:"published"`
}

func (s *QiitaService) SearchArticles(ctx context.Context, req *qiita.SearchRequest) (*qiita.SearchResponse, error) {
	query := req.Query
	if query == "" {
		// TODO Datadog連携
	}

	feedURL := fmt.Sprintf("https://qiita.com/tags/%s/feed", query)

	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feed qiitaRSS
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	var articles []*qiita.Article

	for _, entry := range feed.Entries {
		pubTime, _ := time.Parse(time.RFC3339, entry.Published)
		articles = append(articles, &qiita.Article{
			Title:       entry.Title,
			Author:      entry.Author,
			Url:         entry.Url,
			PublishedAt: pubTime.Format(time.RFC3339),
		})
	}

	return &qiita.SearchResponse{Articles: articles}, nil
}
