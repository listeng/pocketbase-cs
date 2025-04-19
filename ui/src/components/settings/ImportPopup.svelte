<script>
    import { createEventDispatcher } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addSuccessToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import CollectionsDiffTable from "@/components/collections/CollectionsDiffTable.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollections = [];
    let newCollections = [];
    let pairs = [];
    let deleteMissing = false;
    let isImporting = false;

    $: if (Array.isArray(oldCollections) && Array.isArray(newCollections)) {
        loadPairs();
    }

    export function show(oldCollectionsArg, newCollectionsArg, deleteMissingArg = false) {
        oldCollections = oldCollectionsArg;
        newCollections = newCollectionsArg;
        deleteMissing = deleteMissingArg;

        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function loadPairs() {
        pairs = [];

        // add modified and deleted (if deleteMissing is set)
        for (const oldCollection of oldCollections) {
            const newCollection = CommonHelper.findByKey(newCollections, "id", oldCollection.id) || null;
            if (
                (deleteMissing && !newCollection?.id) ||
                (newCollection?.id &&
                    CommonHelper.hasCollectionChanges(oldCollection, newCollection, deleteMissing))
            ) {
                pairs.push({
                    old: oldCollection,
                    new: newCollection,
                });
            }
        }

        // add only new collections
        for (const newCollection of newCollections) {
            const oldCollection = CommonHelper.findByKey(oldCollections, "id", newCollection.id) || null;
            if (!oldCollection?.id) {
                pairs.push({
                    old: oldCollection,
                    new: newCollection,
                });
            }
        }
    }

    function submitConfirm() {
        // find deleted fields
        const deletedFieldNames = [];
        if (deleteMissing) {
            for (const old of oldCollections) {
                const imported = CommonHelper.findByKey(newCollections, "id", old.id);
                if (!imported) {
                    // add all fields
                    deletedFieldNames.push(old.name + ".*");
                } else {
                    // add only deleted fields
                    const fields = Array.isArray(old.fields) ? old.fields : [];
                    for (const field of fields) {
                        if (!CommonHelper.findByKey(imported.fields, "id", field.id)) {
                            deletedFieldNames.push(`${old.name}.${field.name} (${field.id})`);
                        }
                    }
                }
            }
        }

        if (deletedFieldNames.length) {
            confirm(
                `你确定要删除以下集合字段及其相关记录数据吗:\n- ${deletedFieldNames.join(
                    "\n- ",
                )}`,
                () => {
                    submit();
                },
            );
        } else {
            submit();
        }
    }

    async function submit() {
        if (isImporting) {
            return;
        }

        isImporting = true;

        try {
            await ApiClient.collections.import(newCollections, deleteMissing);
            addSuccessToast("成功导入集合配置");
            dispatch("submit");
        } catch (err) {
            ApiClient.error(err);
        }

        isImporting = false;

        hide();
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="full-width-popup import-popup"
    overlayClose={false}
    escClose={!isImporting}
    beforeHide={() => !isImporting}
    popup
    on:show
    on:hide
>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">并排对比</h4>
    </svelte:fragment>

    {#each pairs as pair}
        <CollectionsDiffTable collectionA={pair.old} collectionB={pair.new} {deleteMissing} />
    {/each}

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide} disabled={isImporting}>关闭</button
        >
        <button
            type="button"
            class="btn btn-expanded"
            class:btn-loading={isImporting}
            disabled={isImporting}
            on:click={() => submitConfirm()}
        >
            <span class="txt">确认并导入</span>
        </button>
    </svelte:fragment>
</OverlayPanel>