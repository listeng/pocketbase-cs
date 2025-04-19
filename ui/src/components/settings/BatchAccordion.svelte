<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    export let formSettings;

    $: hasErrors = !CommonHelper.isEmpty($errors?.batch);

    $: isEnabled = !!formSettings.batch?.enabled;
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-archive-stack-line"></i>
            <span class="txt">批量API</span>
        </div>

        <div class="flex-fill" />

        {#if isEnabled}
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

    <Field class="form-field form-field-toggle m-b-sm" name="batch.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={formSettings.batch.enabled} />
        <label for={uniqueId}>启用 <small class="txt-hint">(实验性功能)</small></label>
    </Field>

    <div class="grid">
        <div class="col-lg-4">
            <Field class="form-field {isEnabled ? 'required' : ''}" name="batch.maxRequests" let:uniqueId>
                <label for={uniqueId}>最大允许批量请求数</label>
                <input
                    type="number"
                    id={uniqueId}
                    min="0"
                    step="1"
                    required={isEnabled}
                    bind:value={formSettings.batch.maxRequests}
                />
            </Field>
        </div>

        <div class="col-lg-4">
            <Field class="form-field {isEnabled ? 'required' : ''}" name="batch.timeout" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">最大处理时间(秒)</span>
                </label>
                <input
                    type="number"
                    id={uniqueId}
                    min="0"
                    step="1"
                    required={isEnabled}
                    bind:value={formSettings.batch.timeout}
                />
            </Field>
        </div>

        <div class="col-lg-4">
            <Field class="form-field" name="batch.maxBodySize" let:uniqueId>
                <label for={uniqueId}>最大请求体大小(字节)</label>
                <input
                    type="number"
                    id={uniqueId}
                    min="0"
                    step="1"
                    placeholder="默认为128MB"
                    value={formSettings.batch.maxBodySize || ""}
                    on:input={(e) => (formSettings.batch.maxBodySize = e.target.value << 0)}
                />
            </Field>
        </div>
    </div>
</Accordion>