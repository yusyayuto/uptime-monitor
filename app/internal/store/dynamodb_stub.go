//go:build !aws
// +build !aws

package store

import "context"

// NewDynamo returns an in-memory repository when built without the aws tag.
func NewDynamo(ctx context.Context, region, sitesTable, statusTable string) Repository {
	return NewInMemory()
}
