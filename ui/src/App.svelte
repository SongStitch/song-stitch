<script lang="ts">
  import Header from './lib/Header.svelte';
  import Footer from './lib/Footer.svelte';
  import Form from './lib/Form.svelte';
  import DarkMode from './lib/DarkMode.svelte';
</script>

<main>
  <Header />
  <DarkMode />
  <Form />
  <Footer />
</main>

<style>
  :global(body.dark-mode) {
    background-color: #202124;
    opacity: 1;
    background-image: radial-gradient(
        #5d5d5d 0.6000000000000001px,
        transparent 0.6000000000000001px
      ),
      radial-gradient(
        #5d5d5d 0.6000000000000001px,
        #202124 0.6000000000000001px
      );
    background-size: 24px 24px;
    background-position: 0 0, 12px 12px;
  }
</style>
