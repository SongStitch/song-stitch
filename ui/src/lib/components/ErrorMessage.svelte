<script lang="ts">
  export let message: string;
</script>

<div class="error">
  <span>{message}</span>
</div>

<style>
  .error {
    color: red;
    padding-top: 0;
    margin-top: 0;
    font-size: 0.8em;
    padding-left: 0px;
    padding-bottom: 0;
    margin-bottom: 1em;
    contain: inline-size;
  }
</style>
