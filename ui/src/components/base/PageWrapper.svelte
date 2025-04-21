<script>
    import { superuser } from "@/stores/superuser";

    export let center = false;
    export let showFooter = true;

    let classes = "";
    export { classes as class }; // export reserved keyword
</script>

<div class="page-wrapper {classes}" class:center-content={center}>
    <main class="page-content">
        <slot />
    </main>

    {#if showFooter}
    <footer class="page-footer">
        <slot name="footer" />

        {#if $superuser?.id}
            <a href={import.meta.env.PB_DOCS_URL} target="_blank" rel="noopener noreferrer">
                <i class="ri-book-open-line txt-sm" />
                <span class="txt">文档</span>
            </a>
            <span class="delimiter">|</span>
            <a href={import.meta.env.PB_RELEASES} target="_blank" rel="noopener noreferrer" title="版本发布">
                <span class="txt">PocketBase {import.meta.env.PB_VERSION}</span>
            </a>
        {/if}
    </footer>
    {/if}
</div>