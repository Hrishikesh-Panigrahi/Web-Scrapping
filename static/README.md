# Static Assets

## Structure

```
static/
├── css/
│   └── app.css       # Shared styles (HTMX indicators)
├── js/
│   ├── tailwind-config.js   # Tailwind theme (load before Tailwind CDN)
│   ├── dark-mode.js         # Theme toggle, localStorage, init
│   ├── app.js               # Lucide icons, common setup
│   └── landing.js           # Landing page URL validation
└── README.md
```

## Load Order (templates)

1. `tailwind-config.js` – before Tailwind CDN
2. Tailwind, HTMX, Lucide (CDN)
3. `app.css`
4. `dark-mode.js` – runs init immediately
5. Page-specific scripts (`landing.js` on landing only)
6. `app.js` – Lucide init on DOMContentLoaded

## Scripts

- **dark-mode.js** – Toggle, init from localStorage, applies `dark` class on `html`. Source of truth for theme.
- **landing.js** – Validates Instagram/YouTube URLs before `/download` submit via `htmx:beforeRequest`.
