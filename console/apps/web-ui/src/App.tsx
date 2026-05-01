import { absoluteRouteMap, CoreUI as AgentManagerApp, MountPoints, type ExternalModule } from '@agent-management-platform/core-ui'
import '@agent-management-platform/core-ui/dist/index.css'

const externalModules: ExternalModule[] = [
  {
    kind: "page",
    moduleName: "External Org Page Module Example",
    mountPoint: MountPoints.OrgLevelPage,
    pageComponent: () => <div>This is an external page module mounted at org-level.</div>,
    path: "external-org-page",
  },
  {
    kind: "nav-item",
    moduleName: "External Org Nav Item Example",
    mountPoint: MountPoints.LeftNavItem,
    icon: "🚀",
    title: "External Org Nav Item",
    level: "org",
    route: absoluteRouteMap.children.org.path + "/external-org-page",
  },
  {
    kind: "nav-item",
    moduleName: "External Project Nav Item Example",
    mountPoint: MountPoints.LeftNavItem,
    icon: "🗂️",
    title: "External Project Nav Item",
    level: "project",
    route: absoluteRouteMap.children.org.children.projects.path + "/external-project-page",
  },
  {
    kind: "nav-item",
    moduleName: "External Component Nav Item Example",
    mountPoint: MountPoints.LeftNavItem,
    icon: "🤖",
    title: "External Component Nav Item",
    level: "component",
    route: absoluteRouteMap.children.org.children.projects.children.agents.path + "/external-component-page",
  },
  {
    kind: "page",
    moduleName: "External Project Page Module Example",
    mountPoint: MountPoints.ProjectLevelPage,
    pageComponent: () => <div>This is an external page module mounted at project-level.</div>,
    path: "external-project-page",
  },
  {
    kind: "page",
    moduleName: "External Component Page Module Example",
    mountPoint: MountPoints.ComponentLevelPage,
    pageComponent: () => <div>This is an external page module mounted at component-level.</div>,
    path: "external-component-page",
  },
  {
    kind: "component",
    moduleName: "External Component Module Example",
    mountPoint: MountPoints.TopRightPanel,
    component: () => <div>This is an external component injected into a slot.</div>,
  },
]
function App() {
  return <AgentManagerApp externalPageModules={externalModules} />
}

export default App
