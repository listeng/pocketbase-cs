<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.exceptDomains" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">排除域名</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: '不允许的域名列表。\n如果设置了"仅允许域名"，此字段将被禁用。',
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id={uniqueId}
                        disabled={!CommonHelper.isEmpty(field.onlyDomains)}
                        bind:value={field.exceptDomains}
                    />
                    <div class="help-block">使用逗号分隔。</div>
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="fields.{key}.onlyDomains" let:uniqueId>
                    <label for="{uniqueId}.onlyDomains">
                        <span class="txt">仅允许域名</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: '仅允许的域名列表。\n如果设置了"排除域名"，此字段将被禁用。',
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id="{uniqueId}.onlyDomains"
                        disabled={!CommonHelper.isEmpty(field.exceptDomains)}
                        bind:value={field.onlyDomains}
                    />
                    <div class="help-block">使用逗号分隔。</div>
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>