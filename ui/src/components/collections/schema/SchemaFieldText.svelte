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
                    <label for={uniqueId}>
                        <span class="txt">最小长度</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={"清空字段或设置为0表示无限制。"}
                        />
                    </label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        min="0"
                        max={Number.MAX_SAFE_INTEGER}
                        placeholder="无最小限制"
                        value={field.min || ""}
                        on:input={(e) => (field.min = parseInt(e.target.value, 10))}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.max" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">最大长度</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={"清空字段或设置为0将回退到默认限制。"}
                        />
                    </label>
                    <input
                        type="number"
                        id={uniqueId}
                        step="1"
                        placeholder="默认最多5000字符"
                        min={field.min || 0}
                        max={Number.MAX_SAFE_INTEGER}
                        value={field.max || ""}
                        on:input={(e) => (field.max = parseInt(e.target.value, 10))}
                    />
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.pattern" let:uniqueId>
                    <label for={uniqueId}>验证规则</label>
                    <input type="text" id={uniqueId} bind:value={field.pattern} />
                    <div class="help-block">
                        <p>例如 <code>{"^[a-z0-9]+$"}</code></p>
                    </div>
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.autogeneratePattern" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">自动生成规则</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={"设置并在记录创建时自动生成符合规则的文本。"}
                        />
                    </label>
                    <input type="text" id={uniqueId} bind:value={field.autogeneratePattern} />
                    <div class="help-block">
                        <p>例如 <code>{"[a-z0-9]{30}"}</code></p>
                    </div>
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>