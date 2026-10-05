import { mount } from 'svelte';
import './app.css';
import App from './App.svelte';

function initGoslide() {
  let target = document.getElementById('goslide-app');
  if (!target) {
    target = document.createElement('div');
    target.id = 'goslide-app';
    document.body.appendChild(target);
  }

  mount(App, { target });
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initGoslide);
} else {
  initGoslide();
}
