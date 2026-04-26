import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import styles from './index.module.css';

export default function Home(): ReactNode {
  return (
    <Layout
      title="pgxcli — Interactive PostgreSQL client for your terminal"
      description="Syntax highlighting, smart autocompletion, and a configurable REPL — out of the box. A single Go binary, zero dependencies.">
      <main className={styles.hero}>

        {/* ── Left column ─────────────────────── */}
        <div className={styles.heroLeft}>

            <h1 className={styles.heroTitle}>
            A PostgreSQL client<br />
            written in Go<br />
            for your terminal.
            </h1>

          <p className={styles.heroDesc}>
            Syntax highlighting, autocompletion, and a single Go binary.
          </p>

          <div className={styles.heroCtas}>
            <Link className={styles.btnPrimary} to="/docs/guides/getting-started">
              Get Started
            </Link>
            <Link className={styles.btnSecondary} to="http://localhost:3000/docs/guides/getting-started#installation">
              Download
            </Link>
          </div>
        </div>

        {/* ── Right column — app screenshot ──── */}
        <div className={styles.heroRight}>
          <img
            src="img/app-screenshot.png"
            alt="pgxcli in action"
            className={styles.screenshot}
          />
        </div>

      </main>
    </Layout>
  );
}
