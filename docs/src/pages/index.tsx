import type {ReactNode} from 'react';
import {useState} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

import styles from './index.module.css';

const INSTALL_COMMAND = 'curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash';

function CopyButton({text}: {text: string}) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <button
      onClick={handleCopy}
      style={{
        background: 'transparent',
        border: 'none',
        cursor: 'pointer',
        padding: '0.25rem',
        marginLeft: '0.5rem',
        color: 'inherit',
        opacity: 0.7,
        transition: 'opacity 0.2s',
      }}
      onMouseOver={(e) => e.currentTarget.style.opacity = '1'}
      onMouseOut={(e) => e.currentTarget.style.opacity = '0.7'}
      title={copied ? 'Copied!' : 'Copy to clipboard'}
    >
      {copied ? (
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
      ) : (
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
        </svg>
      )}
    </button>
  );
}

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
            className="button button--secondary button--lg"
            style={{marginLeft: '1rem', backgroundColor: '#2d2a5e', color: '#f5f3ff', border: '2px solid #f5f3ff'}}
            href="https://github.com/KashifKhn/haft">
            GitHub
          </Link>
        </div>

        <div style={{marginTop: '2.5rem'}}>
          <p style={{marginBottom: '0.75rem', opacity: 0.9, fontSize: '0.9rem', color: '#f5f3ff'}}>Install with one command:</p>
          <div style={{
            backgroundColor: 'rgba(0,0,0,0.4)',
            padding: '0.75rem 1.5rem',
            borderRadius: '8px',
            fontSize: '1rem',
            fontFamily: 'monospace',
            display: 'inline-flex',
            alignItems: 'center',
            border: '1px solid rgba(255,255,255,0.2)',
            color: '#f5f3ff',
          }}>
            <code style={{margin: 0, color: '#f5f3ff', background: 'transparent'}}>{INSTALL_COMMAND}</code>
            <CopyButton text={INSTALL_COMMAND} />
          </div>
        </div>
      </div>
    </header>
  );
}

function Feature({title, description}: {title: string; description: string}) {
  return (
    <div className="col col--4">
      <div style={{padding: '1.5rem'}}>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

function HomepageFeatures() {
  return (
    <section style={{padding: '4rem 0'}}>
      <div className="container">
        <div className="row">
          <Feature
            title="Offline First"
            description="No internet required. All Spring Initializr metadata is bundled. Work anywhere, anytime."
          />
          <Feature
            title="Interactive TUI"
            description="Beautiful terminal interface with keyboard navigation. No more copy-pasting from browsers."
          />
          <Feature
            title="Code Generation"
            description="Generate complete CRUD resources with a single command. Entity, Repository, Service, Controller, DTOs."
          />
        </div>
        <div className="row">
          <Feature
            title="Smart Defaults"
            description="Sensible defaults that match industry standards. Java 21, Spring Boot 3.4, YAML config."
          />
          <Feature
            title="Back Navigation"
            description="Made a mistake? Press Esc to go back. No need to restart the entire wizard."
          />
          <Feature
            title="Dependency Search"
            description="Find any Spring starter instantly. Press / to search through all available dependencies."
          />
        </div>
      </div>
    </section>
  );
}

function DemoSection() {
  const demoGifUrl = useBaseUrl('/img/demo.gif');
  return (
    <section style={{padding: '4rem 0', backgroundColor: 'var(--ifm-background-surface-color)'}}>
      <div className="container">
        <Heading as="h2" style={{textAlign: 'center', marginBottom: '2rem'}}>
          See It In Action
        </Heading>
        <div style={{display: 'flex', justifyContent: 'center'}}>
          <img 
            src={demoGifUrl}
            alt="Haft Demo" 
            style={{
              maxWidth: '800px', 
              width: '100%', 
              borderRadius: '12px',
              boxShadow: '0 4px 20px rgba(0,0,0,0.3)'
            }}
          />
        </div>
      </div>
    </section>
  );
}

function ComparisonSection() {
  return (
    <section style={{padding: '4rem 0', backgroundColor: 'var(--ifm-background-surface-color)'}}>
      <div className="container">
        <Heading as="h2" style={{textAlign: 'center', marginBottom: '2rem'}}>
          Why Haft?
        </Heading>
        <div style={{maxWidth: '600px', margin: '0 auto'}}>
          <table style={{width: '100%'}}>
            <thead>
              <tr>
                <th></th>
                <th style={{textAlign: 'center'}}>Spring Initializr</th>
                <th style={{textAlign: 'center'}}>Haft</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Project Bootstrap</td>
                <td style={{textAlign: 'center'}}>Yes</td>
                <td style={{textAlign: 'center'}}>Yes</td>
              </tr>
              <tr>
                <td>Works Offline</td>
                <td style={{textAlign: 'center'}}>No</td>
                <td style={{textAlign: 'center'}}>Yes</td>
              </tr>
              <tr>
                <td>Resource Generation</td>
                <td style={{textAlign: 'center'}}>No</td>
                <td style={{textAlign: 'center'}}>Yes</td>
              </tr>
              <tr>
                <td>Dependency Management</td>
                <td style={{textAlign: 'center'}}>No</td>
                <td style={{textAlign: 'center'}}>Yes</td>
              </tr>
              <tr>
                <td>Terminal UI</td>
                <td style={{textAlign: 'center'}}>No</td>
                <td style={{textAlign: 'center'}}>Yes</td>
              </tr>
              <tr>
                <td>Lifecycle Companion</td>
                <td style={{textAlign: 'center'}}>No</td>
                <td style={{textAlign: 'center'}}>Yes</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
}

function CodeBlock({children}: {children: string}) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(children);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div style={{position: 'relative'}}>
      <pre style={{
        backgroundColor: 'var(--ifm-code-background)',
        padding: '1.5rem',
        borderRadius: '8px',
        overflow: 'auto',
        margin: 0,
      }}>
        <code>{children}</code>
      </pre>
      <button
        onClick={handleCopy}
        style={{
          position: 'absolute',
          top: '0.75rem',
          right: '0.75rem',
          background: 'var(--ifm-color-emphasis-200)',
          border: 'none',
          borderRadius: '4px',
          cursor: 'pointer',
          padding: '0.5rem',
          color: 'var(--ifm-color-content)',
          opacity: 0.7,
          transition: 'opacity 0.2s',
        }}
        onMouseOver={(e) => e.currentTarget.style.opacity = '1'}
        onMouseOut={(e) => e.currentTarget.style.opacity = '0.7'}
        title={copied ? 'Copied!' : 'Copy to clipboard'}
      >
        {copied ? (
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
        ) : (
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
        )}
      </button>
    </div>
  );
}

function QuickStartSection() {
  const quickStartCode = `# Install
curl -fsSL https://raw.githubusercontent.com/KashifKhn/haft/main/install.sh | bash

# Create a new project
haft init

# Generate a CRUD resource
cd my-project
haft generate resource User`;

  return (
    <section style={{padding: '4rem 0'}}>
      <div className="container">
        <Heading as="h2" style={{textAlign: 'center', marginBottom: '2rem'}}>
          Quick Start
        </Heading>
        <div style={{maxWidth: '700px', margin: '0 auto'}}>
          <CodeBlock>{quickStartCode}</CodeBlock>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title="The Spring Boot CLI"
      description="The Spring Boot CLI that Spring forgot to build. Generate projects, resources, and manage dependencies with an interactive TUI. Works offline.">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <DemoSection />
        <ComparisonSection />
        <QuickStartSection />
      </main>
    </Layout>
  );
}
