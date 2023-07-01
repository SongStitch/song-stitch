<script lang="ts">
  import Header from './lib/Header.svelte';
  import Footer from './lib/Footer.svelte';
</script>

<main>
  <Header />
  <h1>Support</h1>
  <h3 style="text-align: center">Thank you for using SongStitch</h3>
  <h4 style="text-align: center">
    If you have any issues or questions, please contact us at
    <a class="href-links" href="mailto:songstitchsupport@theden.sh"
      >supportsongstitch@theden.sh</a
    >
    or create a
    <a
      class="href-links"
      href="https://github.com/SongStitch/song-stitch/issues">GitHub issue</a
    >
  </h4>
  <Footer />
</main>

<style>
  h1 {
    text-align: center;
    font-size: 2em;
  }
  @media (max-width: 600px) {
    h1 {
      font-size: 1.5em;
    }
  }
</style>
