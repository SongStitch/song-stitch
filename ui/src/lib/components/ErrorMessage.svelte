<script lang="ts">
  export let message: string;
</script>

<div>
  <span class="error">{message} </span>
</div>

<style>
  .error {
    color: red;
    padding-top: 0;
    margin-top: 0;
    font-size: 0.9em;
    padding-left: 10px;
    padding-bottom: 0;
    margin-bottom: 2px;
  }
</style>
