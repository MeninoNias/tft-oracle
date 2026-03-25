package tierlist

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
	"github.com/MeninoNias/tft-oracle/backend/internal/crawler"
	"github.com/MeninoNias/tft-oracle/backend/sqlc/generated"
)

var _ tftv1connect.TierListServiceHandler = (*Service)(nil)

// Service handles tier list RPCs.
type Service struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	crawler *crawler.Crawler
}

// NewService creates a new TierList service.
func NewService(db *pgxpool.Pool, c *crawler.Crawler) *Service {
	return &Service{
		db:      db,
		queries: generated.New(db),
		crawler: c,
	}
}

func (s *Service) GetConsolidatedTierList(
	ctx context.Context,
	req *connect.Request[tftv1.GetConsolidatedTierListRequest],
) (*connect.Response[tftv1.GetConsolidatedTierListResponse], error) {
	patch := req.Msg.Patch
	tierFilter := req.Msg.TierFilter

	// If no patch specified, use latest available
	if patch == "" {
		latestPatch, err := s.queries.GetLatestConsolidatedPatch(ctx)
		if err == nil {
			patch = latestPatch
		}
	}

	var rows []generated.ConsolidatedTierList
	var err error

	if tierFilter != "" {
		rows, err = s.queries.GetConsolidatedTierListByTier(ctx, generated.GetConsolidatedTierListByTierParams{
			Patch:            patch,
			ConsolidatedTier: tierFilter,
		})
	} else {
		rows, err = s.queries.GetConsolidatedTierList(ctx, patch)
	}

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("query tier list: %w", err))
	}

	entries := make([]*tftv1.TierListEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, mapToProto(row))
	}

	// Determine last_updated from the most recent entry
	lastUpdated := ""
	if len(rows) > 0 && rows[0].UpdatedAt.Valid {
		lastUpdated = rows[0].UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
	}

	return connect.NewResponse(&tftv1.GetConsolidatedTierListResponse{
		Patch:       patch,
		Entries:     entries,
		LastUpdated: lastUpdated,
	}), nil
}

func (s *Service) TriggerCrawl(
	ctx context.Context,
	req *connect.Request[tftv1.TriggerCrawlRequest],
) (*connect.Response[tftv1.TriggerCrawlResponse], error) {
	if s.crawler == nil {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("crawler not configured"))
	}

	// Run crawl in a goroutine to not block the RPC
	go func() {
		bgCtx := context.Background()
		if err := s.crawler.Run(bgCtx); err != nil {
			log.Printf("tierlist: triggered crawl failed: %v", err)
		}
	}()

	return connect.NewResponse(&tftv1.TriggerCrawlResponse{
		Message: "crawl triggered",
	}), nil
}

func (s *Service) GetCrawlerStatus(
	ctx context.Context,
	req *connect.Request[tftv1.GetCrawlerStatusRequest],
) (*connect.Response[tftv1.GetCrawlerStatusResponse], error) {
	resp := &tftv1.GetCrawlerStatusResponse{}

	if s.crawler != nil {
		statuses := s.crawler.GetStatuses()
		for _, st := range statuses {
			lastScraped := ""
			if !st.LastScraped.IsZero() {
				lastScraped = st.LastScraped.Format("2006-01-02T15:04:05Z")
			}
			resp.Sources = append(resp.Sources, &tftv1.SourceStatus{
				Source:       st.Source,
				LastScraped:  lastScraped,
				Status:       st.Status,
				ErrorMessage: st.ErrorMessage,
			})
		}
	}

	return connect.NewResponse(resp), nil
}
