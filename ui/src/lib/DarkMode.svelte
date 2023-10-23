<script>
  import { onMount } from "svelte";
  import logo from "../assets/images/songstitch_logo.png";
  import logoDark from "../assets/images/songstitch_logo_dark.png";
  import darkModeIcon from "../assets/images/dark_mode.png";
  let darkMode = localStorage.getItem("darkMode") === "true";

  function toggleClassByClassname(className) {
    var elements = document.getElementsByClassName(className);
    for (var i = 0; i < elements.length; i++) {
      elements[i].classList.toggle("dark-mode");
    }
  }

  function applyDarkMode() {
    window.document.body.classList.toggle("dark-mode");
    document.querySelector("html").classList.toggle("dark-mode");
    const headerImg = document.querySelector(".header-img");
    headerImg.classList.toggle("dark-mode");
    headerImg.src = darkMode ? logoDark : logo;
    localStorage.setItem("darkMode", darkMode.toString());
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
  <a href={"#"} class="img-link" on:click={toggle}>
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
  :global(:root.dark-mode) {
    background-color: #202124;
    opacity: 1;
    background-image: radial-gradient(#5d5d5d 0.61px, transparent 0.6px),
      radial-gradient(#5d5d5d 0.6px, #202124 0.6px);
    background-size: 24px 24px;
    background-position:
      0 0,
      12px 12px;
  }

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

  :global(body.dark-mode) .darkmode-icon-img {
    filter: invert(1);
  }
</style>
