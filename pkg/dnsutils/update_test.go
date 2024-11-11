package dnsutils

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/libdns/libdns"
	"github.com/stretchr/testify/require"
)

type MockProvider struct {
	getRecords func(ctx context.Context, zone string) ([]libdns.Record, error)
	setRecords func(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error)
}

func (m *MockProvider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	return m.getRecords(ctx, zone)
}

func (m *MockProvider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return m.setRecords(ctx, zone, recs)
}

func TestPlanUpdate(t *testing.T) {
	var provider *MockProvider

	setup := func(t *testing.T) {
		provider = &MockProvider{}
		provider.getRecords = func(ctx context.Context, zone string) ([]libdns.Record, error) {
			return []libdns.Record{
				{
					ID:    uuid.NewString(),
					Type:  "TXT",
					Name:  "abc.xd.squeeze.one",
					Value: "old.squeeze.one",
				},
				{
					ID:    uuid.NewString(),
					Type:  "A",
					Name:  "abc.xd.squeeze.one",
					Value: "old.squeeze.one",
				},
				{
					ID:    uuid.NewString(),
					Type:  "A",
					Name:  "lol.squeeze.one",
					Value: "old.squeeze.one",
				},
			}, nil
		}

		provider.setRecords = func(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
			t.Error("did not expect setRecords to be called")
			t.FailNow()
			return nil, nil
		}
	}

	t.Run("no records", func(t *testing.T) {
		setup(t)

		provider.getRecords = func(ctx context.Context, zone string) ([]libdns.Record, error) {
			return []libdns.Record{}, nil
		}

		records, err := PlanUpdate(context.Background(), regexp.MustCompile(`.*\.xd\.squeeze\.one`), "A", "squeeze.one", "new.squeeze.one", provider)
		require.NoError(t, err)
		require.Empty(t, records)
	})

	t.Run("plans to update one record without modifying it", func(t *testing.T) {
		setup(t)

		records, err := PlanUpdate(context.Background(), regexp.MustCompile(`.*\.xd\.squeeze\.one`), "A", "squeeze.one", "new.squeeze.one", provider)
		require.NoError(t, err)
		require.Len(t, records, 1)

		require.Equal(t, "old.squeeze.one", records[0].Value)
		require.Equal(t, "A", records[0].Type)
		require.Equal(t, "abc.xd.squeeze.one", records[0].Name)
	})

	t.Run("updates nothing on non-matching regex domain", func(t *testing.T) {
		setup(t)

		records, err := PlanUpdate(context.Background(), regexp.MustCompile(`.*\.xd\.squeeze\.de`), "A", "squeeze.one", "new.squeeze.one", provider)
		require.NoError(t, err)
		require.Empty(t, records)
	})
}
