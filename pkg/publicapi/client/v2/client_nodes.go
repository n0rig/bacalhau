package client

import (
	"context"

	"github.com/bacalhau-project/bacalhau/pkg/publicapi/apimodels"
)

const nodesPath = "/api/v1/orchestrator/nodes"

type Nodes struct {
	client *Client
}

// Nodes returns a handle on the nodes endpoints.
func (c *Client) Nodes() *Nodes {
	return &Nodes{client: c}
}

// Get is used to get a node by ID.
func (c *Nodes) Get(ctx context.Context, r *apimodels.GetNodeRequest) (*apimodels.GetNodeResponse, error) {
	var resp apimodels.GetNodeResponse
	if err := c.client.get(ctx, nodesPath+"/"+r.NodeID, r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// List is used to list all nodes in the cluster.
func (c *Nodes) List(ctx context.Context, r *apimodels.ListNodesRequest) (*apimodels.ListNodesResponse, error) {
	var resp apimodels.ListNodesResponse
	if err := c.client.list(ctx, nodesPath, r, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
