// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package console

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/kversion"
)

type EndpointCompatibility struct {
	KafkaClusterVersion string                          `json:"kafkaVersion"`
	Endpoints           []EndpointCompatibilityEndpoint `json:"endpoints"`
}

type EndpointCompatibilityEndpoint struct {
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	IsSupported bool   `json:"isSupported"`
}

// GetEndpointCompatibility requests API versions from brokers in order to figure out what Kowl endpoints
// can be offered to the frontend. If the broker does not support certain features which are required for a
// Kowl endpoint we can let the frontend know in advance, so that these features will be rendered as
// disabled.
func (s *Service) GetEndpointCompatibility(ctx context.Context) (EndpointCompatibility, error) {
	versionsRes, err := s.kafkaSvc.GetAPIVersions(ctx)
	if err != nil {
		return EndpointCompatibility{}, fmt.Errorf("failed to get kafka api version: %w", err)
	}
	versions := kversion.FromApiVersionsResponse(versionsRes)
	clusterVersion := versions.VersionGuess()

	// Required kafka requests per API endpoint
	type endpoint struct {
		URL      string
		Method   string
		Requests []kmsg.Request
	}
	endpointRequirements := []endpoint{
		{
			URL:      "/api/cluster/config",
			Method:   "GET",
			Requests: []kmsg.Request{&kmsg.DescribeConfigsRequest{}},
		},
		{
			URL:      "/api/consumer-groups",
			Method:   "GET",
			Requests: []kmsg.Request{&kmsg.DescribeGroupsRequest{}, &kmsg.ListGroupsRequest{}},
		},
		{
			URL:      "/api/topics/{topicName}/records",
			Method:   "DELETE",
			Requests: []kmsg.Request{&kmsg.DeleteRecordsRequest{}},
		},
		{
			URL:      "/api/consumer-groups/{groupId}",
			Method:   "PATCH",
			Requests: []kmsg.Request{&kmsg.OffsetCommitRequest{}},
		},
		{
			URL:      "/api/consumer-groups/{groupId}",
			Method:   "DELETE",
			Requests: []kmsg.Request{&kmsg.OffsetDeleteRequest{}},
		},
		{
			URL:      "/api/operations/reassign-partitions",
			Method:   "GET",
			Requests: []kmsg.Request{&kmsg.ListPartitionReassignmentsRequest{}},
		},
		{
			URL:      "/api/operations/reassign-partitions",
			Method:   "PATCH",
			Requests: []kmsg.Request{&kmsg.IncrementalAlterConfigsRequest{}, &kmsg.AlterPartitionAssignmentsRequest{}},
		},
		{
			URL:      "/api/quotas",
			Method:   "GET",
			Requests: []kmsg.Request{&kmsg.DescribeClientQuotasRequest{}},
		},
	}

	endpoints := make([]EndpointCompatibilityEndpoint, 0, len(endpointRequirements))
	for _, endpointReq := range endpointRequirements {
		endpointSupported := true
		for _, req := range endpointReq.Requests {
			_, isSupported := versions.LookupMaxKeyVersion(req.Key())
			if !isSupported {
				endpointSupported = false
			}
		}

		endpoints = append(endpoints, EndpointCompatibilityEndpoint{
			Endpoint:    endpointReq.URL,
			Method:      endpointReq.Method,
			IsSupported: endpointSupported,
		})
	}

	return EndpointCompatibility{
		KafkaClusterVersion: clusterVersion,
		Endpoints:           endpoints,
	}, nil
}
