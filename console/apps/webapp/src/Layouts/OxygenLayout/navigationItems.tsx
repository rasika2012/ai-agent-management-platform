/**
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { BarChart3 as AutoGraphOutlined, Binoculars } from "@wso2/oxygen-ui-icons-react";
import {
  generatePath,
  matchPath,
  useLocation,
  useParams,
} from "react-router-dom";
import { absoluteRouteMap } from "@agent-management-platform/types";
import {
  useGetAgent,
  useListEnvironments,
} from "@agent-management-platform/api-client";
import { metaData as overviewMetadata } from "@agent-management-platform/overview";
import { metaData as buildMetadata } from "@agent-management-platform/build";
import { metaData as testMetadata } from "@agent-management-platform/test";
import { metaData as tracesMetadata } from "@agent-management-platform/traces";
import { metaData as logsMetadata } from "@agent-management-platform/logs";
import { metaData as metricsMetadata } from "@agent-management-platform/metrics";
import { metaData as deploymentMetadata } from "@agent-management-platform/deploy";
import type { NavigationItem, NavigationSection } from "./LeftNavigation";

/**
 * TODO: Use nav bar instead of navigate to the items.
 */

const overviewComponentPage = overviewMetadata.pages.component.componentOverview;
const buildComponentPage = buildMetadata.pages.component.build;
const deploymentComponentPage = deploymentMetadata.pages.component.deploy;
const testComponentPage = testMetadata.pages.component.test;
const tracesComponentPage = tracesMetadata.pages.component.trace;
const logsComponentPage = logsMetadata.pages.component.logs;
const metricsComponentPage = metricsMetadata.pages.component.metrics;

export function useNavigationItems(): Array<NavigationSection | NavigationItem> {
  const { orgId, projectId, agentId, envId } = useParams();
  const { data: agent, isLoading: isLoadingAgent } = useGetAgent({
    agentName: agentId,
    orgName: orgId,
    projName: projectId,
  });
  const { data: environments, isLoading: isLoadingEnvironments } =
    useListEnvironments({
      orgName: orgId,
    });
  const defaultEnv = envId ?? environments?.[0]?.name;
  const { pathname } = useLocation();

  if (isLoadingAgent || (isLoadingEnvironments && agentId)) {
    return [];
  }

  if (
    agent?.provisioning.type === "external" &&
    agentId &&
    projectId &&
    orgId
  ) {
    return [
      {
        label: overviewComponentPage.title,
        type: "item",
        icon: <overviewComponentPage.icon size={20} />,
        isActive: !!matchPath(
          absoluteRouteMap.children.org.children.projects.children.agents.path,
          pathname
        ),
        href: generatePath(
          absoluteRouteMap.children.org.children.projects.children.agents.path,
          { orgId, projectId, agentId }
        ),
      },
      {
        title: "Observability",
        type: "section",
        icon: <AutoGraphOutlined />,
        items: [
          {
            label: tracesComponentPage.title,
            type: "item",
            icon: <tracesComponentPage.icon size={20} />,
            isActive: !!matchPath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.traces.wildPath,
              pathname
            ),
            href: generatePath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.traces.path,
              { orgId, projectId, agentId, envId: defaultEnv }
            ),
          },
        ],
      },
    ];
  }

  if (orgId && projectId && agentId && defaultEnv) {
    return [
      {
        label: overviewComponentPage.title,
        type: "item",
        icon: <overviewComponentPage.icon size={20} />,
        isActive: !!matchPath(
          absoluteRouteMap.children.org.children.projects.children.agents.path,
          pathname
        ),
        href: generatePath(
          absoluteRouteMap.children.org.children.projects.children.agents.path,
          { orgId, projectId, agentId }
        ),
      },
      {
        label: buildComponentPage.title,
        type: "item",
        icon: <buildComponentPage.icon size={20} />,
        isActive: !!matchPath(
          absoluteRouteMap.children.org.children.projects.children.agents
            .children.build.wildPath,
          pathname
        ),
        href: generatePath(
          absoluteRouteMap.children.org.children.projects.children.agents
            .children.build.path,
          { orgId, projectId, agentId }
        ),
      },
      {
        label: deploymentComponentPage.title,
        type: "item",
        icon: <deploymentComponentPage.icon size={20} />,
        isActive: !!matchPath(
          absoluteRouteMap.children.org.children.projects.children.agents
            .children.deployment.wildPath,
          pathname
        ),
        href: generatePath(
          absoluteRouteMap.children.org.children.projects.children.agents
            .children.deployment.path,
          { orgId, projectId, agentId }
        ),
      },
      {
        label: testComponentPage.title,
        type: "item",
        icon: <testComponentPage.icon size={20} />,
        isActive: !!matchPath(
          absoluteRouteMap.children.org.children.projects.children.agents
            .children.environment.children.tryOut.wildPath,
          pathname
        ),
        href: generatePath(
          absoluteRouteMap.children.org.children.projects.children.agents
            .children.environment.children.tryOut.path,
          { orgId, projectId, agentId, envId: defaultEnv }
        ),
      },
      {
        title: "Observability",
        type: "section",
        icon: <Binoculars  />,
        items: [
          {
            label: tracesComponentPage.title,
            type: "item",
            icon: <tracesComponentPage.icon size={20} />,
            isActive: !!matchPath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.traces
                .wildPath,
              pathname
            ),
            href: generatePath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.traces
                .path,
              { orgId, projectId, agentId, envId: defaultEnv }
            ),
          },
          {
            label: logsComponentPage.title,
            type: "item",
            icon: <logsComponentPage.icon size={20} />,
            isActive: !!matchPath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.logs
                .wildPath,
              pathname
            ),
            href: generatePath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.logs
                .path,
              { orgId, projectId, agentId, envId: defaultEnv }
            ),
          },
          {
            label: metricsComponentPage.title,
            type: "item",
            icon: <metricsComponentPage.icon size={20} />,
            isActive: !!matchPath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.metrics
                .wildPath,
              pathname
            ),
            href: generatePath(
              absoluteRouteMap.children.org.children.projects.children.agents
                .children.environment.children.observability.children.metrics
                .path,
              { orgId, projectId, agentId, envId: defaultEnv }
            ),
          },
        ],
      },
    ];
  }
  if (orgId && projectId) {
    return [
      {
        label: "Agents",
        type: "item",
        icon: <overviewComponentPage.icon size={20} />,
        href: generatePath(
          absoluteRouteMap.children.org.children.projects.path,
          { orgId, projectId }
        ),
        isActive:
          !!matchPath(
            absoluteRouteMap.children.org.children.projects.path,
            pathname
          ) ||
          !!matchPath(
            absoluteRouteMap.children.org.children.projects.children.agents
              .wildPath,
            pathname
          ),
      },
    ];
  }
  if (orgId) {
    return [
      {
        label: "Projects",
        type: "item",
        icon: <overviewComponentPage.icon size={20} />,
        href: generatePath(absoluteRouteMap.children.org.path, { orgId }),
        isActive: !!matchPath(absoluteRouteMap.children.org.path, pathname),
      },
    ];
  }
  return [];
}
