<script>
    import { scale, slide } from "svelte/transition";
    import { errors } from "@/stores/errors";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import Accordion from "@/components/base/Accordion.svelte";

    export let collection;

    $: if (collection.type === "auth" && CommonHelper.isEmpty(collection.options)) {
        collection.options = {
            allowEmailAuth: true,
            allowUsernameAuth: true,
            allowOAuth2Auth: true,
            minPasswordLength: 8,
        };
    }

    $: hasUsernameErrors = false;

    $: hasEmailErrors =
        !CommonHelper.isEmpty($errors?.options?.allowEmailAuth) ||
        !CommonHelper.isEmpty($errors?.options?.onlyEmailDomains) ||
        !CommonHelper.isEmpty($errors?.options?.exceptEmailDomains);

    $: hasOAuth2Errors = !CommonHelper.isEmpty($errors?.options?.allowOAuth2Auth);
</script>

<h4 class="section-title">鉴权方式</h4>

<div class="accordions">
    <Accordion single>
        <svelte:fragment slot="header">
            <div class="inline-flex">
                <i class="ri-user-star-line" />
                <span class="txt">用户名/密码</span>
            </div>

            <div class="flex-fill" />

            {#if collection.options.allowUsernameAuth}
                <span class="label label-success">启用</span>
            {:else}
                <span class="label">禁用</span>
            {/if}

            {#if hasUsernameErrors}
                <i
                    class="ri-error-warning-fill txt-danger"
                    transition:scale={{ duration: 150, start: 0.7 }}
                    use:tooltip={{ text: "有错误", position: "left" }}
                />
            {/if}
        </svelte:fragment>

        <Field class="form-field form-field-toggle m-b-0" name="options.allowUsernameAuth" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.options.allowUsernameAuth} />
            <label for={uniqueId}>启用</label>
        </Field>
    </Accordion>

    <Accordion single>
        <svelte:fragment slot="header">
            <div class="inline-flex">
                <i class="ri-mail-star-line" />
                <span class="txt">邮箱/密码</span>
            </div>

            <div class="flex-fill" />

            {#if collection.options.allowEmailAuth}
                <span class="label label-success">启用</span>
            {:else}
                <span class="label">禁用</span>
            {/if}

            {#if hasEmailErrors}
                <i
                    class="ri-error-warning-fill txt-danger"
                    transition:scale={{ duration: 150, start: 0.7 }}
                    use:tooltip={{ text: "Has errors", position: "left" }}
                />
            {/if}
        </svelte:fragment>

        <Field class="form-field form-field-toggle m-0" name="options.allowEmailAuth" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.options.allowEmailAuth} />
            <label for={uniqueId}>启用</label>
        </Field>

        {#if collection.options.allowEmailAuth}
            <div class="grid grid-sm p-t-sm" transition:slide={{ duration: 150 }}>
                <div class="col-lg-6">
                    <Field
                        class="form-field {!CommonHelper.isEmpty(collection.options.onlyEmailDomains)
                            ? 'disabled'
                            : ''}"
                        name="options.exceptEmailDomains"
                        let:uniqueId
                    >
                        <label for={uniqueId}>
                            <span class="txt">域名黑名单</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: '不允许注册的邮箱域名。\n如果设置了“域名白名单”，此字段将被禁用。',
                                    position: "top",
                                }}
                            />
                        </label>
                        <MultipleValueInput
                            id={uniqueId}
                            disabled={!CommonHelper.isEmpty(collection.options.onlyEmailDomains)}
                            bind:value={collection.options.exceptEmailDomains}
                        />
                        <div class="help-block">使用逗号分隔</div>
                    </Field>
                </div>
                <div class="col-lg-6">
                    <Field
                        class="form-field {!CommonHelper.isEmpty(collection.options.exceptEmailDomains)
                            ? 'disabled'
                            : ''}"
                        name="options.onlyEmailDomains"
                        let:uniqueId
                    >
                        <label for={uniqueId}>
                            <span class="txt">域名白名单</span>
                            <i
                                class="ri-information-line link-hint"
                                use:tooltip={{
                                    text: '仅允许这些邮箱域名注册。\n如果设置了“域名黑名单”，此字段将被禁用。',
                                    position: "top",
                                }}
                            />
                        </label>
                        <MultipleValueInput
                            id={uniqueId}
                            disabled={!CommonHelper.isEmpty(collection.options.exceptEmailDomains)}
                            bind:value={collection.options.onlyEmailDomains}
                        />
                        <div class="help-block">使用逗号分隔</div>
                    </Field>
                </div>
            </div>
        {/if}
    </Accordion>

    <Accordion single>
        <svelte:fragment slot="header">
            <div class="inline-flex">
                <i class="ri-shield-star-line" />
                <span class="txt">OAuth2</span>
            </div>

            <div class="flex-fill" />

            {#if collection.options.allowOAuth2Auth}
                <span class="label label-success">启用</span>
            {:else}
                <span class="label">禁用</span>
            {/if}

            {#if hasOAuth2Errors}
                <i
                    class="ri-error-warning-fill txt-danger"
                    transition:scale={{ duration: 150, start: 0.7 }}
                    use:tooltip={{ text: "Has errors", position: "left" }}
                />
            {/if}
        </svelte:fragment>

        <Field class="form-field form-field-toggle m-b-0" name="options.allowOAuth2Auth" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.options.allowOAuth2Auth} />
            <label for={uniqueId}>启用</label>
        </Field>

        {#if collection.options.allowOAuth2Auth}
            <div class="block" transition:slide={{ duration: 150 }}>
                <div class="flex p-t-base">
                    <a href="#/settings/auth-providers" target="_blank" class="btn btn-sm btn-outline">
                        <span class="txt">管理第三方登录</span>
                    </a>
                </div>
            </div>
        {/if}
    </Accordion>
</div>

<hr />

<h4 class="section-title">常规</h4>

<Field class="form-field required" name="options.minPasswordLength" let:uniqueId>
    <label for={uniqueId}>最小密码长度</label>
    <input
        type="number"
        id={uniqueId}
        required
        min="6"
        max="72"
        bind:value={collection.options.minPasswordLength}
    />
</Field>

<Field class="form-field form-field-toggle m-b-sm" name="options.requireEmail" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={collection.options.requireEmail} />
    <label for={uniqueId}>
        <span class="txt">Email总是必须</span>
        <i
            class="ri-information-line txt-sm link-hint"
            use:tooltip={{
                text: "该约束仅适用于新记录。\n另请注意，一些 OAuth2 提供商（如 Twitter）不会返回电子邮件，如果电子邮件字段是必填项，身份验证可能会失败。",
                position: "right",
            }}
        />
    </label>
</Field>

<Field class="form-field form-field-toggle m-b-sm" name="options.onlyVerified" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={collection.options.onlyVerified} />
    <label for={uniqueId}>
        <span class="txt">禁用未验证用户的权限</span>
        <i
            class="ri-information-line txt-sm link-hint"
            use:tooltip={{
                text: [
                    "如果启用，它会对新的未验证用户身份验证请求返回 403。",
                    "如果您需要更细粒度的控制，请不要启用此选项，而是在您指定的数据集中使用 `@request.auth.verified = true` 规则。",
                ].join("\n"),
                position: "right",
            }}
        />
    </label>
</Field>
