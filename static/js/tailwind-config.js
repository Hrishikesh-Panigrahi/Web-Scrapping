/**
 * Tailwind config – shared across pages.
 * Must load before cdn.tailwindcss.com.
 */
(function () {
  'use strict';
  window.tailwind = window.tailwind || {};
  window.tailwind.config = {
    darkMode: 'class',
    theme: {
      extend: {
        fontFamily: { sans: ['Inter', 'system-ui', 'sans-serif'] },
        animation: { fadeIn: 'fadeIn 0.35s ease-out' },
        keyframes: {
          fadeIn: {
            '0%': { opacity: '0', transform: 'translateY(6px)' },
            '100%': { opacity: '1', transform: 'translateY(0)' },
          },
        },
      },
    },
  };
})();
