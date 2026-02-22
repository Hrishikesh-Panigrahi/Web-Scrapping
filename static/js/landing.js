/**
 * Landing page – download form and URL validation.
 * Validates Instagram Reels and YouTube video/shorts URLs before submit.
 */
(function () {
  'use strict';

  var DOWNLOAD_PATH = '/download';
  var URL_INPUT_ID = 'landing_url';
  var RESULT_ID = 'landingResult';

  /**
   * Check if URL is a valid Instagram Reel or YouTube video/shorts URL.
   * @param {string} url - Full URL (with protocol).
   * @returns {boolean}
   */
  function isValidInstagramOrYouTubeUrl(url) {
    try {
      var u = new URL(url);
      var host = u.hostname.toLowerCase().replace(/^www\./, '');

      if (host === 'instagram.com') {
        return /^\/(reel|p)\//.test(u.pathname);
      }

      if (host === 'youtube.com' || host === 'youtu.be') {
        if (host === 'youtu.be') return u.pathname.length > 1;
        return /^\/watch(\?|$)/.test(u.pathname) || /^\/shorts\//.test(u.pathname);
      }

      return false;
    } catch (e) {
      return false;
    }
  }

  /**
   * Normalize URL (add https:// if missing).
   */
  function normalizeUrl(url) {
    return url.match(/^https?:\/\//i) ? url : 'https://' + url;
  }

  /**
   * Show error message in result area.
   */
  function showError(message) {
    var el = document.getElementById(RESULT_ID);
    if (el) {
      el.innerHTML = '<p class="text-amber-500 dark:text-amber-400">' + message + '</p>';
    }
  }

  /**
   * Validate URL and block HTMX request if invalid.
   */
  function onBeforeRequest(evt) {
    var pathInfo = evt.detail && evt.detail.pathInfo;
    if (!pathInfo || pathInfo.requestPath !== DOWNLOAD_PATH) return;

    var input = document.getElementById(URL_INPUT_ID);
    var url = (input && input.value || '').trim();

    if (!url) {
      evt.preventDefault();
      showError('Please enter a URL.');
      return;
    }

    var normalized = normalizeUrl(url);
    if (!isValidInstagramOrYouTubeUrl(normalized)) {
      evt.preventDefault();
      showError('Please enter a valid Instagram Reel or YouTube video/shorts URL.');
    }
  }

  document.body.addEventListener('htmx:beforeRequest', onBeforeRequest);
})();
