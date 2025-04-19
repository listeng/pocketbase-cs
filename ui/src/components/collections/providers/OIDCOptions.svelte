<script>
    import { slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import MultipleValueInput from "@/components/base/MultipleValueInput.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let key = "";
    export let config = {};

    const userInfoOptions = [
        { label: "用户信息URL", value: true },
        { label: "ID令牌", value: false },
    ];

    let hasUserInfoURL = !!config.userInfoURL;

    if (CommonHelper.isEmpty(config.pkce)) {
        config.pkce = true;
    }

    if (!config.displayName) {
        config.displayName = "OIDC";
    }

    if (!config.extra) {
        config.extra = {};
        hasUserInfoURL = true;
    }

    $: if (typeof hasUserInfoURL !== undefined) {
        refreshUserInfoState();
    }

    function refreshUserInfoState() {
        if (!hasUserInfoURL) {
            config.userInfoURL = "";
            config.extra = config.extra || {};
        } else {
            config.extra = {};
        }
    }
</script>

<Field class="form-field required" name="{key}.displayName" let:uniqueId>
    <label for={uniqueId}>显示名称</label>
    <input type="text" id={uniqueId} bind:value={config.displayName} required />
</Field>

<div class="section-title">端点</div>

<Field class="form-field required" name="{key}.authURL" let:uniqueId>
    <label for={uniqueId}>认证URL</label>
    <input type="url" id={uniqueId} bind:value={config.authURL} required />
</Field>

<Field class="form-field required" name="{key}.tokenURL" let:uniqueId>
    <label for={uniqueId}>令牌URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenURL} required />
</Field>

<Field class="form-field m-b-xs" let:uniqueId>
    <label for={uniqueId}>获取用户信息方式</label>
    <ObjectSelect id={uniqueId} items={userInfoOptions} bind:keyOfSelected={hasUserInfoURL} />
</Field>

<div class="sub-panel m-b-base">
    {#if hasUserInfoURL}
        <div class="content" transition:slide={{ delay: 10, duration: 150 }}>
            <Field class="form-field required" name="{key}.userInfoURL" let:uniqueId>
                <label for={uniqueId}>用户信息URL</label>
                <input type="url" id={uniqueId} bind:value={config.userInfoURL} required />
            </Field>
        </div>
    {:else}
        <div class="content" transition:slide={{ delay: 10, duration: 150 }}>
            <p class="txt-hint txt-sm m-b-xs">
                <em>
                    这两个字段都被视为可选，因为解析的<code>id_token</code>
                    是可信服务器代码->令牌交换响应的直接结果。
                </em>
            </p>
            <Field class="form-field m-b-xs" name="{key}.extra.jwksURL" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">JWKS验证URL</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "公开令牌验证密钥的URL",
                            position: "top",
                        }}
                    />
                </label>
                <input type="url" id={uniqueId} bind:value={config.extra.jwksURL} />
            </Field>
            <Field class="form-field" name="{key}.extra.issuers" let:uniqueId>
                <label for={uniqueId}>
                    <span class="txt">签发者</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: "用于iss令牌声明验证的接受值列表，以逗号分隔",
                            position: "top",
                        }}
                    />
                </label>
                <MultipleValueInput id={uniqueId} bind:value={config.extra.issuers} />
            </Field>
        </div>
    {/if}
</div>

<Field class="form-field" name="{key}.pkce" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={config.pkce} />
    <label for={uniqueId}>
        <span class="txt">支持PKCE</span>
        <i
            class="ri-information-line link-hint"
            use:tooltip={{
                text: "通常可以安全地始终启用，因为大多数提供商如果不支持PKCE，只会忽略额外的查询参数。",
                position: "right",
            }}
        />
    </label>
</Field>