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

  url: 'https://kashifkhn.github.io',
  baseUrl: '/haft/',

  organizationName: 'KashifKhn',
  projectName: 'haft',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

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

  themeConfig: {
    image: 'img/haft-social-card.jpg',
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
      additionalLanguages: ['bash', 'java', 'go', 'yaml', 'json'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
