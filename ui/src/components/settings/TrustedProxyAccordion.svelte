<script>
    import tooltip from "@/actions/tooltip";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import { errors } from "@/stores/errors";
    import CommonHelper from "@/utils/CommonHelper";
    import { scale } from "svelte/transition";

    const commonProxyHeaders = ["X-Forwarded-For", "Fly-Client-IP", "CF-Connecting-IP"];

    export let formSettings;
    export let healthData;

    let initialSettingsHash = "";

    $: settingsHash = JSON.stringify(formSettings);

    $: if (initialSettingsHash != settingsHash) {
        initialSettingsHash = settingsHash;
    }

    $: hasChanges = initialSettingsHash != settingsHash;

    $: hasErrors = !CommonHelper.isEmpty($errors?.trustedProxy);

    $: isEnabled = !CommonHelper.isEmpty(formSettings.trustedProxy.headers);

    $: suggestedProxyHeaders = !healthData.possibleProxyHeader
        ? commonProxyHeaders
        : [healthData.possibleProxyHeader].concat(
              commonProxyHeaders.filter((h) => h != healthData.possibleProxyHeader),
          );

    function setHeader(val) {
        formSettings.trustedProxy.headers = [val];
    }

    const ipOptions = [
        { label: "使用最左侧IP", value: true },
        { label: "使用最右侧IP", value: false },
    ];
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-route-line"></i>
            <span class="txt">用户IP代理头</span>
            {#if !isEnabled && healthData.possibleProxyHeader}
                <i
                    class="ri-alert-line txt-sm txt-warning"
                    use:tooltip={"检测到代理头。\n建议将其列为受信任的。"}
                />
            {:else if isEnabled && !hasChanges && !formSettings.trustedProxy.headers.includes(healthData.possibleProxyHeader)}
                <i
                    class="ri-alert-line txt-sm txt-hint"
                    use:tooltip={"配置的代理头与检测到的不匹配。"}
                />
            {/if}
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

    <div class="alert alert-info m-b-sm">
        <div class="content">
            <div class="inline-flex flex-gap-5">
                <span>解析的用户IP:</span>
                <strong>{healthData.realIP || "N/A"}</strong>
                <i
                    class="ri-information-line txt-sm link-hint"
                    use:tooltip={"应显示您的实际IP。\n如果不是，请设置正确的代理头。"}
                />
            </div>
            <br />
            <div class="inline-flex flex-gap-5">
                <span>检测到的代理头:</span>
                <strong>{healthData.possibleProxyHeader || "N/A"}</strong>
            </div>
        </div>
    </div>

    <div class="content m-b-sm">
        <p>
            当PocketBase部署在Fly等平台上或通过NGINX等代理访问时，来自不同用户的请求将源自相同的IP地址（连接到您的PocketBase应用的代理IP）。
        </p>
        <p>
            在这种情况下，为了获取实际的用户IP（用于速率限制、日志记录等），您需要正确配置代理，并在下方列出PocketBase可用于提取用户IP的受信任头。
        </p>
        <p class="txt-bold">使用此类代理时，为避免欺骗，建议：</p>
        <ul class="m-t-0 txt-bold">
            <li>使用仅由代理控制且不能由用户手动设置的头</li>
            <li>确保PocketBase服务器只能通过代理访问</li>
        </ul>
        <p>如果PocketBase未部署在代理后面，可以清除头字段。</p>
    </div>

    <div class="grid grid-sm">
        <div class="col-lg-9">
            <Field class="form-field m-b-0" name="trustedProxy.headers" let:uniqueId>
                <label for={uniqueId}>受信任的代理头</label>
                <MultipleValueInput
                    id={uniqueId}
                    placeholder="留空以禁用"
                    bind:value={formSettings.trustedProxy.headers}
                />
                <div class="form-field-addon">
                    <button
                        type="button"
                        class="btn btn-sm btn-hint btn-transparent btn-clear"
                        class:hidden={CommonHelper.isEmpty(formSettings.trustedProxy.headers)}
                        on:click={() => (formSettings.trustedProxy.headers = [])}
                    >
                        清除
                    </button>
                </div>
                <div class="help-block">
                    <p>
                        逗号分隔的头列表，例如：
                        {#each suggestedProxyHeaders as header}
                            <button
                                type="button"
                                class="label label-sm link-primary txt-mono"
                                on:click={() => setHeader(header)}
                            >
                                {header}
                            </button>&nbsp;
                        {/each}
                    </p>
                </div>
            </Field>
        </div>
        <div class="col-lg-3">
            <Field class="form-field m-0" name="trustedProxy.useLeftmostIP" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">IP优先级选择</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "当代理返回多个IP作为头值时使用。通常认为最右侧的IP更可信，但这可能因代理而异。",
                            position: "right",
                        }}
                    />
                </label>
                <ObjectSelect
                    items={ipOptions}
                    bind:keyOfSelected={formSettings.trustedProxy.useLeftmostIP}
                />
            </Field>
        </div>
    </div>
</Accordion>