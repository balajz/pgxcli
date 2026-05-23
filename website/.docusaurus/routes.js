import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/markdown-page',
    component: ComponentCreator('/markdown-page', '53a'),
    exact: true
  },
  {
    path: '/search',
    component: ComponentCreator('/search', '822'),
    exact: true
  },
  {
    path: '/support',
    component: ComponentCreator('/support', 'ce4'),
    exact: true
  },
  {
    path: '/docs',
    component: ComponentCreator('/docs', '449'),
    routes: [
      {
        path: '/docs',
        component: ComponentCreator('/docs', 'b39'),
        routes: [
          {
            path: '/docs',
            component: ComponentCreator('/docs', '558'),
            routes: [
              {
                path: '/docs/guides/comparison-with-pgcli',
                component: ComponentCreator('/docs/guides/comparison-with-pgcli', '4b5'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/guides/configuration',
                component: ComponentCreator('/docs/guides/configuration', 'c46'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/guides/connecting',
                component: ComponentCreator('/docs/guides/connecting', 'ee4'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/guides/features',
                component: ComponentCreator('/docs/guides/features', '01f'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/guides/getting-started',
                component: ComponentCreator('/docs/guides/getting-started', '7b9'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/guides/special-commands',
                component: ComponentCreator('/docs/guides/special-commands', 'e29'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/reference/cli-reference',
                component: ComponentCreator('/docs/reference/cli-reference', 'f4a'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/docs/reference/environment-variables',
                component: ComponentCreator('/docs/reference/environment-variables', '746'),
                exact: true,
                sidebar: "docsSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/',
    component: ComponentCreator('/', '2e1'),
    exact: true
  },
  {
    path: '/',
    component: ComponentCreator('/', 'e5f'),
    exact: true
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];
