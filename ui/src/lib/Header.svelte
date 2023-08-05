<script>
  import logo from '../assets/images/songstitch_logo.png';
</script>

<div class="header-img-container" id="top">
  <a class="img-link" href="/">
    <img
      src={logo}
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
  :global(body.dark-mode) .subheading {
    background-color: #202124;
    color: #bfc2c7;
  }
</style>
