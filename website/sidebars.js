const sidebars = {
    docsSidebar: [
        {
            type: 'doc',
            id: 'guides/getting-started',
            label: 'Getting Started',
        },
        {
            type: 'doc',
            id: 'guides/comparison-with-pgcli',
            label: 'pgxcli vs pgcli',
        },
        {
            type: 'category',
            label: 'Usage Guides',
            items: [
                'guides/connecting',
                'guides/features',
                'guides/configuration',
                'guides/special-commands',
            ],
        },
        {
            type: 'category',
            label: 'Reference',
            items: [
                'reference/cli-reference',
                'reference/environment-variables',
            ],
        },
    ],
};
export default sidebars;
