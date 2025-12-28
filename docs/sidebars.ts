import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    'getting-started',
    'installation',
    'why-haft',
    {
      type: 'category',
      label: 'Commands',
      collapsed: false,
      items: [
        'commands/init',
        'commands/generate',
        'commands/security',
        'commands/add',
        'commands/remove',
        'commands/completion',
        'commands/dev',
        'commands/info',
        'commands/routes',
        'commands/stats',
      ],
    },
    {
      type: 'category',
      label: 'Guides',
      collapsed: false,
      items: [
        'guides/wizard-navigation',
        'guides/dependencies',
        'guides/project-structure',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      collapsed: true,
      items: [
        'reference/configuration',
        'reference/templates',
      ],
    },
    'contributing',
    'roadmap',
  ],
};

export default sidebars;
