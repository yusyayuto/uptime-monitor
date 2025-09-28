//go:build aws
// +build aws

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func NewDynamo(ctx context.Context, region, sitesTable, statusTable string) Repository {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return NewInMemory()
	}
	client := dynamodb.NewFromConfig(cfg)
	return &dynamoRepo{
		client:      client,
		sitesTable:  sitesTable,
		statusTable: statusTable,
	}
}

type dynamoRepo struct {
	client      *dynamodb.Client
	sitesTable  string
	statusTable string
}

func (d *dynamoRepo) PutSite(ctx context.Context, site Site) error {
	item := map[string]types.AttributeValue{
		"id":         &types.AttributeValueMemberS{Value: site.ID},
		"url":        &types.AttributeValueMemberS{Value: site.URL},
		"created_at": &types.AttributeValueMemberS{Value: site.CreatedAt.Format(time.RFC3339Nano)},
	}
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.sitesTable,
		Item:      item,
	})
	return err
}

func (d *dynamoRepo) ListSites(ctx context.Context) ([]Site, error) {
	output, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &d.sitesTable,
	})
	if err != nil {
		return nil, err
	}
	sites := make([]Site, 0, len(output.Items))
	for _, item := range output.Items {
		site := Site{}
		if v, ok := item["id"].(*types.AttributeValueMemberS); ok {
			site.ID = v.Value
		}
		if v, ok := item["url"].(*types.AttributeValueMemberS); ok {
			site.URL = v.Value
		}
		if v, ok := item["created_at"].(*types.AttributeValueMemberS); ok {
			if ts, err := time.Parse(time.RFC3339Nano, v.Value); err == nil {
				site.CreatedAt = ts
			}
		}
		sites = append(sites, site)
	}
	return sites, nil
}

func (d *dynamoRepo) PutStatus(ctx context.Context, status Status) error {
	payload, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("marshal status: %w", err)
	}
	item := map[string]types.AttributeValue{
		"site_id":    &types.AttributeValueMemberS{Value: status.SiteID},
		"checked_at": &types.AttributeValueMemberS{Value: status.CheckedAt.Format(time.RFC3339Nano)},
		"payload":    &types.AttributeValueMemberB{Value: payload},
	}
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.statusTable,
		Item:      item,
	})
	return err
}

func (d *dynamoRepo) LatestStatus(ctx context.Context, siteID string) (Status, error) {
	key := map[string]types.AttributeValue{
		"site_id": &types.AttributeValueMemberS{Value: siteID},
	}
	output, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.statusTable,
		Key:       key,
	})
	if err != nil {
		return Status{}, err
	}
	if output.Item == nil {
		return Status{}, ErrNotFound
	}
	payloadAttr, ok := output.Item["payload"].(*types.AttributeValueMemberB)
	if !ok {
		return Status{}, fmt.Errorf("missing payload field")
	}
	var status Status
	if err := json.Unmarshal(payloadAttr.Value, &status); err != nil {
		return Status{}, fmt.Errorf("unmarshal payload: %w", err)
	}
	return status, nil
}
