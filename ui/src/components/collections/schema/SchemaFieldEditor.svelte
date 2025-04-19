<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <Field class="form-field m-b-sm" name="fields.{key}.maxSize" let:uniqueId>
            <label for={uniqueId}>最大大小 <small>(字节)</small></label>
            <input
                type="number"
                id={uniqueId}
                step="1"
                min="0"
                max={Number.MAX_SAFE_INTEGER}
                value={field.maxSize || ""}
                on:input={(e) => (field.maxSize = parseInt(e.target.value, 10))}
                placeholder="默认最大约5MB"
            />
        </Field>

        <Field class="form-field form-field-toggle" name="fields.{key}.convertURLs" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={field.convertURLs} />
            <label for={uniqueId}>
                <span class="txt">去除URL域名</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `这有助于使编辑器内容在不同环境之间更具可移植性，因为不需要替换本地基础URL。`,
                    }}
                />
            </label>
        </Field>
    </svelte:fragment>
</SchemaField>