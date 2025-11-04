package connector

import (
	"context"
	"fmt"

	baseConnector "github.com/conductorone/baton-metabase/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/actions"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"google.golang.org/protobuf/types/known/structpb"
)

func (c *Connector) RegisterActionManager(ctx context.Context) (connectorbuilder.CustomActionManager, error) {
	actionManager := actions.NewActionManager(ctx)

	err := actionManager.RegisterAction(ctx, baseConnector.EnableUserAction.Name, baseConnector.EnableUserAction, c.EnableUserV056)
	if err != nil {
		return nil, err
	}

	err = actionManager.RegisterAction(ctx, baseConnector.DisableUserAction.Name, baseConnector.DisableUserAction, c.DisableUserV056)
	if err != nil {
		return nil, err
	}

	return actionManager, nil
}

func (c *Connector) EnableUserV056(ctx context.Context, args *structpb.Struct) (*structpb.Struct, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	ann := annotations.New()

	userIdField, ok := args.Fields["userId"]
	if !ok || userIdField == nil {
		return nil, ann, fmt.Errorf("userId field is required")
	}
	userId := userIdField.GetStringValue()
	if userId == "" {
		return nil, ann, fmt.Errorf("userId cannot be empty")
	}

	user, rateLimitDesc, err := c.vBaseClient.GetUserByID(ctx, userId)
	if rateLimitDesc != nil {
		ann.WithRateLimiting(rateLimitDesc)
	}
	if err != nil {
		return nil, ann, fmt.Errorf("failed to fetch user %s: %w", userId, err)
	}

	if user.IsActive {
		l.Debug("user already active, skipping enable", zap.String("userId", userId))
		return &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"success": structpb.NewBoolValue(true),
			},
		}, ann, nil
	}

	resp, ann2, err := c.vBaseConnector.EnableUser(ctx, args)
	if ann2 != nil {
		ann.Merge(ann2...)
	}
	if err != nil {
		return nil, ann, fmt.Errorf("failed to enable user %s: %w", userId, err)
	}

	return resp, ann, nil
}

func (c *Connector) DisableUserV056(ctx context.Context, args *structpb.Struct) (*structpb.Struct, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	ann := annotations.New()

	userIdField, ok := args.Fields["userId"]
	if !ok || userIdField == nil {
		return nil, ann, fmt.Errorf("userId field is required")
	}
	userId := userIdField.GetStringValue()
	if userId == "" {
		return nil, ann, fmt.Errorf("userId cannot be empty")
	}

	user, rateLimitDesc, err := c.vBaseClient.GetUserByID(ctx, userId)
	if rateLimitDesc != nil {
		ann.WithRateLimiting(rateLimitDesc)
	}
	if err != nil {
		return nil, ann, fmt.Errorf("failed to fetch user %s: %w", userId, err)
	}

	if !user.IsActive {
		l.Debug("user already inactive, skipping disable", zap.String("userId", userId))
		return &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"success": structpb.NewBoolValue(true),
			},
		}, ann, nil
	}

	resp, ann2, err := c.vBaseConnector.DisableUser(ctx, args)
	if ann2 != nil {
		ann.Merge(ann2...)
	}
	if err != nil {
		return nil, ann, fmt.Errorf("failed to disable user %s: %w", userId, err)
	}

	return resp, ann, nil
}
