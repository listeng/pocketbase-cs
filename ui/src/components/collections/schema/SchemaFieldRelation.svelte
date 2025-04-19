<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import CollectionUpsertPanel from "@/components/collections/CollectionUpsertPanel.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import { collections } from "@/stores/collections";

    export let field;
    export let key = "";

    const isSingleOptions = [
        { label: "单选", value: true },
        { label: "多选", value: false },
    ];

    const defaultOptions = [
        { label: "否", value: false },
        { label: "是", value: true },
    ];

    let upsertPanel = null;
    let isSingle = field.maxSelect <= 1;
    let oldIsSingle = isSingle;

    $: selectCollections = $collections.filter((c) => !c.system && c.type != "view");

    // load defaults
    $: if (typeof field.maxSelect == "undefined") {
        loadDefaults();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.minSelect = 0;
            field.maxSelect = 1;
        } else {
            field.maxSelect = 999;
        }
    }

    $: selectedColection = $collections.find((c) => c.id == field.collectionId) || null;

    function loadDefaults() {
        field.maxSelect = 1;
        field.collectionId = null;
        field.cascadeDelete = false;
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
            name="fields.{key}.collectionId"
            let:uniqueId
        >
            <ObjectSelect
                id={uniqueId}
                searchable={selectCollections.length > 5}
                selectPlaceholder={"选择集合 *"}
                noOptionsText="未找到集合"
                selectionKey="id"
                items={selectCollections}
                readonly={!interactive || field.id}
                bind:keyOfSelected={field.collectionId}
            >
                <svelte:fragment slot="afterOptions">
                    <hr />
                    <button
                        type="button"
                        class="btn btn-transparent btn-block btn-sm"
                        on:click={() => upsertPanel?.show()}
                    >
                        <i class="ri-add-line" />
                        <span class="txt">新建集合</span>
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
                    <Field class="form-field" name="fields.{key}.minSelect" let:uniqueId>
                        <label for={uniqueId}>最小选择数</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            min="0"
                            placeholder="无最小限制"
                            value={field.minSelect || ""}
                            on:input={(e) => (field.minSelect = e.target.value << 0)}
                        />
                    </Field>
                </div>
                <div class="col-sm-6">
                    <Field class="form-field" name="fields.{key}.maxSelect" let:uniqueId>
                        <label for={uniqueId}>最大选择数</label>
                        <input
                            type="number"
                            id={uniqueId}
                            step="1"
                            placeholder="默认为单选"
                            min={field.minSelect || 1}
                            bind:value={field.maxSelect}
                        />
                    </Field>
                </div>
            {/if}

            <div class="col-sm-12">
                <Field class="form-field" name="fields.{key}.cascadeDelete" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">级联删除</span>
                        <!-- prettier-ignore -->
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: [
                                    `当${selectedColection?.name || "关联"}记录被删除时，是否同时删除当前对应的集合记录。`,
                                    !isSingle ? `对于"多选"关系字段，仅当从对应记录中移除所有${selectedColection?.name || "关联"}ID时才会触发级联删除。` : null
                                ].filter(Boolean).join("\n\n"),
                                position: "top",
                            }}
                        />
                    </label>
                    <ObjectSelect
                        id={uniqueId}
                        items={defaultOptions}
                        bind:keyOfSelected={field.cascadeDelete}
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
            field.collectionId = e.detail.collection.id;
        }
    }}
/>