<script>
    import { onMount } from "svelte";
    import ApiClient from "@/utils/ApiClient";

    let pages = [];
    let selectedPage = null;
    let iframeUrl = "";

    function groupItemsByCtype(items) {
        const grouped = items.reduce((acc, item) => {
            const ctype = item.ctype;

            if (!acc[ctype]) {
                acc[ctype] = [];
            }

            acc[ctype].push(item);

            return acc;
        }, {});

        const result = Object.entries(grouped).map(([ctype, items]) => {
            return {
                ctype: ctype,
                pages: items,
            };
        });

        return result;
    }

    async function loadPageList() {
        const pageData = await ApiClient.collection("Page").getList(1, 100, {
            filter: "show=true",
            fields: "id,name,remark,icon,ctype",
            sort: "+sort",
        });

        pages = groupItemsByCtype(pageData.items);
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
        {#each pages as items}
            <div class="sidebar-title">{items.ctype}</div>

            {#each items.pages as page}
                <button
                    class="sidebar-list-item"
                    on:click={() => selectPage(page)}
                    class:page-selected={page.name === selectedPage}
                >
                    <i class={page.icon ? page.icon : "ri-article-line"} aria-hidden="true" />
                    <span class="txt">{page.remark}</span>
                </button>
            {/each}
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
