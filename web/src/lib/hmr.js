/**
 * Goslide HMR (Hot Module Replacement) and Live Reload Client
 * Listens to Server-Sent Events (/events) and triggers incremental slide DOM patching.
 */

import { deck } from '../stores/deck.svelte.js';
import { hmr } from '../stores/hmr.svelte.js';

export function initHMR() {
  if (typeof window === 'undefined') return;
  if (!window.EventSource || !window.location.protocol.startsWith('http')) {
    return;
  }

  const es = new EventSource('/events');

  es.onmessage = (e) => {
    const data = e.data;
    if (!data) return;

    let payload;
    try {
      payload = JSON.parse(data);
    } catch {
      if (data === 'reload') {
        payload = { type: 'reload' };
      }
    }

    if (!payload) return;

    if (payload.type === 'reload') {
      hmr.clearAll();
      console.log('[goslide] Live reload triggered');
      location.reload();
    } else if (payload.type === 'patch') {
      hmr.clearError();
      const patches = payload.patches || [];
      patches.forEach((patch) => {
        const oldEl = document.querySelector(`.slide-card[data-slide="${patch.index}"]`);
        if (!oldEl) return;

        const temp = document.createElement('div');
        temp.innerHTML = (patch.html || '').trim();
        const newEl = temp.firstElementChild;
        if (!newEl) return;

        if (oldEl.classList.contains('active')) {
          newEl.classList.add('active');
        }

        newEl.style.transition = 'opacity 0.15s ease-in-out';
        newEl.style.opacity = '0.7';

        oldEl.replaceWith(newEl);
        requestAnimationFrame(() => {
          newEl.style.opacity = '1';
        });

        // Notify Svelte DeckStore to synchronize slide DOM element reference
        deck.patchSlide(patch.index - 1, newEl);
        window.dispatchEvent(
          new CustomEvent('goslide:slide-patched', {
            detail: { index: patch.index - 1, element: newEl }
          })
        );

        // Re-execute Math (KaTeX) and Diagrams (Mermaid) if loaded
        if (typeof window.renderMathInElement === 'function') {
          window.renderMathInElement(newEl, {
            delimiters: [
              { left: '$$', right: '$$', display: true },
              { left: '$', right: '$', display: false }
            ],
            throwOnError: false
          });
        }

        if (typeof window.mermaid !== 'undefined' && newEl.querySelector('.mermaid')) {
          window.mermaid.run({ nodes: newEl.querySelectorAll('.mermaid') });
        }
      });
    } else if (payload.type === 'error') {
      console.error('[goslide] Build error:', payload.message);
      hmr.setError({
        title: payload.title || 'Build Error',
        message: payload.message || 'Unknown error occurred',
        file: payload.file || ''
      });
    } else if (payload.type === 'warning') {
      console.warn('[goslide] DSL warnings:', payload.diagnostics);
      hmr.setWarnings(payload.diagnostics || []);
    }
  };

  es.onerror = () => {
    // EventSource will automatically retry connecting
  };

  return () => {
    es.close();
  };
}
