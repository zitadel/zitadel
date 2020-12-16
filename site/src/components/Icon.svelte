<script>
  export let name;
  export let size = '20px';
</script>

<style>
  :global(i) {
    font-size: 20px;
    height: 20px;
    position: relative;
    overflow: hidden;
    vertical-align: middle;
	-o-object-fit: contain;
    object-fit: contain;
    -webkit-transform-origin: center center;
    transform-origin: center center;
    color: white;
  }
</style>

<i class={name} style="font-size: {size}; height: {size};"><use xlink:href="#{name}" /></i>
