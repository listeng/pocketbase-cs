<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.min" let:uniqueId>
                    <label for={uniqueId}>最小值</label>
                    <input type="number" id={uniqueId} bind:value={field.min} />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.max" let:uniqueId>
                    <label for={uniqueId}>最大值</label>
                    <input type="number" id={uniqueId} min={field.min} bind:value={field.max} />
                </Field>
            </div>
        </div>
    </svelte:fragment>

    <svelte:fragment slot="optionsFooter">
        <Field class="form-field form-field-toggle" name="fields.{key}.onlyInt" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={field.onlyInt} />
            <label for={uniqueId}>
                <span class="txt">仅整数</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `已有小数不会受到影响。`,
                    }}
                />
            </label>
        </Field>
    </svelte:fragment>
</SchemaField>