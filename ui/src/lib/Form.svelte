<script lang="ts">
  import Checkbox from './Checkbox.svelte';
  let username = '';
  let method = 'album';

  $: showTrack = method == 'track';
  $: showAlbum = method != 'artist';

  function collageTypeChanged() {
    console.log(method);
  }
</script>

<form action="/collage" method="get" id="form">
  <label class="form-heading" for="username">Generate a collage for</label>
  <br />
  <input
    bind:value={username}
    required
    autocomplete="on"
    placeholder="*Last.FM username"
    title="Username must be between 2 to 15 characters, begin with a letter and contain only letters, numbers, '_' or '-'"
    pattern={String.raw`^[a-zA-Z][a-zA-Z0-9_\-]{(1, 14)}$`}
    type="text"
    id="username"
    name="username"
  />
  <br />
  <label class="form-heading" for="method">With</label><br />
  <select
    name="method"
    id="method"
    bind:value={method}
    on:change={collageTypeChanged}
  >
    <option value="album">Top Albums</option>
    <option value="artist">Top Artists</option>
    <option value="track">Top Tracks</option></select
  ><br />
  <label class="form-heading" for="period">For the time period</label><br />
  <select name="period" id="period" on:change={collageTypeChanged}>
    <option value="7day">7 Days</option>
    <option value="1month">1 Month</option>
    <option value="3month">3 Months</option>
    <option value="6month">6 Months</option>
    <option value="12month">Year</option>
    <option value="overall">All Time</option></select
  ><br />
  <fieldset id="fieldset">
    <legend class="legend">Collage Options</legend>
    <Checkbox name="Track" visible={showTrack} />
    <Checkbox name="Artist" visible={true} />
    <Checkbox name="Album" visible={showAlbum} />
    <Checkbox name="Playcount" visible={true} />
  </fieldset>
</form>

<style>
  form {
    max-width: 500px;
    padding: 20px;
    background: #fff;
    border-radius: 10px;
    box-shadow: 0px 0px 20px 0px rgba(0, 0, 0, 0.1);
    font-family: 'Poppins';
    margin: auto;
  }
  #username,
  input[type='text'],
  select,
  input[type='number'] {
    appearance: none;
    -moz-appearance: none;
    -webkit-appearance: none;
    width: 100%;
    padding: 12px 20px;
    margin: 8px 0;
    display: inline-block;
    border-radius: 10px;
    box-sizing: border-box;
    font-size: 1em;
    background-color: white;
    background: none;
    color: black;
    font-family: 'Poppins';
    line-height: 20px;
    min-height: 28px;
    border: 2px solid transparent;
    box-shadow: rgb(0 0 0 / 12%) 0px 1px 3px, rgb(0 0 0 / 24%) 0px 1px 2px;
    transition: all 0.1s ease 0s;
  }
  input[type='text'],
  select {
    background: url("data:image/svg+xml,<svg height='10px' width='10px' viewBox='0 0 16 16' fill='%23000000' xmlns='http://www.w3.org/2000/svg'><path d='M7.247 11.14 2.451 5.658C1.885 5.013 2.345 4 3.204 4h9.592a1 1 0 0 1 .753 1.659l-4.796 5.48a1 1 0 0 1-1.506 0z'/></svg>")
      no-repeat;
    background-position: calc(100% - 0.75rem) center !important;
    -moz-appearance: none !important;
    -webkit-appearance: none !important;
    appearance: none !important;
    padding-right: 2rem !important;
  }
  input[type='submit'],
  input[type='button'] {
    font-family: 'Poppins';
    font-weight: bold;
  }
  input[type='submit'],
  input[type='button'] {
    width: 100%;
    background-color: #4caf50;
    color: white;
    padding: 14px 20px;
    margin: 8px 0;
    border: none;
    border-radius: 10px;
    cursor: pointer;
    font-size: 1em;
  }
  input[type='submit']:hover {
    background-color: #45a049;
  }
  input:focus {
    outline: none;
  }
  .legend,
  .form-heading {
    font-size: 1.2em;
    font-weight: bold;
  }
</style>
