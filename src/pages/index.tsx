import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import styles from './index.module.css';

export default function Home(): ReactNode {
  return (
    <Layout
      title="pgxcli — Interactive PostgreSQL client for your terminal"
      description="Syntax highlighting, smart autocompletion, and a configurable REPL — out of the box. A single Go binary, zero dependencies.">
      <main>

        {/* ── Hero ─────────────────────────────── */}
        <section className={styles.hero}>

          {/* Screenshot — large, full width */}
          <div className={styles.heroRight}>
            <img
              src="img/app-screenshot.png"
              alt="pgxcli in action"
              className={styles.screenshot}
            />
          </div>

          {/* Text below screenshot — brand text on left, content centered */}
          <div className={styles.heroBottom}>
            <span className={styles.brandSide}>pgxcli</span>

            <div className={styles.heroLeft}>
              <h1 className={styles.heroTitle}>
                A PostgreSQL client<br />
                written in <span className={styles.goText}>Go</span><br />
                for your terminal.
              </h1>

              <p className={styles.heroDesc}>
                Syntax highlighting, autocompletion, and a single Go binary.
              </p>

              <div className={styles.heroCtas}>
                <Link className={styles.btnPrimary} to="/docs/guides/getting-started">
                  Get Started
                </Link>
                <Link className={styles.btnSecondary} to="/docs/guides/getting-started#installation">
                  Download
                </Link>
              </div>
            </div>
          </div>
        </section>

        {/* ── Orca mascot ──────────────────────── */}
        <section className={styles.orcaSection}>
          <img
            src="img/logo2.png"
            alt="pgxcli mascot — an orca"
            className={styles.orcaImg}
          />
          <p className={styles.orcaCaption}>
            Why an orca and not the elephant? Honestly, I just love orcas.
          </p>
        </section>

      </main>
    </Layout>
  );
}
