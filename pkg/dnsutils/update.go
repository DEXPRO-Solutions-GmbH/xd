package dnsutils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/libdns/libdns"
)

func PlanUpdate(ctx context.Context, domain *regexp.Regexp, recordType, zone, newValue string, provider libdns.RecordGetter) ([]libdns.Record, error) {
	records, err := provider.GetRecords(ctx, zone)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve records of zone %v, %w", zone, err)
	}

	var recordsToUpdate []libdns.Record

	for _, record := range records {
		if record.Type != recordType {
			continue
		}

		if !domain.MatchString(record.Name) {
			continue
		}

		recordsToUpdate = append(recordsToUpdate, record)
	}

	return recordsToUpdate, nil
}

func Update(ctx context.Context, zone string, records []libdns.Record, newValue string, provider libdns.RecordSetter) error {
	for i := 0; i < len(records); i++ {
		records[i].Value = newValue
	}

	_, err := provider.SetRecords(ctx, zone, records)
	if err != nil {
		return fmt.Errorf("failed to update %d records of zone %v, %w", len(records), zone, err)
	}

	return nil
}
