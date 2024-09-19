<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import Select from "@/components/base/Select.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import { collections } from "@/stores/collections";

    export let field;
    export let key = "";

    const isSingleOptions = [
        { label: "单个", value: true },
        { label: "多个", value: false },
    ];

    const defaultOptions = [
        { label: "禁用", value: false },
        { label: "启用", value: true },
    ];

    let upsertPanel = null;
    let isSingle = field.options?.maxSelect == 1;
    let oldIsSingle = isSingle;

    $: selectCollections = $collections.filter((c) => c.type != "view");

    // load defaults
    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.options.minSelect = null;
            field.options.maxSelect = 1;
        } else {
            field.options.maxSelect = null;
        }
    }

    $: selectedColection = $collections.find((c) => c.id == field.options.collectionId) || null;

    function loadDefaults() {
        field.options = {
            maxSelect: 1,
            collectionId: null,
            cascadeDelete: false,
        };
        isSingle = true;
        oldIsSingle = isSingle;
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment let:interactive>
        <div class="separator" />

        <Field
            class="form-field required {!interactive ? 'readonly' : ''}"
            inlineError
            name="schema.{key}.options.collectionId"
            let:uniqueId
        >
            <ObjectSelect
                id={uniqueId}
                searchable={selectCollections.length > 5}
                selectPlaceholder={"选择数据集 *"}
                noOptionsText="没有数据集"
                selectionKey="id"
                items={selectCollections}
                readonly={!interactive || field.id}
                bind:keyOfSelected={field.options.collectionId}
            >
                <svelte:fragment slot="afterOptions">
                    <hr />
                    <button
                        type="button"
                        class="btn btn-transparent btn-block btn-sm"
                        on:click={() => upsertPanel?.show()}
                    >
                        <i class="ri-add-line" />
                        <span class="txt">新数据集</span>
                    </button>
                </svelte:fragment>
            </ObjectSelect>
        </Field>

        <div class="separator" />

        <Field
            class="form-field form-field-single-multiple-select {!interactive ? 'readonly' : ''}"
            inlineError
            let:uniqueId
        >
            <ObjectSelect
                id={uniqueId}
                items={isSingleOptions}
                readonly={!interactive}
                bind:keyOfSelected={isSingle}
            />
        </Field>

        <div class="separator" />
    </svelte:fragment>

    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            {#if !isSingle}
                <div class="col-sm-6">
                    <Field class="form-field" name="schema.{key}.options.minSelect" let:uniqueId>
                        <label for={uniqueId}>最少选择数量</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            min="1"
                            placeholder="没有最少限制"
                            bind:value={field.options.minSelect}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field class="form-field" name="schema.{key}.options.maxSelect" let:uniqueId>
                        <label for={uniqueId}>最多选择数量</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            placeholder="没有最多限制"
                            min={field.options.minSelect || 2}
                            bind:value={field.options.maxSelect}
                        />
                    </Field>
                </div>
            {/if}

            <div class="col-sm-12">
                <Field class="form-field" name="schema.{key}.options.cascadeDelete" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">级联删除</span>
                        <!-- prettier-ignore -->
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: [
                                    `是否在删除 ${selectedColection?.name || "relation"} 记录时，同时删除当前对应的集合记录。`,
                                    !isSingle ? `对于“多重”关联字段，级联删除仅在从对应记录中移除所有 ${selectedColection?.name || "relation"} ID 时触发。` : null
                                ].filter(Boolean).join("\n\n"),
                                position: "top",
                            }}
                        />
                    </label>
                    <ObjectSelect
                        id={uniqueId}
                        items={defaultOptions}
                        bind:keyOfSelected={field.options.cascadeDelete}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>

<CollectionUpsertPanel
    bind:this={upsertPanel}
    on:save={(e) => {
        if (e?.detail?.collection?.id && e.detail.collection.type != "view") {
            field.options.collectionId = e.detail.collection.id;
        }
    }}
/>
