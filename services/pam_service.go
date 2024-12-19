package services

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	privilegedaccessmanager "cloud.google.com/go/privilegedaccessmanager/apiv1"
	"cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
)

type PAMService struct {
	client *privilegedaccessmanager.Client
}

func NewPAMService() *PAMService {
	return &PAMService{}
}

func (p *PAMService) RequestGrant(ctx context.Context, token oauth2.TokenSource, projectId, entitlement, justification string, duration int64) (*privilegedaccessmanagerpb.Grant, error) {
	client, err := privilegedaccessmanager.NewClient(ctx, option.WithTokenSource(token))
	if err != nil {
		return nil, fmt.Errorf("failed to create PAM client: %v", err)
	}

	req := &privilegedaccessmanagerpb.CreateGrantRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global/entitlements/%s", projectId, entitlement),
		Grant: &privilegedaccessmanagerpb.Grant{
			RequestedDuration: &durationpb.Duration{Seconds: duration},
			Justification: &privilegedaccessmanagerpb.Justification{
				Justification: &privilegedaccessmanagerpb.Justification_UnstructuredJustification{
					UnstructuredJustification: justification,
				},
			},
		},
	}

	return client.CreateGrant(ctx, req)
}
