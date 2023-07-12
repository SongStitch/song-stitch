<script>
  import { onMount } from 'svelte';
  import logo from '../assets/images/songstitch_logo.png';
  import logoDark from '../assets/images/songstitch_logo_dark.png';
  import darkModeIcon from '../assets/images/dark_mode.png';
  let darkMode = localStorage.getItem('darkMode') === 'true';

  function toggleClassByClassname(className) {
    var elements = document.getElementsByClassName(className);
    for (var i = 0; i < elements.length; i++) {
      elements[i].classList.toggle('dark-mode');
    }
  }

  function applyDarkMode() {
    document.body.classList.toggle('dark-mode');
    document.querySelector('html').classList.toggle('dark-mode');
    document.getElementsByTagName('dialog')[0].classList.toggle('dark-mode');
    document.getElementsByTagName('fieldset')[0].classList.toggle('dark-mode');
    document.querySelector('#app > main > form').classList.toggle('dark-mode');

    const selectElements = document.getElementsByTagName('select');
    const labelElements = document.getElementsByTagName('label');
    const headerImg = document.querySelector('.header-img');
    headerImg.classList.toggle('dark-mode');
    headerImg.src = darkMode ? logoDark : logo;

    const classNames = [
      'username',
      'number-input',
      'nonbold',
      'href-links',
      'btn-grad',
      'btn-grad-embed',
      'reset-text',
      'appstore-icon',
      'gh-footer',
      'dvanced-option-label',
      'darkmode-icon-img',
      'loader-container',
    ];

    classNames.forEach((className) => toggleClassByClassname(className));

    Array.from(selectElements).forEach((selectElement) => {
      selectElement.classList.toggle('dark-mode');
      const optionElements = selectElement.getElementsByTagName('option');
      Array.from(optionElements).forEach((optionElement) => {
        optionElement.classList.toggle('dark-mode');
      });
    });

    Array.from(labelElements).forEach((labelElement) => {
      labelElement.classList.toggle('dark-mode');
    });

    localStorage.setItem('darkMode', darkMode.toString());
  }

  function toggle() {
    darkMode = !darkMode;
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

  :global(:root.dark-mode) {
    background-color: #1a1b1a;
    background: none;
  }

  :global(
      body.dark-mode,
      .href-links.dark-mode,
      label.dark-mode,
      option.dark-mode,
      .username.dark-mode,
      .number-input.dark-mode,
      .nonbold.dark-mode,
      .reset-text.dark-mode,
      .advanced-option-label.dark-mode
    ) {
    background-color: #1a1b1a !important;
    color: #bfc2c7 !important;
  }

  :global(dialog.dark-mode) {
    background-color: #1a1b1a;
    color: #bfc2c7;
    border: solid 1px #bfc2c7 !important;
  }

  :global(fieldset.dark-mode) {
    background-color: #1a1b1a;
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

  :global(form.dark-mode) {
    background-color: #1a1b1a !important;
    color: #bfc2c7 !important;
    box-shadow: 0 0 2px #bfc2c7;
  }

  :global(input.dark-mode) {
    mix-blend-mode: exclusion;
  }

  :global(select.dark-mode) {
    background-color: #1a1b1a !important;
    color: #bfc2c7 !important;
    mix-blend-mode: exclusion;
  }
</style>
