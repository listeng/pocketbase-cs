<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    }

    function loadDefaults() {
        field.options = {
            convertUrls: false,
        };
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="optionsFooter">
        <Field class="form-field form-field-toggle" name="schema.{key}.options.convertUrls" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={field.options.convertUrls} />
            <label for={uniqueId}>
                <span class="txt">去除URL的域名</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `这可以帮助使编辑器内容在不同环境之间更加便携，因为不再需要替换本地基础 URL`,
                    }}
                />
            </label>
        </Field>
    </svelte:fragment>
</SchemaField>
