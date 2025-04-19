<script>
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let field;
    export let key = "";

    $: if (CommonHelper.isEmpty(field.id)) {
        loadDefaults();
    }

    function loadDefaults() {
        field.cost = 11;
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.min" let:uniqueId>
                    <label for={uniqueId}>最小长度</label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        min="0"
                        placeholder="无最小限制"
                        value={field.min || ""}
                        on:input={(e) => (field.min = e.target.value << 0)}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.max" let:uniqueId>
                    <label for={uniqueId}>最大长度</label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        placeholder="最多71个字符"
                        min={field.min || 0}
                        max="71"
                        value={field.max || ""}
                        on:input={(e) => (field.max = e.target.value << 0)}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.cost" let:uniqueId>
                    <label for={uniqueId}>Bcrypt强度</label>
                    <input
                        type="number"
                        id={uniqueId}
                        placeholder="默认为10"
                        step="1"
                        min="6"
                        max="31"
                        value={field.cost || ""}
                        on:input={(e) => (field.cost = e.target.value << 0)}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.pattern" let:uniqueId>
                    <label for={uniqueId}>验证规则</label>
                    <input type="text" id={uniqueId} placeholder="例如：^\w+$" bind:value={field.pattern} />
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>