<script>
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import DynamicOptionsSelect from "@/components/base/DynamicOptionsSelect.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    const isSingleOptions = [
        { label: "单选", value: true },
        { label: "多选", value: false },
    ];

    let isSingle = field.maxSelect <= 1;
    let oldIsSingle = isSingle;

    $: if (typeof field.maxSelect == "undefined") {
        loadDefaults();
    }

    $: if (oldIsSingle != isSingle) {
        oldIsSingle = isSingle;
        if (isSingle) {
            field.maxSelect = 1;
        } else {
            field.maxSelect = field.values?.length || 2;
        }
    }

    function loadDefaults() {
        field.maxSelect = 1;
        field.values = [];
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
            name="fields.{key}.values"
            let:uniqueId
        >
            <DynamicOptionsSelect
                id={uniqueId}
                emptyPlaceholder={"添加选项 *"}
                bind:items={field.values}
            />
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
        {#if !isSingle}
            <Field class="form-field" name="fields.{key}.maxSelect" let:uniqueId>
                <label for={uniqueId}>最多选择</label>
                <input
                    id={uniqueId}
                    type="number"
                    step="1"
                    min="2"
                    max={field.values.length}
                    placeholder="默认为单选"
                    bind:value={field.maxSelect}
                />
            </Field>
        {/if}
    </svelte:fragment>
</SchemaField>