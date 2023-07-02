<script lang="ts">
  import Checkbox from './components/Checkbox.svelte';
  import { createForm } from 'felte';
  import { validator } from '@felte/validator-zod';
  import { extender } from '@felte/extender-persist';
  import { z } from 'zod';
  import NumberInput from './components/NumberInput.svelte';

  let method = 'album';
  let maxRows: number;
  let maxColumns: number;
  let showTrack = false;
  let showAlbum = false;
  $: {
    showTrack = method === 'track';
    showAlbum = method !== 'artist';
    maxRows = method === 'track' ? 5 : method === 'artist' ? 10 : 15;
    maxColumns = method === 'track' ? 5 : method === 'artist' ? 10 : 15;
  }
  let showAdvancedOptions = false;
  let showTextSize = false;

  let maintainAspectRatio = false;

  const schema = z.object({
    username: z
      .string()
      .nonempty('Please fill out this field')
      .regex(
        /^[a-zA-Z][a-zA-Z0-9_-]{0,15}$/,
        "Username must be between 2 to 15 characters, begin with a letter and contain only letters, numbers, '_' or '-'"
      ),
    method: z.string().nonempty(),
    period: z.string().nonempty(),
    track: z.boolean().optional(),
    artist: z.boolean().optional(),
    album: z.boolean().optional(),
    playcount: z.boolean().optional(),
    rows: z.number().int().min(1),
    columns: z.number().int().min(1),
    advancedOptions: z.boolean().optional(),
    showTextSize: z.boolean().optional(),
    lossyCompression: z.boolean().optional(),
    textSize: z.string().optional(),
  });

  const { form, errors } = createForm<z.infer<typeof schema>>({
    extend: [validator({ schema }), extender({ id: 'songstitchform' })],
    onSubmit: async (values) => {
      console.log(values);
      const params = new URLSearchParams();
      params.append('username', values.username);
      params.append('method', values.method);
      params.append('period', values.period);
      if (showTrack) params.append('track', values.track.toString());
      params.append('artist', values.artist.toString());
      if (showAlbum) params.append('album', values.album.toString());
      params.append('playcount', values.playcount.toString());
      let rows = values.rows;
      if (rows > maxRows) {
        rows = maxRows;
      }
      params.append('rows', rows.toString());
      let columns = values.columns;
      if (columns > maxColumns) {
        columns = maxColumns;
      }
      params.append('columns', columns.toString());

      if (values.advancedOptions) {
        if (values.showTextSize) {
          params.append('fontsize', values.textSize);
        }
        if (values.lossyCompression) {
          params.append('compress', values.lossyCompression.toString());
        }
      }

      const url = `/collage?${params.toString()}`;
      console.log(url);
      window.open(url, '_self');
    },
  });
</script>

<form use:form on:submit|preventDefault>
  <label class="form-heading" for="username">Generate a collage for</label>
  <br />
  <input
    class="username"
    type="text"
    name="username"
    autocomplete="on"
    placeholder="*Last.FM username"
    style={$errors.username ? 'border-color: #ff0000' : ''}
  />
  {#if $errors.username}
    <div>
      <span class="error">{$errors.username[0]} </span>
    </div>
  {/if}
  <br />
  <label class="form-heading" for="method">With</label>
  <br />
  <select name="method" id="method" bind:value={method}>
    <option value="album">Top Albums</option>
    <option value="artist">Top Artists</option>
    <option value="track">Top Tracks</option></select
  >
  <br />
  <label class="form-heading" for="period">For the time period</label><br />
  <select name="period" id="period">
    <option value="7day">7 Days</option>
    <option value="1month">1 Month</option>
    <option value="3month">3 Months</option>
    <option value="6month">6 Months</option>
    <option value="12month">Year</option>
    <option value="overall">All Time</option></select
  >
  <br />
  <fieldset id="fieldset">
    <legend class="legend">Collage Options</legend>
    <Checkbox text="Display Track Name" visible={showTrack} name="track" />
    <Checkbox text="Display Artist Name" visible={true} name="artist" />
    <Checkbox text="Display Album Name" visible={showAlbum} name="album" />
    <Checkbox text="Display Playcount" visible={true} name="playcount" />
    <br />
    <NumberInput label="Number of Rows" name="rows" max={maxRows} value={3} />
    {#if $errors.rows}
      <p class="error">{$errors.rows[0]}</p>
    {/if}
    <br />
    <NumberInput
      label="Number of Columns"
      name="columns"
      max={maxColumns}
      value={3}
    />
    {#if $errors.columns}
      <p class="error">{$errors.columns[0]}</p>
    {/if}
    <br />
    <Checkbox
      text="Show Advanced Options"
      visible={true}
      name="advancedOptions"
      bind:checked={showAdvancedOptions}
    />
    {#if showAdvancedOptions}
      <div class="advanced-options">
        <Checkbox
          text="Show Text Font Size"
          visible={showAdvancedOptions}
          name="showTextSize"
          bind:checked={showTextSize}
        />
        {#if showTextSize}
          <div id="fontsize-options">
            <label class="advanced-option-label" for="fontsize"
              >Text Font Size</label
            ><br />
            <select name="textSize">
              <option selected value={12}>Small (default)</option>
              <option value={15}>Medium</option>
              <option value={18}>Large</option></select
            ><br />
          </div>
        {/if}
        <Checkbox
          text="Lossy Compress Image"
          visible={showAdvancedOptions}
          name="lossyCompression"
          checked={maintainAspectRatio}
        />
      </div>
    {/if}
  </fieldset>
  <div class="loader-container">
    <div class="loader" />
  </div>
  <input name="submit" class="btn-grad" type="submit" value="Generate" />
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
  input[type='submit'] {
    font-family: 'Poppins';
    font-weight: bold;
  }
  input[type='submit'] {
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
  fieldset {
    border: 1px solid #ccc;
    border-radius: 4px;
    margin-top: 1em;
  }
  .error {
    color: red;
    padding-top: 0;
    margin-top: 0;
    font-size: 0.9em;
    padding-left: 10px;
    padding-bottom: 0;
    margin-bottom: 2px;
  }
  .advanced-options {
    color: darkgrey;
    padding-top: 1em;
    margin-left: 1em;
    padding-bottom: 1em;
  }
  #image-resolution-options {
    margin-left: 1em;
    padding-bottom: 1em;
    padding-top: 1em;
  }
  .advanced-option-label {
    color: black;
    font-size: 1em;
    font-weight: bold;
  }
  .advanced-option-label span {
    color: darkgrey;
  }
  #fontsize-options {
    padding-left: 1em;
    padding-top: 1em;
  }
  .username,
  input[type='text'],
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
</style>
