<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let showModal: boolean;
  export let message: string;

  let dialog: HTMLDialogElement;
  let copyButtonText = 'Copy';

  $: if (dialog && showModal) dialog.showModal();
  const dispatch = createEventDispatcher();

  const copyCode = () => {
    copyButtonText = 'Copied!';
    navigator.clipboard
      .writeText(message)
      .then(
        () => dispatch('copy', message),
        (_) => dispatch('fail')
      )
      .then(() =>
        setTimeout(() => {
          copyButtonText = 'Copy';
        }, 2000)
      );
  };
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<dialog
  bind:this={dialog}
  on:close={() => (showModal = false)}
  on:click|self={() => dialog.close()}
>
  <div on:click|stopPropagation>
    <span class="close" on:click={() => dialog.close()}>&times;</span>
    <slot />
  </div>
  <div class="highlight" id="highlight">
    <button class="copy-code-button" type="button" on:click={copyCode}
      >{copyButtonText}</button
    >
    <pre class="chroma"><code id="embedUrl">{message}</code></pre>
  </div>
</dialog>

<style>
  dialog {
    border: none;
    margin: 20% auto;
    padding: 20px;
    width: 40%;
    border-radius: 10px;
  }
  @media (max-width: 600px) {
    dialog {
      width: 80% !important;
    }
  }
  dialog::backdrop {
    background: rgba(0, 0, 0, 0.3);
  }
  dialog > div {
    padding: 1em;
  }
  dialog[open] {
    animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }
  @keyframes zoom {
    from {
      transform: scale(0.95);
    }
    to {
      transform: scale(1);
    }
  }
  dialog[open]::backdrop {
    animation: fade 0.2s ease-out;
  }
  @keyframes fade {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
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
  .highlight {
    position: relative;
    z-index: 0;
    padding: 0;
    margin: 0;
    border-radius: 4px;
  }
  .highlight > .chroma {
    color: #d0d0d0;
    background-color: black;
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
  :global(body.dark-mode) dialog {
    background-color: #202124;
    color: #bfc2c7;
    border: solid 1px #bfc2c7;
  }
  :global(body.dark-mode) .close:hover,
  :global(body.dark-mode) .close:focus {
    color: white;
  }
</style>
