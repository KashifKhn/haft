import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Haft',
  tagline: 'The Spring Boot CLI that Spring forgot to build',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://haft.kashifkhan.dev',
  baseUrl: '/',

  organizationName: 'KashifKhn',
  projectName: 'haft',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  scripts: [
    {
      src: 'https://cloud.umami.is/script.js',
      defer: true,
      'data-website-id': 'bd96a9a7-e568-40f2-ab08-cce8b92e7a31',
    },
  ],

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/KashifKhn/haft/tree/main/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  headTags: [
    {
      tagName: 'meta',
      attributes: {
        name: 'google-site-verification',
        content: 'GsnVvhnHNdZC3YMYInZxanh97DATIV3-Ooup2YHYcS8',
      },
    },
  ],

  themeConfig: {
    image: 'img/haft-social-card.jpg',
    metadata: [
      {name: 'keywords', content: 'haft, spring boot, cli, java, spring, scaffolding, code generator, spring initializr'},
      {name: 'description', content: 'Haft is a fast Spring Boot CLI for scaffolding projects, generating components, and managing Spring Boot applications from the terminal.'},
      {name: 'robots', content: 'index, follow'},
      {property: 'og:title', content: 'Haft - The Spring Boot CLI'},
      {property: 'og:description', content: 'Haft is a fast Spring Boot CLI for scaffolding projects, generating components, and managing Spring Boot applications from the terminal.'},
      {property: 'og:type', content: 'website'},
      {property: 'og:url', content: 'https://haft.kashifkhan.dev'},
      {property: 'og:image', content: 'https://haft.kashifkhan.dev/img/haft-social-card.jpg'},
      {name: 'twitter:card', content: 'summary_large_image'},
      {name: 'twitter:title', content: 'Haft - The Spring Boot CLI'},
      {name: 'twitter:description', content: 'Haft is a fast Spring Boot CLI for scaffolding projects, generating components, and managing Spring Boot applications from the terminal.'},
      {name: 'twitter:image', content: 'https://haft.kashifkhan.dev/img/haft-social-card.jpg'},
    ],
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: false,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: 'Haft',
      logo: {
        alt: 'Haft Logo',
        src: 'img/logo-light.png',
        srcDark: 'img/logo-dark.png',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'dropdown',
          label: 'Integrations',
          position: 'left',
          items: [
            {
              label: 'Neovim',
              to: '/docs/integrations/neovim',
            },
            {
              label: 'VS Code',
              to: '/docs/integrations/vscode',
            },
            {
              label: 'IntelliJ IDEA',
              to: '/docs/integrations/intellij',
            },
          ],
        },
        {
          href: 'https://github.com/KashifKhn/haft',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Documentation',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/getting-started',
            },
            {
              label: 'Installation',
              to: '/docs/installation',
            },
            {
              label: 'Commands',
              to: '/docs/commands/init',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub Discussions',
              href: 'https://github.com/KashifKhn/haft/discussions',
            },
            {
              label: 'Issues',
              href: 'https://github.com/KashifKhn/haft/issues',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/KashifKhn/haft',
            },
            {
              label: 'Contributing',
              href: 'https://github.com/KashifKhn/haft/blob/main/CONTRIBUTING.md',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Haft. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'java', 'go', 'yaml', 'json', 'lua', 'vim'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
