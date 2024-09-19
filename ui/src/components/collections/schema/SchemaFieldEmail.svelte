<script>
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <div class="grid grid-sm">
            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.exceptDomains" let:uniqueId>
                    <label for={uniqueId}>
                        <span class="txt">域名黑名单</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: '不允许这些域名。\n如果设置了“域名白名单”，此字段将被禁用。',
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id={uniqueId}
                        disabled={!CommonHelper.isEmpty(field.options.onlyDomains)}
                        bind:value={field.options.exceptDomains}
                    />
                    <div class="help-block">使用逗号分隔</div>
                </Field>
            </div>

            <div class="col-sm-6">
                <Field class="form-field" name="schema.{key}.options.onlyDomains" let:uniqueId>
                    <label for="{uniqueId}.options.onlyDomains">
                        <span class="txt">域名白名单</span>
                        <i
                            class="ri-information-line link-hint"
                            use:tooltip={{
                                text: '仅允许这些域名。\n如果设置了“域名黑名单”，此字段将被禁用。',
                                position: "top",
                            }}
                        />
                    </label>
                    <MultipleValueInput
                        id="{uniqueId}.options.onlyDomains"
                        disabled={!CommonHelper.isEmpty(field.options.exceptDomains)}
                        bind:value={field.options.onlyDomains}
                    />
                    <div class="help-block">使用逗号分隔</div>
                </Field>
            </div>
        </div>
    </svelte:fragment>
</SchemaField>
