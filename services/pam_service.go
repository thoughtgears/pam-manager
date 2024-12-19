package services

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/iterator"

	privilegedaccessmanager "cloud.google.com/go/privilegedaccessmanager/apiv1"
	"cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
)

type PAMService struct {
	client *privilegedaccessmanager.Client
}

func NewPAMService(ctx context.Context, token oauth2.TokenSource) (*PAMService, error) {
	// If token is nil, create a PAM client without a token source
	if token == nil {
		client, err := privilegedaccessmanager.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create PAM client: %v", err)
		}

		return &PAMService{
			client: client,
		}, nil
	}

	// Create a PAM client with a token source that impersonates the user
	client, err := privilegedaccessmanager.NewClient(ctx, option.WithTokenSource(token))
	if err != nil {
		return nil, fmt.Errorf("failed to create PAM client: %v", err)
	}
	return &PAMService{
		client: client,
	}, nil
}

func (p *PAMService) GetGrants(ctx context.Context, project, entitlement string) ([]*privilegedaccessmanagerpb.Grant, error) {
	req := &privilegedaccessmanagerpb.ListGrantsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global/entitlements/%s", project, entitlement),
	}

	itr := p.client.ListGrants(ctx, req)

	var grants []*privilegedaccessmanagerpb.Grant
	for {
		grant, err := itr.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get grant: %v", err)
		}

		grants = append(grants, grant)
	}

	return grants, nil
}

func (p *PAMService) RequestGrant(ctx context.Context, projectId, entitlement, reason string, duration int64) (*privilegedaccessmanagerpb.Grant, error) {
	req := &privilegedaccessmanagerpb.CreateGrantRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global/entitlements/%s", projectId, entitlement),
		Grant: &privilegedaccessmanagerpb.Grant{
			RequestedDuration: &durationpb.Duration{Seconds: duration},
			Justification: &privilegedaccessmanagerpb.Justification{
				Justification: &privilegedaccessmanagerpb.Justification_UnstructuredJustification{
					UnstructuredJustification: reason,
				},
			},
		},
	}

	return p.client.CreateGrant(ctx, req)
}

func (p *PAMService) ApproveGrant(ctx context.Context, id, projectId, entitlement, reason string) (*privilegedaccessmanagerpb.Grant, error) {
	req := &privilegedaccessmanagerpb.ApproveGrantRequest{
		Name:   fmt.Sprintf("projects/%s/locations/global/entitlements/%s/grants/%s", projectId, entitlement, id),
		Reason: reason,
	}

	return p.client.ApproveGrant(ctx, req)
}

func (p *PAMService) RevokeGrant(ctx context.Context, id, projectId, entitlement, reason string) (*privilegedaccessmanagerpb.Grant, error) {
	req := &privilegedaccessmanagerpb.RevokeGrantRequest{
		Name:   fmt.Sprintf("projects/%s/locations/global/entitlements/%s/grants/%s", projectId, entitlement, id),
		Reason: reason,
	}

	op, err := p.client.RevokeGrant(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke grant: %v", err)
	}

	return op.Wait(ctx)
}
