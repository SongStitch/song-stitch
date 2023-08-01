<script>
  import { onMount } from 'svelte';
  import darkModeIcon from '../assets/images/dark_mode.png';
  import { darkMode } from '../stores';

  function applyDarkMode() {
    window.document.body.classList.toggle('dark-mode');
    localStorage.setItem('darkMode', darkMode.toString());
  }

  function toggle() {
    $darkMode = !$darkMode;
    console.log($darkMode);
    applyDarkMode();
  }

  onMount(async () => {
    if (darkMode) {
      applyDarkMode();
    }
  });
</script>

<div class="dark-mode-icon">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <a href={'#'} class="img-link" on:click={toggle}>
    <img
      src={darkModeIcon}
      class="darkmode-icon-img"
      alt="darkmode icon"
      width="512"
      height="512"
    />
  </a>
</div>

<style>
  .dark-mode-icon {
    cursor: pointer;
    position: relative;
    text-align: right;
    display: inline;
    float: right;
    padding-right: 1.5em;
    padding-top: 1em;
  }

  .darkmode-icon-img {
    max-width: 1.5em;
    width: auto;
    height: auto;
  }

  :global(dialog.dark-mode) {
    background-color: #202124;
    color: #bfc2c7;
    border: solid 1px #bfc2c7 !important;
  }

  :global(fieldset.dark-mode) {
    background-color: #202124;
    color: #bfc2c7;
    border: none;
    box-shadow: 0 0 2px #bfc2c7;
  }

  :global(.btn-grad.dark-mode, .btn-grad-embed.dark-mode) {
    box-shadow: none;
  }

  :global(
      .appstore-icon.dark-mode,
      .gh-footer.dark-mode,
      .darkmode-icon-img.dark-mode,
      .loader-container.dark-mode
    ) {
    filter: invert(1);
  }

  :global(input.dark-mode) {
    mix-blend-mode: exclusion;
  }

  :global(select.dark-mode) {
    background-color: #202124 !important;
    color: #bfc2c7 !important;
    mix-blend-mode: exclusion;
  }
</style>
