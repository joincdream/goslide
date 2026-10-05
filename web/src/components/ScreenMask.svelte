<script>
  import { deck } from '../stores/deck.svelte.js';
</script>

{#if deck.isBlackout}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 bg-black z-[9998] cursor-pointer"
    onclick={() => deck.toggleBlackout(false)}
  ></div>
{/if}

{#if deck.isWhiteout}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="fixed inset-0 bg-white z-[9998] cursor-pointer"
    onclick={() => deck.toggleWhiteout(false)}
  ></div>
{/if}
