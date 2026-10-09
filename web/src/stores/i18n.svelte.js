/**
 * Goslide Frontend Micro i18n Store (Svelte 5 Runes)
 * Zero-dependency, offline-first translation store powered by Go backend HTML injection.
 */

const defaultData = (typeof window !== 'undefined' && window.__GOSLIDE_I18N__) || {
  locale: 'en',
  messages: {}
};

let currentLocale = $state(defaultData.locale || 'en');
let catalog = $state(defaultData.messages || {});

/**
 * Returns the translated string for a given key, with fallback and positional formatting.
 * @param {string} key - The message key (e.g. 'ui.toolbar.pen')
 * @param {string} [fallback=''] - Fallback text if the key is not in catalog
 * @param {...any} args - Format arguments to replace %s or %d
 * @returns {string}
 */
export function t(key, fallback = '', ...args) {
  let msg = catalog[key] || fallback || key;
  if (args.length > 0) {
    args.forEach(arg => {
      msg = msg.replace(/%[sd]/, String(arg));
    });
  }
  return msg;
}

/**
 * Returns the currently active locale.
 * @returns {string}
 */
export function getLocale() {
  return currentLocale;
}

/**
 * Updates the locale and optionally replaces the translation catalog.
 * @param {string} locale
 * @param {Record<string, string>} [newMessages]
 */
export function setLocale(locale, newMessages) {
  currentLocale = locale;
  if (newMessages) {
    catalog = newMessages;
  }
}

/**
 * Returns a copy of the entire translation catalog.
 * @returns {Record<string, string>}
 */
export function getMessages() {
  return { ...catalog };
}
