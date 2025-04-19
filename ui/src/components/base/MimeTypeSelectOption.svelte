<script>
    export let item = {}; // {ext, mimeType}
</script>

<span class="txt">{item.ext || "无"}</span>
<small class="txt-hint">{item.mimeType}</small>