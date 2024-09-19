<script>
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import Field from "@/components/base/Field.svelte";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-3">
                <Field class="form-field" name="schema.{key}.options.min" let:uniqueId>
                    <label for={uniqueId}>最小长度</label>
                    <input type="number" id={uniqueId} step="1" min="0" bind:value={field.options.min} />
                </Field>
            </div>

            <div class="col-sm-3">
                <Field class="form-field" name="schema.{key}.options.max" let:uniqueId>
                    <label for={uniqueId}>最大长度</label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        min={field.options.min || 0}
                        bind:value={field.options.max}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.pattern" let:uniqueId>
                    <label for={uniqueId}>正则验证</label>
                    <input
                        type="text"
                        id={uniqueId}
                        placeholder={"使用 Go 表达式来验证, 比如： ^\\w+$"}
                        bind:value={field.options.pattern}
                    />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
