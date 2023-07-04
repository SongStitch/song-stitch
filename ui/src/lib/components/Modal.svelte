<script>
  export let showModal; // boolean

  let dialog; // HTMLDialogElement

  $: if (dialog && showModal) dialog.showModal();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
<dialog
  bind:this={dialog}
  on:close={() => (showModal = false)}
  on:click|self={() => dialog.close()}
>
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div on:click|stopPropagation>
    <span class="close" on:click={() => dialog.close()}>&times;</span>
    <slot />
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
</style>
