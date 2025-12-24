# Haft Documentation

This is the documentation site for Haft, built with [Docusaurus](https://docusaurus.io/).

## Development

```bash
# Install dependencies
pnpm install

# Start development server
pnpm start

# Build for production
pnpm build

# Serve production build locally
pnpm serve
```

## Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch.

Manual deployment:

```bash
USE_SSH=true pnpm deploy
```

## Structure

```
docs/
├── docs/           # Documentation pages (Markdown)
├── src/
│   ├── components/ # React components
│   ├── css/        # Custom styles
│   └── pages/      # Custom pages
├── static/         # Static assets
└── docusaurus.config.ts
```
