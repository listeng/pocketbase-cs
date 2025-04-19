<script>
    import { tick } from "svelte";
    import { querystring } from "svelte-spa-router";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import CollectionDocsPanel from "@/components/collections/CollectionDocsPanel.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import CollectionsSidebar from "@/components/collections/CollectionsSidebar.svelte";
    import RecordPreviewPanel from "@/components/records/RecordPreviewPanel.svelte";
    import RecordUpsertPanel from "@/components/records/RecordUpsertPanel.svelte";
    import RecordsCount from "@/components/records/RecordsCount.svelte";
    import RecordsList from "@/components/records/RecordsList.svelte";
    import { hideControls, pageTitle } from "@/stores/app";
    import {
        activeCollection,
        changeActiveCollectionByIdOrName,
        collections,
        isCollectionsLoading,
        loadCollections,
    } from "@/stores/collections";

    const initialQueryParams = new URLSearchParams($querystring);

    let collectionUpsertPanel;
    let collectionDocsPanel;
    let recordUpsertPanel;
    let recordPreviewPanel;
    let recordsList;
    let recordsCount;
    let filter = initialQueryParams.get("filter") || "";
    let sort = initialQueryParams.get("sort") || "-@rowid";
    let selectedCollectionIdOrName = initialQueryParams.get("collection") || $activeCollection?.id;
    let totalCount = 0; // 用于手动更改计数而无需重新加载recordsCount组件

    loadCollections(selectedCollectionIdOrName);

    $: reactiveParams = new URLSearchParams($querystring);

    $: collectionQueryParam = reactiveParams.get("collection");

    $: if (
        !$isCollectionsLoading &&
        collectionQueryParam &&
        collectionQueryParam != selectedCollectionIdOrName &&
        collectionQueryParam != $activeCollection?.id &&
        collectionQueryParam != $activeCollection?.name
    ) {
        changeActiveCollectionByIdOrName(collectionQueryParam);
    }

    // 在集合变更时重置筛选和排序
    $: if (
        $activeCollection?.id &&
        selectedCollectionIdOrName != $activeCollection.id &&
        selectedCollectionIdOrName != $activeCollection.name
    ) {
        reset();
    }

    $: if ($activeCollection?.id) {
        normalizeSort();
    }

    $: if (!$isCollectionsLoading && initialQueryParams.get("recordId")) {
        showRecordById(initialQueryParams.get("recordId"));
    }

    // 保持URL参数同步
    $: if (!$isCollectionsLoading && (sort || filter || $activeCollection?.id)) {
        updateQueryParams();
    }

    $: $pageTitle = $activeCollection?.name || "集合";

    async function showRecordById(recordId) {
        await tick(); // 确保响应式组件参数已解析

        $activeCollection?.type === "view"
            ? recordPreviewPanel.show(recordId)
            : recordUpsertPanel?.show(recordId);
    }

    function reset() {
        selectedCollectionIdOrName = $activeCollection?.id;
        filter = "";
        sort = "-@rowid";

        normalizeSort();

        updateQueryParams({ recordId: null });

        // 关闭任何打开的集合面板
        collectionUpsertPanel?.forceHide();
        collectionDocsPanel?.hide();
    }

    // 确保排序字段存在于集合中
    async function normalizeSort() {
        if (!sort) {
            return; // 无需规范化
        }

        const collectionFields = CommonHelper.getAllCollectionIdentifiers($activeCollection);

        const sortFields = sort.split(",").map((f) => {
            if (f.startsWith("+") || f.startsWith("-")) {
                return f.substring(1);
            }
            return f;
        });

        // 无效的排序表达式或缺失的排序字段
        if (sortFields.filter((f) => collectionFields.includes(f)).length != sortFields.length) {
            if ($activeCollection?.type != "view") {
                sort = "-@rowid"; // 除视图外的所有集合都有此字段
            } else if (collectionFields.includes("created")) {
                // 常见的自动日期字段
                sort = "-created";
            } else {
                sort = "";
            }
        }
    }

    function updateQueryParams(extra = {}) {
        const queryParams = Object.assign(
            {
                collection: $activeCollection?.id || "",
                filter: filter,
                sort: sort,
            },
            extra,
        );

        CommonHelper.replaceHashQueryParams(queryParams);
    }
</script>

{#if $isCollectionsLoading && !$collections.length}
    <PageWrapper center>
        <div class="placeholder-section m-b-base">
            <span class="loader loader-lg" />
            <h1>正在加载集合...</h1>
        </div>
    </PageWrapper>
{:else if !$collections.length}
    <PageWrapper center>
        <div class="placeholder-section m-b-base">
            <div class="icon">
                <i class="ri-database-2-line" />
            </div>
            {#if $hideControls}
                <h1 class="m-b-10">您还没有任何集合。</h1>
            {:else}
                <h1 class="m-b-10">创建您的第一个集合来添加记录！</h1>
                <button
                    type="button"
                    class="btn btn-expanded-lg btn-lg"
                    on:click={() => collectionUpsertPanel?.show()}
                >
                    <i class="ri-add-line" />
                    <span class="txt">创建新集合</span>
                </button>
            {/if}
        </div>
    </PageWrapper>
{:else}
    <CollectionsSidebar />

    <PageWrapper class="flex-content">
        <header class="page-header">
            <nav class="breadcrumbs">
                <div class="breadcrumb-item">集合</div>
                <div class="breadcrumb-item">{$activeCollection.name}</div>
            </nav>

            <div class="inline-flex gap-5">
                {#if !$hideControls}
                    <button
                        type="button"
                        aria-label="编辑集合"
                        class="btn btn-transparent btn-circle"
                        use:tooltip={{ text: "编辑集合", position: "right" }}
                        on:click={() => collectionUpsertPanel?.show($activeCollection)}
                    >
                        <i class="ri-settings-4-line" />
                    </button>
                {/if}

                <RefreshButton
                    on:refresh={() => {
                        recordsList?.load();
                        recordsCount?.reload();
                    }}
                />
            </div>

            <div class="btns-group">
                <button
                    type="button"
                    class="btn btn-outline"
                    on:click={() => collectionDocsPanel?.show($activeCollection)}
                >
                    <i class="ri-code-s-slash-line" />
                    <span class="txt">API预览</span>
                </button>

                {#if $activeCollection.type !== "view"}
                    <button type="button" class="btn btn-expanded" on:click={() => recordUpsertPanel?.show()}>
                        <i class="ri-add-line" />
                        <span class="txt">新建记录</span>
                    </button>
                {/if}
            </div>
        </header>

        <Searchbar
            value={filter}
            autocompleteCollection={$activeCollection}
            on:submit={(e) => (filter = e.detail)}
        />

        <div class="clearfix m-b-sm" />

        <RecordsList
            bind:this={recordsList}
            collection={$activeCollection}
            bind:filter
            bind:sort
            on:select={(e) => {
                updateQueryParams({
                    recordId: e.detail.id,
                });

                let showModel = e.detail._partial ? e.detail.id : e.detail;

                $activeCollection.type === "view"
                    ? recordPreviewPanel?.show(showModel)
                    : recordUpsertPanel?.show(showModel);
            }}
            on:delete={() => {
                recordsCount?.reload();
            }}
            on:new={() => recordUpsertPanel?.show()}
        />

        <svelte:fragment slot="footer">
            <RecordsCount
                bind:this={recordsCount}
                class="m-r-auto txt-sm txt-hint"
                collection={$activeCollection}
                {filter}
                bind:totalCount
            />
        </svelte:fragment>
    </PageWrapper>
{/if}

<CollectionUpsertPanel
    bind:this={collectionUpsertPanel}
    on:truncate={() => {
        recordsList?.load();
        recordsCount?.reload();
    }}
/>

<CollectionDocsPanel bind:this={collectionDocsPanel} />

<RecordUpsertPanel
    bind:this={recordUpsertPanel}
    collection={$activeCollection}
    on:hide={() => {
        updateQueryParams({ recordId: null });
    }}
    on:save={(e) => {
        if (filter) {
            // 如果有筛选条件，重新加载计数，因为
            // 保存后我们不知道记录是否满足条件
            recordsCount?.reload();
        } else if (e.detail.isNew) {
            totalCount++;
        }

        recordsList?.reloadLoadedPages();
    }}
    on:delete={(e) => {
        if (!filter || recordsList?.hasRecord(e.detail.id)) {
            totalCount--;
        }

        recordsList?.reloadLoadedPages();
    }}
/>

<RecordPreviewPanel
    bind:this={recordPreviewPanel}
    collection={$activeCollection}
    on:hide={() => {
        updateQueryParams({ recordId: null });
    }}
/>