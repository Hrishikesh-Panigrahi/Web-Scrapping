/**
 * Shared app initialization.
 * Lucide icons and other common setup.
 */
(function () {
  'use strict';

  function initLucide() {
    if (typeof lucide !== 'undefined') {
      lucide.createIcons();
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initLucide);
  } else {
    initLucide();
  }

  // Re-init Lucide after dark mode toggle (icons may need refresh)
  document.body.addEventListener('htmx:afterSwap', initLucide);
})();
