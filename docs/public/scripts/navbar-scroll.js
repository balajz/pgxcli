// src/scripts/navbar-scroll.js
(function () {
  const header = document.querySelector('header.header');
  if (!header) return;

  // Translucency styles injected via JS so they apply after Starlight's styles
  Object.assign(header.style, {
    position: 'sticky',
    top: '0',
    zIndex: '50',
    transition: 'transform 0.3s ease, background 0.3s ease',
    backdropFilter: 'blur(12px)',
    WebkitBackdropFilter: 'blur(12px)',
    background: 'rgba(220, 232, 248, 0.78)', // light mode
  });

  // Respect explicit site theme first; only fall back to OS preference.
  const mq = window.matchMedia('(prefers-color-scheme: dark)');
  function applyThemeBg() {
    const currentTheme = document.documentElement.dataset.theme;
    const isDark = currentTheme ? currentTheme === 'dark' : mq.matches;
    header.style.background = isDark
      ? 'rgba(26, 39, 68, 0.78)'
      : 'rgba(220, 232, 248, 0.78)';
  }
  applyThemeBg();

  // Re-apply when the active site theme or OS preference changes.
  new MutationObserver(applyThemeBg).observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme'],
  });
  mq.addEventListener('change', applyThemeBg);

  // Hide on scroll down, show on scroll up
  let lastY = 0;
  let ticking = false;

  window.addEventListener('scroll', () => {
    if (!ticking) {
      requestAnimationFrame(() => {
        const currentY = window.scrollY;
        if (currentY <= 60) {
          // Always show at top of page
          header.style.transform = 'translateY(0)';
        } else if (currentY > lastY) {
          // Scrolling down — hide
          header.style.transform = 'translateY(-100%)';
        } else {
          // Scrolling up — show
          header.style.transform = 'translateY(0)';
        }
        lastY = currentY;
        ticking = false;
      });
      ticking = true;
    }
  });
})();
