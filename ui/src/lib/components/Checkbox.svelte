<script lang="ts">
  export let visible: boolean = true;
  export let checked: boolean = true;
  export let text: string;
  export let name: string;

  $: color = checked ? 'black' : 'darkgrey';
  $: display = visible ? 'block' : 'none';
</script>

<div class="checkbox-wrapper" style="display: {display}">
  <input checked type="checkbox" class="switch" {name} id={name} />
  <label style="color: {color};" for={name}>{text}</label>
</div>

<style>
  .checkbox-wrapper {
    padding-top: 1em;
  }
  @supports (-webkit-appearance: none) or (-moz-appearance: none) {
    .checkbox-wrapper input[type='checkbox'] {
      --active: #275efe;
      --active-inner: #fff;
      --focus: 2px rgba(39, 94, 254, 0.3);
      --border: #bbc1e1;
      --border-hover: #275efe;
      --background: #fff;
      --disabled: #f6f8ff;
      --disabled-inner: #e1e6f9;
      -webkit-appearance: none;
      -moz-appearance: none;
      height: 21px;
      outline: none;
      display: inline-block;
      vertical-align: top;
      position: relative;
      margin: 0;
      cursor: pointer;
      border: 1px solid var(--bc, var(--border));
      background: var(--b, var(--background));
      transition: background 0.3s, border-color 0.3s, box-shadow 0.2s;
    }
    .checkbox-wrapper input[type='checkbox']:after {
      content: '';
      display: block;
      left: 0;
      top: 0;
      position: absolute;
      transition: transform var(--d-t, 0.3s) var(--d-t-e, ease),
        opacity var(--d-o, 0.2s);
    }
    .checkbox-wrapper input[type='checkbox']:checked {
      --b: var(--active);
      --bc: var(--active);
      --d-o: 0.3s;
      --d-t: 0.6s;
      --d-t-e: cubic-bezier(0.2, 0.85, 0.32, 1.2);
    }
    .checkbox-wrapper input[type='checkbox']:disabled {
      --b: var(--disabled);
      cursor: not-allowed;
      opacity: 0.9;
    }
    .checkbox-wrapper input[type='checkbox']:disabled:checked {
      --b: var(--disabled-inner);
      --bc: var(--border);
    }
    .checkbox-wrapper input[type='checkbox']:disabled + label {
      cursor: not-allowed;
    }
    .checkbox-wrapper
      input[type='checkbox']:hover:not(:checked):not(:disabled) {
      --bc: var(--border-hover);
    }
    .checkbox-wrapper input[type='checkbox']:focus {
      box-shadow: 0 0 0 var(--focus);
    }
    .checkbox-wrapper input[type='checkbox']:not(.switch) {
      width: 21px;
    }
    .checkbox-wrapper input[type='checkbox']:not(.switch):after {
      opacity: var(--o, 0);
    }
    .checkbox-wrapper input[type='checkbox']:not(.switch):checked {
      --o: 1;
    }
    .checkbox-wrapper input[type='checkbox'] + label {
      display: inline-block;
      vertical-align: middle;
      cursor: pointer;
      margin-left: 4px;
    }
    .checkbox-wrapper input[type='checkbox']:not(.switch) {
      border-radius: 7px;
    }
    .checkbox-wrapper input[type='checkbox']:not(.switch):after {
      width: 5px;
      height: 9px;
      border: 2px solid var(--active-inner);
      border-top: 0;
      border-left: 0;
      left: 7px;
      top: 4px;
      transform: rotate(var(--r, 20deg));
    }
    .checkbox-wrapper input[type='checkbox']:not(.switch):checked {
      --r: 43deg;
    }
    .checkbox-wrapper input[type='checkbox'].switch {
      width: 38px;
      border-radius: 11px;
    }
    .checkbox-wrapper input[type='checkbox'].switch:after {
      left: 2px;
      top: 2px;
      border-radius: 50%;
      width: 17px;
      height: 17px;
      background: var(--ab, var(--border));
      transform: translatex(var(--x, 0));
    }
    .checkbox-wrapper input[type='checkbox'].switch:checked {
      --ab: var(--active-inner);
      --x: 17px;
    }
    .checkbox-wrapper
      input[type='checkbox'].switch:disabled:not(:checked):after {
      opacity: 0.6;
    }
  }
  .checkbox-wrapper * {
    box-sizing: inherit;
  }
  .checkbox-wrapper *:before,
  .checkbox-wrapper *:after {
    box-sizing: inherit;
  }
  .checkbox-wrapper * {
    box-sizing: inherit;
  }
  .checkbox-wrapper *:before,
  .checkbox-wrapper *:after {
    box-sizing: inherit;
  }
  .checkbox-wrapper {
    padding-top: 1em;
  }
</style>
