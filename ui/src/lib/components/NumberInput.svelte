<script lang="ts">
  import ErrorMessage from './ErrorMessage.svelte';

  export let max: number;
  export let name: string;
  export let value: number;
  export let label: string;
  export let errorMessage: string = '';
  let min = 0;

  $: {
    if (value > max) {
      value = max;
    }
    if (value < min) {
      value = min;
    }
  }
</script>

<label class="label" for={name}>
  {label}
  <span class="limit">(max. {max})</span>
</label>
<br />
<input
  class="number-input"
  inputmode="decimal"
  type="number"
  bind:value
  {max}
  {min}
  {name}
  class:error={errorMessage}
/>
{#if errorMessage}
  <ErrorMessage message={errorMessage} />
{:else}
  <div style="display: none;" />
{/if}

<style>
  .number-input {
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
  .label {
    color: black;
    font-size: 1em;
    font-weight: bold;
  }
  .limit {
    color: darkgrey;
  }
  .error {
    border-color: red;
  }
</style>
