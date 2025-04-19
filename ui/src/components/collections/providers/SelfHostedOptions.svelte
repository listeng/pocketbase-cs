<script>
    import Field from "@/components/base/Field.svelte";

    export let key = "";
    export let config = {};
    export let required = false;
    export let title = "提供商端点";

    $: isRequired = required && config?.enabled;
</script>

<div class="section-title">{title}</div>
<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.authURL" let:uniqueId>
    <label for={uniqueId}>认证URL</label>
    <input type="url" id={uniqueId} bind:value={config.authURL} required={isRequired} />
</Field>
<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.tokenURL" let:uniqueId>
    <label for={uniqueId}>令牌URL</label>
    <input type="url" id={uniqueId} bind:value={config.tokenURL} required={isRequired} />
</Field>
<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.userInfoURL" let:uniqueId>
    <label for={uniqueId}>用户信息URL</label>
    <input type="url" id={uniqueId} bind:value={config.userInfoURL} required={isRequired} />
</Field>