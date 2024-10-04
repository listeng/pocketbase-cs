<script>
    import { onMount } from "svelte";
    import ApiClient from "@/utils/ApiClient";

    let pages = [];
    let selectedPage = null;
    let iframeUrl = "";

    async function loadPageList() {
        const pageData = await ApiClient.collection("Page").getList(1, 100, {
            filter: "show=true",
            fields: "id,name,remark",
        });

        pages = pageData.items;
    }

    onMount(async () => {
        await loadPageList();
    });

    async function selectPage(page) {
        selectedPage = page.name;
        iframeUrl = "./libs/page/?_=" + Math.random() + "/#/page/" + page.id;
    }
</script>

<aside class="page-sidebar settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">扩展页面</div>

        {#each pages as page}
            <button
                class="sidebar-list-item"
                on:click={() => selectPage(page)}
                class:page-selected={page.name === selectedPage}
            >
                <i class="ri-article-line" aria-hidden="true" />
                <span class="txt">{page.remark}</span>
            </button>
        {/each}
    </div>
</aside>

<div class="page-wrapper flex-content">
    <main class="page-content" style="padding: 0;">
        {#if iframeUrl}
            <iframe class="page" src={iframeUrl} title=""></iframe>
        {:else}
            <div class="no-page-selected">请选择一个页面</div>
        {/if}
    </main>
</div>

<style>
    .page {
        height: 100vh;
        border: 0;
    }

    .page-selected {
        background-color: #007bff;
        color: white;
    }

    .page-selected:hover {
        background-color: #007bff;
        color: white;
    }
</style>
