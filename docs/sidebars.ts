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
        'commands/dev',
        'commands/info',
        'commands/routes',
        'commands/stats',
        'commands/template',
        'commands/upgrade',
        'commands/completion',
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
        'guides/custom-templates',
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
    {
      type: 'category',
      label: 'Editor Integrations',
      collapsed: true,
      link: {
        type: 'doc',
        id: 'integrations/overview',
      },
      items: [
        {
          type: 'category',
          label: 'Neovim',
          collapsed: true,
          link: {
            type: 'doc',
            id: 'integrations/neovim/index',
          },
          items: [
            'integrations/neovim/installation',
            'integrations/neovim/configuration',
            'integrations/neovim/usage',
            'integrations/neovim/troubleshooting',
          ],
        },
        {
          type: 'doc',
          id: 'integrations/vscode/index',
          label: 'VS Code',
        },
        {
          type: 'doc',
          id: 'integrations/intellij/index',
          label: 'IntelliJ IDEA',
        },
      ],
    },
    'contributing',
    'roadmap',
  ],
};

export default sidebars;
