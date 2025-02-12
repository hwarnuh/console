/**
 * Copyright 2022 Redpanda Data, Inc.
 *
 * Use of this software is governed by the Business Source License
 * included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
 *
 * As of the Change Date specified in that file, in accordance with
 * the Business Source License, use of this software will be governed
 * by the Apache License, Version 2.0
 */

import React from 'react';
import { render } from '@testing-library/react';
import AclList from './AclList';
import { observable } from 'mobx';
import { ResourcePatternType } from '../../../../state/restInterfaces';

it('renders an empty table when no data is present', () => {
    const store = observable({
        isAuthorizerEnabled: true,
        aclResources: [],
    });

    const { getByText } = render(<AclList acl={store} />);
    expect(getByText('No Data')).toBeInTheDocument();
});

it('a table with one entry', () => {
    const store = observable({
        isAuthorizerEnabled: true,
        aclResources: [
            {
                resourceType: 'Topic',
                resourceName: 'Test Topic',
                resourcePatternType: ResourcePatternType.UNKNOWN,
                acls: [
                    {
                        principal: 'test principal',
                        host: 'test host',
                        operation: 'test operation',
                        permissionType: 'test permission type',
                    },
                ],
            },
        ],
    });

    const { getByText } = render(<AclList acl={store} />);

    expect(getByText('Topic')).toBeInTheDocument();
    expect(getByText('Test Topic')).toBeInTheDocument();
    expect(getByText('0')).toBeInTheDocument();
    expect(getByText('test principal')).toBeInTheDocument();
    expect(getByText('test host')).toBeInTheDocument();
    expect(getByText('test operation')).toBeInTheDocument();
    expect(getByText('test permission type')).toBeInTheDocument();
});

it('informs user about missing permission to view ACLs', () => {
    const { getByText } = render(<AclList acl={null} />);
    expect(getByText('You do not have the necessary permissions to view ACLs')).toBeInTheDocument();
});

it('informs user about missing authorizer config in Kafka cluster', () => {
    const store = observable({
        isAuthorizerEnabled: false,
        aclResources: [],
    });

    const { getByText } = render(<AclList acl={store} />);
    expect(getByText('There\'s no authorizer configured in your Kafka cluster')).toBeInTheDocument();
});
