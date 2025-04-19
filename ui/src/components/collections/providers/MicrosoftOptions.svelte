<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
</script>

<div class="section-title">Azure AD 端点</div>
<Field class="form-field required" name="{key}.authURL" let:uniqueId>
    <label for={uniqueId}>认证URL</label>
    <input type="url" id={uniqueId} bind:value={config.authURL} required />
    <div class="help-block">
        例如：{`https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/authorize`}
    </div>
</Field>
<Field class="form-field required" name="{key}.tokenURL" let:uniqueId>
    <label for={uniqueId}>令牌URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenURL} required />
    <div class="help-block">
        例如：{`https://login.microsoftonline.com/YOUR_DIRECTORY_TENANT_ID/oauth2/v2.0/token`}
    </div>
</Field>