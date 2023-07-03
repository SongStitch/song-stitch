<script lang="ts">
  import Checkbox from './components/Checkbox.svelte';
  import { createForm } from 'felte';
  import { validator } from '@felte/validator-zod';
  import { extender } from '@felte/extender-persist';
  import { z } from 'zod';
  import NumberInput from './components/NumberInput.svelte';
  import ErrorMessage from './components/ErrorMessage.svelte';

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
      .nonempty('Username is required')
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
    rows: z
      .number({ required_error: 'Number of rows is required' })
      .int()
      .min(1),
    columns: z
      .number({ required_error: 'Number of columns is required' })
      .int()
      .min(1),
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
    style={$errors.username ? 'border: 2px solid red' : ''}
  />
  {#if $errors.username}
    <ErrorMessage message={$errors.username[0]} />
  {/if}
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
    <NumberInput
      label="Number of Rows"
      name="rows"
      max={maxRows}
      value={3}
      errorMessage={$errors.rows ? $errors.rows[0] : ''}
    />
    <NumberInput
      label="Number of Columns"
      name="columns"
      max={maxColumns}
      value={3}
      errorMessage={$errors.columns ? $errors.columns[0] : ''}
    />
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
  <input
    name="embed"
    class="btn-grad-embed"
    type="button"
    value="Share/embed"
  />
</form>
<div id="modal" class="modal">
  <div class="modal-content">
    <div class="modal-text" id="imageUrl" />
    <p class="modal-text">
      Or use this HTML code to embed your configured collage. The latest collage
      will automatically be shown whenever viewed! 🎉
    </p>
    <span class="close">&times;</span>
    <div class="highlight">
      <pre class="chroma"><code id="embedUrl" /></pre>
    </div>
  </div>
</div>

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
  .btn-grad {
    background-image: linear-gradient(
      to right,
      #da22ff 0%,
      #9733ee 51%,
      #da22ff 100%
    );
    margin: 10px;
    padding: 15px 45px;
    text-align: center;
    text-transform: uppercase;
    transition: 0.5s;
    background-size: 200% auto;
    color: white;
    box-shadow: 0 0 20px #eee;
    border-radius: 10px;
    display: block;
  }
  .btn-grad:hover {
    background-position: right center;
    color: #fff;
    text-decoration: none;
  }
  .btn-grad-embed {
    background-image: linear-gradient(
      to right,
      #dd5e89 0%,
      #f7bb97 51%,
      #dd5e89 100%
    );
    margin: 10px;
    padding: 15px 45px;
    text-align: center;
    text-transform: uppercase;
    transition: 0.5s;
    background-size: 200% auto;
    color: white;
    box-shadow: 0 0 20px #eee;
    border-radius: 10px;
    display: block;
  }
  .btn-grad-embed:hover {
    background-position: right center;
    color: #fff;
    text-decoration: none;
  }
  .modal {
    display: none;
    position: fixed;
    z-index: 1;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    overflow: auto;
    background-color: rgb(0, 0, 0);
    background-color: rgba(0, 0, 0, 0.4);
  }
  .modal-content {
    position: relative;
    background-color: #fefefe;
    margin: 20% auto;
    padding: 20px;
    width: 50%;
    border-radius: 4px;
  }
  .modal-text {
    text-align: center;
  }
  .close {
    color: #aaa;
    float: right;
    font-size: 28px;
    font-weight: bold;
    position: absolute;
    top: 0;
    right: 15px;
    transition: 0.3s;
  }
  .close:hover,
  .close:focus {
    color: black;
    text-decoration: none;
    cursor: pointer;
  }
  pre {
    overflow-x: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
  .btn-copy {
    font-size: 1em;
    padding: 10px;
    color: #fff;
    background-color: #4caf50;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }
  .btn-copy:hover {
    background-color: #45a049;
  }
  .highlight-wrapper {
    display: block;
  }
  .highlight {
    position: relative;
    z-index: 0;
    padding: 0;
    margin: 0;
    border-radius: 4px;
  }
  .highlight > .chroma {
    color: #d0d0d0;
    background-color: #212121;
    position: static;
    z-index: 1;
    border-radius: 4px;
    padding: 2em;
  }
  .chroma {
    overflow: auto;
  }
  .chroma .lntable {
    display: table;
    width: 100%;
    padding: 0 0 5px;
    margin: 0;
    border-spacing: 0;
    border: 0;
    overflow: auto;
  }
  .chroma .lntd:first-child {
    padding: 7px 7px 7px 10px;
    margin: 0;
  }
  .chroma .lntd:last-child {
    padding: 7px 10px 7px 7px;
    margin: 0;
  }
  .copy-code-button {
    position: absolute;
    z-index: 2;
    right: 0;
    top: 0;
    font-size: 13px;
    font-weight: 700;
    line-height: 14px;
    width: 65px;
    color: #232326;
    background-color: #b3b3b3;
    border: 1.25px solid #232326;
    border-top-left-radius: 0;
    border-top-right-radius: 4px;
    border-bottom-right-radius: 0;
    border-bottom-left-radius: 4px;
    white-space: nowrap;
    padding: 4px 4px 5px 4px;
    margin: 0 0 0 1px;
    cursor: pointer;
  }
  .copy-code-button:hover,
  .copy-code-button:focus,
  .copy-code-button:active,
  .copy-code-button:active:hover {
    color: #222225;
    background-color: #b3b3b3;
    opacity: 0.8;
  }
  .copyable-text-area {
    position: absolute;
    height: 0;
    z-index: -1;
    opacity: 0.01;
  }
  .loader {
    width: 48px;
    height: 48px;
    border: 5px solid black;
    border-bottom-color: transparent;
    border-radius: 50%;
    display: inline-block;
    box-sizing: border-box;
    animation: rotation 0.6s linear infinite;
    margin-top: 1em;
    display: none;
  }
  @keyframes rotation {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }
  .loader-container {
    display: grid;
    place-items: center;
    display: none;
  }
</style>
