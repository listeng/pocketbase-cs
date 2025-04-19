<script>
    import { onDestroy } from "svelte";
    import tooltip from "@/actions/tooltip";
    import { collections } from "@/stores/collections";
    import Draggable from "@/components/base/Draggable.svelte";
    import Field from "@/components/base/Field.svelte";
    import RecordInfo from "@/components/records/RecordInfo.svelte";
    import RecordsPicker from "@/components/records/RecordsPicker.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";

    const batchSize = 100;

    export let field;
    export let value;
    export let picker;

    let fieldRef;
    let list = [];
    let isLoading = false;
    let loadTimeoutId;
    let invalidIds = [];

    $: isMultiple = field.maxSelect > 1;

    $: if (typeof value != "undefined") {
        fieldRef?.changed();
    }

    $: if (needLoad(list, value)) {
        isLoading = true;
        // 将加载函数移到执行队列末尾
        //
        // 这有助于减少布局偏移（关系字段有固定高度的骨架加载器）
        // 并允许其他表单字段更快加载
        clearTimeout(loadTimeoutId);
        loadTimeoutId = setTimeout(load, 0);
    }

    function needLoad() {
        if (isLoading) {
            return false;
        }

        const ids = CommonHelper.toArray(value);

        list = list.filter((item) => ids.includes(item.id));

        return ids.length != list.length;
    }

    async function load() {
        const ids = CommonHelper.toArray(value);

        // 重置
        list = [];
        invalidIds = [];

        if (!field?.collectionId || !ids.length) {
            isLoading = false;
            return;
        }

        isLoading = true;

        let expands = [];
        const presentableRelFields = $collections
            .find((c) => c.id == field.collectionId)
            ?.fields?.filter((f) => !f.hidden && f.presentable && f.type == "relation");
        for (const field of presentableRelFields) {
            expands = expands.concat(CommonHelper.getExpandPresentableRelFields(field, $collections, 2));
        }

        // 批量加载所有选中的记录以避免解析器堆栈溢出错误
        const filterIds = ids.slice();
        const loadPromises = [];
        while (filterIds.length > 0) {
            const filters = [];
            for (const id of filterIds.splice(0, batchSize)) {
                filters.push(`id="${id}"`);
            }

            loadPromises.push(
                ApiClient.collection(field.collectionId).getFullList(batchSize, {
                    filter: filters.join("||"),
                    fields: "*:excerpt(200)",
                    expand: expands.join(","),
                    requestKey: null,
                }),
            );
        }

        try {
            let loadedItems = [];
            await Promise.all(loadPromises).then((values) => {
                loadedItems = loadedItems.concat(...values);
            });

            // 保持选中顺序
            for (const id of ids) {
                const rel = CommonHelper.findByKey(loadedItems, "id", id);
                if (rel) {
                    list.push(rel);
                } else {
                    invalidIds.push(id);
                }
            }

            list = list;

            // 确保在请求期间被删除的任何记录
            // 也会从关系值中移除
            listToValue();
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    function remove(rel) {
        CommonHelper.removeByKey(list, "id", rel.id);
        list = list;

        listToValue();
    }

    function listToValue() {
        if (isMultiple) {
            value = list.map((r) => r.id);
        } else {
            value = list[0]?.id || "";
        }
    }

    onDestroy(() => {
        clearTimeout(loadTimeoutId);
    });
</script>

<Field
    bind:this={fieldRef}
    class="form-field form-field-list {field.required ? 'required' : ''}"
    name={field.name}
    let:uniqueId
>
    <FieldLabel {uniqueId} {field}>
        {#if invalidIds.length}
            <i
                class="ri-error-warning-line link-hint m-l-auto flex-order-10"
                use:tooltip={{
                    position: "left",
                    text:
                        "以下关联ID因缺失或无效已从列表中移除: " +
                        invalidIds.join(", "),
                }}
            />
        {/if}
    </FieldLabel>

    <div class="list">
        <div class="relations-list">
            {#each list as record, i (record.id)}
                <Draggable
                    bind:list
                    group={field.name + "_relation"}
                    index={i}
                    disabled={!isMultiple}
                    let:dragging
                    let:dragover
                    on:sort={() => {
                        listToValue();
                    }}
                >
                    <div class="list-item" class:dragging class:dragover>
                        <div class="content">
                            <RecordInfo {record} />
                        </div>
                        <div class="actions">
                            <button
                                type="button"
                                class="btn btn-transparent btn-hint btn-sm btn-circle btn-remove"
                                use:tooltip={"移除"}
                                on:click={() => remove(record)}
                            >
                                <i class="ri-close-line" />
                            </button>
                        </div>
                    </div>
                </Draggable>
            {:else}
                {#if isLoading}
                    {#each CommonHelper.toArray(value).slice(0, 10) as _}
                        <div class="list-item">
                            <div class="skeleton-loader" />
                        </div>
                    {/each}
                {/if}
            {/each}
        </div>

        <div class="list-item list-item-btn">
            <button
                type="button"
                class="btn btn-transparent btn-sm btn-block"
                on:click={() => picker?.show()}
            >
                <i class="ri-magic-line" />
                <!-- <i class="ri-layout-line" /> -->
                <span class="txt">打开选择器</span>
            </button>
        </div>
    </div>
</Field>

<RecordsPicker
    bind:this={picker}
    {value}
    {field}
    on:save={(e) => {
        list = e.detail || [];
        value = isMultiple ? list.map((r) => r.id) : list[0]?.id || "";
    }}
/>

<style lang="scss">
    .relations-list {
        max-height: 300px;
        overflow: auto; /* 回退 */
        overflow: overlay;
    }
</style>