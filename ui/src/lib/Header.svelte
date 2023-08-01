<script lang="ts">
  import logo from '../assets/images/songstitch_logo.png';
  import darkmodeLogo from '../assets/images/songstitch_logo_dark.png';
  import { darkMode } from '../stores';
</script>

<div class="header-img-container" id="top">
  <a class="img-link" href="/">
    <img
      src={$darkMode ? darkmodeLogo : logo}
      class="header-img"
      alt="SongStitch Logo"
      width="418"
      height="100"
    />
  </a>
</div>
<p class="subheading">
  Share your most played Last.FM music with an image collage!
</p>

<style>
  .header-img-container {
    text-align: center;
    text-decoration: none;
  }
  .header-img {
    border-radius: 0.5em;
    width: 100%;
    max-width: 40vh;
    display: block;
    height: auto;
    display: inline;
  }
  .subheading {
    text-align: center;
    margin-top: 0px;
  }
  a.img-link:active,
  a.img-link:link,
  a.img-link:visited,
  a.img-link:hover {
    color: unset;
    background: none;
  }
  :global(body.dark-node) .subheading,
  .href-links,
  label,
  option,
  .username,
  .number-input,
  .nonbold,
  .reset-text,
  .advanced-option-label {
    background-color: #202124 !important;
    color: #bfc2c7 !important;
  }
</style>
