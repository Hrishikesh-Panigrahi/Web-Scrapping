/**
 * Dark mode toggle – source of truth for site-wide theme.
 * Persists preference in localStorage, applies on load.
 * Use with Tailwind darkMode: 'class'.
 */
(function () {
  'use strict';

  var STORAGE_KEY = 'darkMode';

  /**
   * Apply theme from localStorage or system preference.
   */
  function init() {
    var stored = localStorage.getItem(STORAGE_KEY);
    var prefersDark =
      stored === 'true' ||
      (stored === null && window.matchMedia('(prefers-color-scheme: dark)').matches);

    if (prefersDark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  /**
   * Toggle dark mode and persist.
   */
  function toggle() {
    document.documentElement.classList.toggle('dark');
    localStorage.setItem(STORAGE_KEY, document.documentElement.classList.contains('dark'));
    if (typeof lucide !== 'undefined') lucide.createIcons();
  }

  /**
   * Bind toggle button by ID.
   */
  function bindToggle(buttonId) {
    var btn = document.getElementById(buttonId);
    if (btn) {
      btn.addEventListener('click', toggle);
    }
  }

  // Run init immediately (before paint)
  init();

  // Bind on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () {
      bindToggle('darkModeToggle');
    });
  } else {
    bindToggle('darkModeToggle');
  }
})();
