import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import icon from 'astro-icon';

// https://astro.build/config
export default defineConfig({
	site: 'https://balaji01-4d.github.io',
	base: '/pgxcli/',
	integrations: [
		starlight({
			title: 'Pgxcli',
			customCss: [
				'./src/styles/tokens.css',
				'./src/styles/custom.css',
			],
			head: [
				{
					tag: 'script',
					attrs: { src: '/pgxcli/scripts/navbar-scroll.js', defer: true },
				},
			],
			expressiveCode: {
				themes: ['github-dark'],
				styleOverrides: {
					codeBackground: '#0e1628',
					frames: {
						frameBoxShadowCssValue: 'none',
						editorActiveTabBackground: '#0e1628',
						editorTabBarBackground: '#1a2744',
						terminalTitlebarBackground: '#1a2744',
						terminalBackground: '#0e1628',
					},
				},
			},
			social: [
				{
					label: 'GitHub',
					href: 'https://github.com/balaji01-4d/pgxcli',
					icon: 'github',
				},
			],
			sidebar: [
				{ label: 'Getting Started', slug: 'guides/getting-started' },
				{ label: 'pgxcli vs pgcli', slug: 'guides/comparison-with-pgcli' },
				{
					label: 'Usage Guides',
					items: [
						{ label: 'Connecting', slug: 'guides/connecting' },
						{ label: 'Configuration', slug: 'guides/configuration' },
						{ label: 'Special Commands', slug: 'guides/special-commands' },
					],
				},
				{
					label: 'Reference',
					items: [
						{ label: 'CLI Flags', slug: 'reference/cli-reference' },
						{ label: 'Environment Variables', slug: 'reference/environment-variables' },
					],
				},
			],
		}),
		icon(),
	],
});
