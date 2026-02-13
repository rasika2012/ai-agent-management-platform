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



import { lazy, type ComponentType, type FC } from "react";
import { metaData as overviewMetadata } from "@agent-management-platform/overview";
import { metaData as buildMetadata } from "@agent-management-platform/build";
import { metaData as deploymentMetadata } from "@agent-management-platform/deploy";
import { metaData as testMetadata } from "@agent-management-platform/test";
import { metaData as tracesMetadata } from "@agent-management-platform/traces";
import { metaData as logsMetadata } from "@agent-management-platform/logs";
import { metaData as metricsMetadata } from "@agent-management-platform/metrics";

export * from "./Login";

// Navigation pages - imported normally (needed upfront for nav)
export const LazyOverviewOrg =
  overviewMetadata.pages.organization.organizationOverview.component as FC;
export const LazyOverviewProject =
  overviewMetadata.pages.project.projectOverview.component as FC;
export const LazyOverviewComponent =
  overviewMetadata.pages.component.componentOverview.component as FC;
export const LazyBuildComponent = buildMetadata.pages.component.build.component as FC;
export const LazyDeploymentComponent: FC = deploymentMetadata.pages.component.deploy.component;
export const LazyTestComponent = testMetadata.pages.component.test.component as FC;
export const LazyTracesComponent = tracesMetadata.pages.component.trace.component as FC;
export const LazyLogsComponent = logsMetadata.pages.component.logs.component as FC;
export const LazyMetricsComponent = metricsMetadata.pages.component.metrics.component as FC;

// Create pages - lazy loaded (only needed when user creates something)
export const LazyAddNewAgent = lazy(() =>
  import("@agent-management-platform/add-new-agent").then((module) => ({
    default: module.metaData.pages.component.addNewAgent.component as ComponentType,
  }))
);

export const LazyAddNewProject = lazy(() =>
  import("@agent-management-platform/add-new-project").then((module) => ({
    default: module.metaData.pages.component.addNewProject.component as ComponentType,
  }))
);


