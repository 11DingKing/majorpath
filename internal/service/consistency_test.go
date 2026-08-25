package service

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"testing"
)

func TestValidateRecommendation(t *testing.T) {
	s := ConsistencyService{}
	good := domain.Recommendation{ID: "r", StudentID: "s", Version: 1}
	if e := s.ValidateRecommendation(context.Background(), good); e != nil {
		t.Fatal(e)
	}
	for _, r := range []domain.Recommendation{{}, {ID: "r", Version: 1}, {ID: "r", StudentID: "s"}, {ID: "r", StudentID: "s", Version: 0}} {
		if e := s.ValidateRecommendation(context.Background(), r); e == nil {
			t.Fatal("invalid accepted")
		}
	}
}
func TestConsistencyContext(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	if e := (ConsistencyService{}).ValidateRecommendation(ctx, domain.Recommendation{ID: "r", StudentID: "s", Version: 1}); e == nil {
		t.Fatal("cancel ignored")
	}
}
