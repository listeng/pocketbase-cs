<script>
    import { scale } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import CommonHelper from "@/utils/CommonHelper";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import AutocompleteInput from "@/components/base/AutocompleteInput.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import { errors, setErrors, removeError } from "@/stores/errors";
    import { collections, loadCollections } from "@/stores/collections";

    export let formSettings;

    const audienceOptions = [
        { value: "", label: "全部" },
        { value: "@guest", label: "仅游客" },
        { value: "@auth", label: "仅认证用户" },
    ];

    const basePredefinedTags = [
        { value: "*:list" },
        { value: "*:view" },
        { value: "*:create" },
        { value: "*:update" },
        { value: "*:delete" },
        { value: "*:file", description: "针对文件下载端点" },
        { value: "*:listAuthMethods" },
        { value: "*:authRefresh" },
        { value: "*:auth", description: "针对所有认证方法" },
        { value: "*:authWithPassword" },
        { value: "*:authWithOAuth2" },
        { value: "*:authWithOTP" },
        { value: "*:requestOTP" },
        { value: "*:requestPasswordReset" },
        { value: "*:confirmPasswordReset" },
        { value: "*:requestVerification" },
        { value: "*:confirmVerification" },
        { value: "*:requestEmailChange" },
        { value: "*:confirmEmailChange" },
    ];

    let predefinedTags = basePredefinedTags;

    let formatInfoPanel;

    $: hasErrors = !CommonHelper.isEmpty($errors?.rateLimits);

    loadPredefinedTags();

    async function loadPredefinedTags() {
        await loadCollections();

        predefinedTags = [];

        for (let collection of $collections) {
            if (collection.system) {
                continue;
            }

            predefinedTags.push({ value: collection.name + ":list" });
            predefinedTags.push({ value: collection.name + ":view" });

            if (collection.type != "view") {
                predefinedTags.push({ value: collection.name + ":create" });
                predefinedTags.push({ value: collection.name + ":update" });
                predefinedTags.push({ value: collection.name + ":delete" });
            }

            if (collection.type == "auth") {
                predefinedTags.push({ value: collection.name + ":listAuthMethods" });
                predefinedTags.push({ value: collection.name + ":authRefresh" });
                predefinedTags.push({ value: collection.name + ":auth" });
                predefinedTags.push({ value: collection.name + ":authWithPassword" });
                predefinedTags.push({ value: collection.name + ":authWithOAuth2" });
                predefinedTags.push({ value: collection.name + ":authWithOTP" });
                predefinedTags.push({ value: collection.name + ":requestOTP" });
                predefinedTags.push({ value: collection.name + ":requestPasswordReset" });
                predefinedTags.push({ value: collection.name + ":confirmPasswordReset" });
                predefinedTags.push({ value: collection.name + ":requestVerification" });
                predefinedTags.push({ value: collection.name + ":confirmVerification" });
                predefinedTags.push({ value: collection.name + ":requestEmailChange" });
                predefinedTags.push({ value: collection.name + ":confirmEmailChange" });
            }

            if (collection.fields.find((f) => f.type == "file")) {
                predefinedTags.push({ value: collection.name + ":file" });
            }
        }

        predefinedTags = predefinedTags.concat(basePredefinedTags);
    }

    function newRule() {
        setErrors({}); // reset

        if (!Array.isArray(formSettings.rateLimits.rules)) {
            formSettings.rateLimits.rules = [];
        }

        formSettings.rateLimits.rules.push({
            label: "",
            maxRequests: 300,
            duration: 10,
            audience: "",
        });

        formSettings.rateLimits.rules = formSettings.rateLimits.rules;

        if (formSettings.rateLimits.rules.length == 1) {
            formSettings.rateLimits.enabled = true;
        }
    }

    function removeRule(i) {
        setErrors({}); // reset

        formSettings.rateLimits.rules.splice(i, 1);
        formSettings.rateLimits.rules = formSettings.rateLimits.rules;

        if (!formSettings.rateLimits.rules.length) {
            formSettings.rateLimits.enabled = false;
        }
    }
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-pulse-fill"></i>
            <span class="txt">速率限制</span>
        </div>

        <div class="flex-fill" />

        {#if hasErrors}
            <i
                class="ri-error-warning-fill txt-danger"
                transition:scale={{ duration: 150, start: 0.7 }}
                use:tooltip={{ text: "存在错误", position: "left" }}
            />
        {/if}

        {#if formSettings.rateLimits.enabled}
            <span class="label label-success">已启用</span>
        {:else}
            <span class="label">已禁用</span>
        {/if}
    </svelte:fragment>

    <Field class="form-field form-field-toggle m-b-xs" name="rateLimits.enabled" let:uniqueId>
        <input type="checkbox" id={uniqueId} bind:checked={formSettings.rateLimits.enabled} />
        <label for={uniqueId}>启用 <small class="txt-hint">(实验性)</small></label>
    </Field>

    {#if !CommonHelper.isEmpty(formSettings.rateLimits.rules)}
        <table class="rate-limit-table">
            <thead>
                <tr>
                    <th class="col-label">速率限制标签</th>
                    <th class="col-requests">最大请求数<br /><small>(每IP)</small></th>
                    <th class="col-duration">时间间隔<br /><small>(秒)</small></th>
                    <th class="col-audience">目标用户</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                {#each formSettings.rateLimits.rules || [] as rule, i}
                    <tr class="rate-limit-row">
                        <td class="col-label">
                            <Field class="form-field" name={"rateLimits.rules." + i + ".label"} inlineError>
                                <AutocompleteInput
                                    required
                                    placeholder="标签 (users:create) 或路径 (/api/)"
                                    options={predefinedTags}
                                    bind:value={rule.label}
                                />
                            </Field>
                        </td>
                        <td class="col-requests">
                            <Field
                                class="form-field"
                                name={"rateLimits.rules." + i + ".maxRequests"}
                                inlineError
                            >
                                <input
                                    type="number"
                                    required
                                    placeholder="最大请求数*"
                                    min="1"
                                    step="1"
                                    bind:value={rule.maxRequests}
                                />
                            </Field>
                        </td>
                        <td class="col-duration">
                            <Field
                                class="form-field"
                                name={"rateLimits.rules." + i + ".duration"}
                                inlineError
                            >
                                <input
                                    type="number"
                                    required
                                    placeholder="时间间隔*"
                                    min="1"
                                    step="1"
                                    bind:value={rule.duration}
                                />
                            </Field>
                        </td>
                        <td class="col-audience">
                            <Field
                                class="form-field"
                                name={"rateLimits.rules." + i + ".audience"}
                                inlineError
                            >
                                <ObjectSelect
                                    items={audienceOptions}
                                    bind:keyOfSelected={rule.audience}
                                    on:change={() => {
                                        removeError("rateLimits.rules." + i); // reset rule errors
                                    }}
                                />
                            </Field>
                        </td>
                        <td class="col-action">
                            <button
                                type="button"
                                title="移除规则"
                                aria-label="移除规则"
                                class="btn btn-xs btn-circle btn-hint btn-transparent"
                                on:click={() => removeRule(i)}
                            >
                                <i class="ri-close-line"></i>
                            </button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {/if}

    <div class="flex m-t-sm">
        <button
            type="button"
            class="btn btn-sm btn-secondary m-r-auto"
            class:btn-danger={$errors?.rateLimits?.rules?.message}
            on:click={() => newRule()}
        >
            <i class="ri-add-line"></i>
            <span class="txt">添加速率限制规则</span>
        </button>

        <button type="button" class="txt-nowrap txt-sm link-hint" on:click={() => formatInfoPanel?.show()}>
            <em>了解更多关于速率限制规则的信息</em>
        </button>
    </div>
</Accordion>

<OverlayPanel bind:this={formatInfoPanel}>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">速率限制标签格式</h4>
    </svelte:fragment>

    <p>速率限制规则按以下顺序解析(在第一个匹配项处停止):</p>
    <ol>
        <li>精确标签(例如 <code>users:create</code>)</li>
        <li>通配符标签(例如 <code>*:create</code>)</li>
        <li>方法 + 精确路径(例如 <code>POST /a/b</code>)</li>
        <li>方法 + 前缀路径(例如 <code>POST /a/b<strong>/</strong></code>)</li>
        <li>精确路径(例如 <code>/a/b</code>)</li>
        <li>前缀路径(例如 <code>/a/b<strong>/</strong></code>)</li>
    </ol>
    <p>
        如果有多个相同标签但目标用户群体不同的规则(例如"游客"与"认证用户"),只考虑匹配的群体规则。
    </p>

    <hr class="m-t-xs m-b-xs" />

    <p>速率限制标签可以是以下格式之一:</p>
    <ul>
        <li class="m-b-sm">
            <code>[方法 ]/my/path</code> - 完全匹配路由(<strong>不能有尾部斜杠</strong>;"方法"可选)。
            <br /> 例如:
            <ul class="m-0">
                <li class="m-0">
                    <code>/hello</code> - 匹配 <code>GET /hello</code>, <code>POST /hello</code> 等
                </li>
                <li class="m-0">
                    <code>POST /hello</code> - 仅匹配 <code>POST /hello</code>
                </li>
            </ul>
        </li>
        <li class="m-b-sm">
            <code>[方法 ]/my/prefix<strong>/</strong></code> - 路径前缀(<strong>必须以斜杠结尾</strong>;"方法"可选)。例如:
            <ul class="m-0">
                <li class="m-0">
                    <code>/hello/</code> - 匹配 <code>GET /hello</code>,
                    <code>POST /hello/a/b/c</code> 等
                </li>
                <li class="m-0">
                    <code>POST /hello/</code> - 匹配 <code>POST /hello</code>,
                    <code>POST /hello/a/b/c</code> 等
                </li>
            </ul>
        </li>
        <li>
            <code>collectionName:predefinedTag</code> - 针对单个集合的特定操作。要对所有集合应用规则，可以使用<code>*</code>通配符。例如:
            <code>posts:create</code>, <code>users:listAuthMethods</code>, <code>*:auth</code>。
            <br />
            预定义的集合标签有(<em>开始输入时会有自动补全</em>):
            <ul>
                {#each basePredefinedTags as tag}
                    <li class="m-0">
                        {tag.value.replace("*:", ":")}
                        {#if tag.description}
                            <em class="txt-hint">({tag.description})</em>
                        {/if}
                    </li>
                {/each}
            </ul>
        </li>
    </ul>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => formatInfoPanel?.hide()}>
            关闭
        </button>
    </svelte:fragment>
</OverlayPanel>