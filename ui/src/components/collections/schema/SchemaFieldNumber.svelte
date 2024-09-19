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
                <Field class="form-field" name="schema.{key}.options.min" let:uniqueId>
                    <label for={uniqueId}>最小</label>
                    <input type="number" id={uniqueId} bind:value={field.options.min} />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.max" let:uniqueId>
                    <label for={uniqueId}>最大</label>
                    <input
                        type="number"
                        id={uniqueId}
                        min={field.options.min}
                        bind:value={field.options.max}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>

    <svelte:fragment slot="optionsFooter">
        <Field class="form-field form-field-toggle" name="schema.{key}.options.noDecimal" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={field.options.noDecimal} />
            <label for={uniqueId}>
                <span class="txt">忽略小数</span>
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `现有的小数将不受影响`,
                    }}
                />
            </label>
        </Field>
    </svelte:fragment>
</SchemaField>
