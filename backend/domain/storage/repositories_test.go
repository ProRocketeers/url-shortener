package storage

import (
	"context"
	"testing"
	"time"

	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/ProRocketeers/url-shortener/domain/query"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRequestInfoRepo(t *testing.T) *RequestInfoRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.RequestInfo{}); err != nil {
		t.Fatalf("failed to migrate request_info table: %v", err)
	}

	return &RequestInfoRepository{
		Repository: Repository{
			DB: query.Use(db),
		},
	}
}

func seedRequestInfo(t *testing.T, repo *RequestInfoRepository, path string, requestID string) {
	t.Helper()

	now := time.Now().UTC()
	info := &model.RequestInfo{
		RequestId: requestID,
		Timestamp: now,
		Path:      path,
		Method:    "GET",
	}

	if err := repo.Create(context.Background(), info); err != nil {
		t.Fatalf("failed to seed request info: %v", err)
	}
}

func TestRequestInfoRepository_ListBySlug(t *testing.T) {
	repo := setupRequestInfoRepo(t)
	seedRequestInfo(t, repo, "/v1/abc123", "req-1")
	seedRequestInfo(t, repo, "/v1/xyz999", "req-2")
	seedRequestInfo(t, repo, "/v1/abc123", "req-3")

	infos, total, err := repo.ListBySlug(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if len(infos) != 2 {
		t.Fatalf("expected 2 items, got %d", len(infos))
	}

	for _, info := range infos {
		if info.Path != "/v1/abc123" {
			t.Fatalf("expected path /v1/abc123, got %s", info.Path)
		}
	}
}

func TestRequestInfoRepository_PaginatedListBySlug(t *testing.T) {
	repo := setupRequestInfoRepo(t)
	seedRequestInfo(t, repo, "/v1/slug-a", "req-1")
	seedRequestInfo(t, repo, "/v1/slug-a", "req-2")
	seedRequestInfo(t, repo, "/v1/slug-b", "req-3")

	infos, total, err := repo.PaginatedListBySlug(context.Background(), "slug-a", 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	if len(infos) != 1 {
		t.Fatalf("expected 1 item for pagination, got %d", len(infos))
	}

	if infos[0].Path != "/v1/slug-a" {
		t.Fatalf("expected path /v1/slug-a, got %s", infos[0].Path)
	}
}
