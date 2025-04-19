<script>
    import { scale } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let identityFieldsOptions = [];
    let oldIndexes = "";

    $: isSuperusers = collection?.system && collection?.name === "_superusers";

    $: if (CommonHelper.isEmpty(collection?.passwordAuth)) {
        collection.passwordAuth = {
            enabled: true,
            identityFields: ["email"],
        };
    }

    $: hasErrors = !CommonHelper.isEmpty($errors?.passwordAuth);

    $: if (collection && oldIndexes != collection.indexes.join("")) {
        refreshIdentityFieldsOptions();
    }

    function refreshIdentityFieldsOptions() {
        // email is always available in auth collections
        identityFieldsOptions = [{ value: "email" }];

        const fields = collection?.fields || [];
        const indexes = collection?.indexes || [];

        oldIndexes = indexes.join("");

        for (let idx of indexes) {
            const parsed = CommonHelper.parseIndex(idx);
            if (!parsed.unique || parsed.columns.length != 1 || parsed.columns[0].name == "email") {
                continue;
            }

            const field = fields.find((f) => {
                return !f.hidden && f.name.toLowerCase() == parsed.columns[0].name.toLowerCase();
            });
            if (field) {
                identityFieldsOptions.push({ value: field.name });
            }
        }
    }
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-lock-password-line"></i>
            <span class="txt">身份认证/密码</span>
        </div>

        <div class="flex-fill" />

        {#if collection.passwordAuth.enabled}
            <span class="label label-success">已启用</span>
        {:else}
            <span class="label">已禁用</span>
        {/if}

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "存在错误", position: "left" }}
            />
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle" name="passwordAuth.enabled" let:uniqueId>
        <input
            type="checkbox"
            id={uniqueId}
            bind:checked={collection.passwordAuth.enabled}
            disabled={isSuperusers}
        />
        <label for={uniqueId}>启用</label>
        {#if isSuperusers}
            <i
                class="ri-information-line link-hint"
                use:tooltip={{
                    text: "超级用户必须启用密码认证。",
                    position: "right",
                }}
            />
        {/if}
    </Field>

    <Field class="form-field required m-0" name="passwordAuth.identityFields" let:uniqueId>
        <label for={uniqueId}>
            <span class="txt">唯一身份字段</span>
        </label>
        <ObjectSelect
            items={identityFieldsOptions}
            multiple
            bind:keyOfSelected={collection.passwordAuth.identityFields}
        />
    </Field>
</Accordion>