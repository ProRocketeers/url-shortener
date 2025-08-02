package domain

import (
	"context"
	"time"

	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CleanupTask struct {
	Context  context.Context
	DB       *gorm.DB
	Interval time.Duration
}

func (t *CleanupTask) Run() {
	ticker := time.NewTicker(t.Interval)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				deleted, err := gorm.G[model.ShortLink](t.DB).Where("expires_at < ?", time.Now().UTC()).Delete(context.Background())
				if err != nil {
					log.Warn().Err(err).Msg("failed to cleanup expired links")
				} else {
					log.Info().Int("deleted", deleted).Msg("expired link cleanup")
				}
			case <-t.Context.Done():
				log.Info().Msg("stopping expired link cleanup")
				return
			}
		}
	}()
}
