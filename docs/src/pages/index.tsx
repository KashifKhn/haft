import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/getting-started">
            Get Started
          </Link>
          <Link
            className="button button--outline button--secondary button--lg"
            style={{marginLeft: '1rem'}}
            href="https://github.com/KashifKhn/haft">
            GitHub
          </Link>
        </div>
        <div style={{marginTop: '2rem', display: 'flex', flexDirection: 'column', gap: '0.75rem', alignItems: 'center'}}>
          <code style={{
            backgroundColor: 'rgba(255,255,255,0.1)',
            padding: '0.5rem 1rem',
            borderRadius: '4px',
            fontSize: '0.95rem',
            fontFamily: 'monospace'
          }}>
            curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash
          </code>
          <span style={{fontSize: '0.85rem', opacity: 0.8}}>or</span>
          <code style={{
            backgroundColor: 'rgba(255,255,255,0.1)',
            padding: '0.5rem 1rem',
            borderRadius: '4px',
            fontSize: '0.95rem',
            fontFamily: 'monospace'
          }}>
            go install github.com/KashifKhn/haft/cmd/haft@latest
          </code>
        </div>
      </div>
    </header>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title="The Spring Boot CLI"
      description="The Spring Boot CLI that Spring forgot to build. Generate projects, resources, and manage dependencies with an interactive TUI.">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
