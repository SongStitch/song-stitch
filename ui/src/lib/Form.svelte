<script lang="ts">
  import Checkbox from './components/Checkbox.svelte';
  import NumberInput from './components/NumberInput.svelte';
  import ErrorMessage from './components/ErrorMessage.svelte';
  import Modal from './components/Modal.svelte';

  import { createForm } from 'felte';
  import { validator } from '@felte/validator-zod';
  import { extender } from '@felte/extender-persist';
  import { z } from 'zod';

  let showEmbedModal = false;
  let url = '';
  let embedHTML = '';
  let submitting = false;

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
    showBoldtext: z.boolean().optional(),
    textSize: z.string().optional(),
  });

  const generateUrl = (values: z.infer<typeof schema>) => {
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
      if (values.showBoldtext) {
        params.append('boldfont', values.showBoldtext.toString());
      }
    }

    const url = `/collage?${params.toString()}`;
    return url;
  };

  const { form, errors, data, reset } = createForm<z.infer<typeof schema>>({
    extend: [validator({ schema }), extender({ id: 'songstitchform' })],
    onSubmit: async (values) => {
      submitting = true;
      const url = generateUrl(values);
      window.open(url, '_self');
    },
    initialValues: {
      method: 'album',
      period: '7day',
      track: false,
      album: true,
      playcount: true,
      rows: 3,
      columns: 3,
      advancedOptions: false,
      showTextSize: false,
      textSize: '12',
      showBoldtext: false,
      lossyCompression: false,
    },
  });

  const embedOnClick = () => {
    let values = $data;
    url = 'https://songstitch.art' + generateUrl(values);
    embedHTML = `<img class="songstitch-collage" src="${url}">`;
    showEmbedModal = true;
  };

  let maxRows: number;
  let maxColumns: number;
  let showTrack = false;
  let showAlbum = false;
  $: {
    showTrack = $data.method === 'track';
    showAlbum = $data.method !== 'artist';
    maxRows =
      $data.method === 'track' ? 5 : $data.method === 'artist' ? 10 : 15;
    maxColumns =
      $data.method === 'track' ? 5 : $data.method === 'artist' ? 10 : 15;
  }
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
  <select name="method" id="method">
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
    <Checkbox
      text="Display Track Name"
      visible={showTrack}
      name="track"
      bind:checked={$data.track}
    />
    <Checkbox
      text="Display Artist Name"
      visible={true}
      name="artist"
      bind:checked={$data.artist}
    />
    <Checkbox
      text="Display Album Name"
      visible={showAlbum}
      name="album"
      bind:checked={$data.album}
    />
    <Checkbox
      text="Display Playcount"
      visible={true}
      name="playcount"
      bind:checked={$data.playcount}
    />
    <br />
    <NumberInput
      label="Number of Rows"
      name="rows"
      max={maxRows}
      bind:value={$data.rows}
      errorMessage={$errors.rows ? $errors.rows[0] : ''}
    />
    <NumberInput
      label="Number of Columns"
      name="columns"
      max={maxColumns}
      bind:value={$data.columns}
      errorMessage={$errors.columns ? $errors.columns[0] : ''}
    />
    <Checkbox
      text="Show Advanced Options"
      visible={true}
      name="advancedOptions"
      bind:checked={$data.advancedOptions}
    />
    {#if $data.advancedOptions}
      <div class="advanced-options">
        <Checkbox
          text="Use Bold Text"
          name="showBoldtext"
          bind:checked={$data.showBoldtext}
        />
        <Checkbox
          text="Show Text Font Size"
          name="showTextSize"
          bind:checked={$data.showTextSize}
        />
        {#if $data.showTextSize}
          <div id="fontsize-options">
            <label class="advanced-option-label" for="fontsize"
              >Text Font Size</label
            ><br />
            <select name="textSize">
              <option selected value={10}>Extra Small</option>
              <option selected value={12}>Small (default)</option>
              <option value={15}>Medium</option>
              <option value={18}>Large</option></select
            ><br />
          </div>
        {/if}
        <Checkbox
          text="Lossy Compress Image"
          name="lossyCompression"
          bind:checked={$data.lossyCompression}
        />
      </div>
    {/if}
  </fieldset>
  {#if submitting}
    <div class="loader-container">
      <div class="loader" />
    </div>
  {/if}
  <input name="submit" class="btn-grad" type="submit" value="Generate" />
  <input
    name="embed"
    class="btn-grad-embed"
    type="button"
    value="Share/embed"
    on:click={embedOnClick}
  />
</form>
<input type="button" on:click={reset} value="Reset" />
<Modal bind:showModal={showEmbedModal}>
  <div class="modal-text" slot="header">Share/Embed</div>
  <div class="modal-text">
    <a class="href-links" href={url}>Share this link to the collage</a>
    <p>
      Or use this HTML code to embed your configured collage. The latest collage
      will automatically be shown whenever viewed! 🎉
    </p>
    <div class="highlight">
      <button class="copy-code-button" type="button">Copy</button>
      <pre class="chroma"><code id="embedUrl">{embedHTML}</code></pre>
    </div>
  </div>
  <script>
    async function copyCodeToClipboard(button, highlightDiv) {
      const codeToCopy = document.getElementById('embedUrl').innerText;
      try {
        result = await navigator.permissions.query({ name: 'clipboard-write' });
        if (result.state == 'granted' || result.state == 'prompt') {
          await navigator.clipboard.writeText(codeToCopy);
        } else {
          copyCodeBlockExecCommand(codeToCopy, highlightDiv);
        }
      } catch (_) {
        copyCodeBlockExecCommand(codeToCopy, highlightDiv);
      } finally {
        codeWasCopied(button);
      }
    }
    function copyCodeBlockExecCommand(codeToCopy, highlightDiv) {
      const textArea = document.createElement('textArea');
      textArea.contentEditable = 'true';
      textArea.readOnly = 'false';
      textArea.value = codeToCopy;
      highlightDiv.insertBefore(textArea, highlightDiv.firstChild);
      const range = document.createRange();
      range.selectNodeContents(textArea);
      const sel = window.getSelection();
      sel.removeAllRanges();
      sel.addRange(range);
      textArea.setSelectionRange(0, 999999);
      document.execCommand('copy');
      highlightDiv.removeChild(textArea);
    }

    function codeWasCopied(button) {
      button.blur();
      button.innerText = 'Copied!';
      setTimeout(function () {
        button.innerText = 'Copy';
      }, 2000);
    }
    function createCopyButton(highlightDiv) {
      const button = document.getElementsByClassName('copy-code-button')[0];
      document
        .getElementsByClassName('copy-code-button')[0]
        .addEventListener('click', () =>
          copyCodeToClipboard(button, highlightDiv)
        );
      addCopyButtonToDom(button, highlightDiv);
    }
    function addCopyButtonToDom(button, highlightDiv) {
      highlightDiv.insertBefore(button, highlightDiv.firstChild);
      const wrapper = document.createElement('div');
      wrapper.className = 'highlight-wrapper';
      highlightDiv.parentNode.insertBefore(wrapper, highlightDiv);
      wrapper.appendChild(highlightDiv);
    }
    document
      .querySelectorAll('.highlight')
      .forEach((highlightDiv) => createCopyButton(highlightDiv));
  </script>
</Modal>

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
  select {
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
  .advanced-options {
    color: darkgrey;
    padding-top: 1em;
    margin-left: 1em;
    padding-bottom: 1em;
  }
  .advanced-option-label {
    color: black;
    font-size: 1em;
    font-weight: bold;
  }
  #fontsize-options {
    padding-left: 1em;
    padding-top: 1em;
  }
  .username,
  input[type='text'] {
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
  pre {
    overflow-x: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
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
  }
  .modal-text {
    text-align: center;
  }
</style>
