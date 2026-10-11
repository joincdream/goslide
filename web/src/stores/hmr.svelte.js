/**
 * Goslide HMR and Build Error/Warning Reactive Store (Svelte 5 Runes)
 */

export class HMRStore {
  error = $state(null); // { title: string, message: string, file: string } | null
  warnings = $state([]); // Array<Diagnostic>
  warningTimer = null;

  setError(err) {
    this.error = err;
  }

  clearError() {
    this.error = null;
  }

  setWarnings(diagnostics) {
    if (this.warningTimer) {
      clearTimeout(this.warningTimer);
      this.warningTimer = null;
    }
    this.warnings = diagnostics || [];
    if (this.warnings.length > 0) {
      this.warningTimer = setTimeout(() => {
        this.clearWarnings();
      }, 6000);
    }
  }

  clearWarnings() {
    if (this.warningTimer) {
      clearTimeout(this.warningTimer);
      this.warningTimer = null;
    }
    this.warnings = [];
  }

  clearAll() {
    this.clearError();
    this.clearWarnings();
  }
}

export const hmr = new HMRStore();
