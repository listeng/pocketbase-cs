<script>
    import ApiClient from "@/utils/ApiClient";
    import Field from "@/components/base/Field.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let key = "";
    export let config = {};

    $: isRequired = !!config.enabled;
    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    if (!config.displayName) {
        config.displayName = "单点登录";
    }
</script>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.displayName" let:uniqueId>
    <label for={uniqueId}>显示名称</label>
    <input type="text" id={uniqueId} bind:value={config.displayName} required={isRequired} />
</Field>

<div class="section-title">接口端点</div>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.loginUrl" let:uniqueId>
    <label for={uniqueId}>Login URL</label>
    <input type="url" id={uniqueId} bind:value={config.loginUrl} required={isRequired} />
</Field>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.callbackUrl" let:uniqueId>
    <label for={uniqueId}>Callback URL</label>
    <input type="url" id={uniqueId} bind:value={config.callbackUrl} required={isRequired} />

    <div class="help-block">一般是：{backendAbsUrl}api/collections/(:users)/auth-with-cas[u]</div>
    <div class="help-block">users 可以是users或者某个用户数据集</div>
    <div class="help-block">auth-with-cas 是作为管理员登录，auth-with-casu 是作为普通用户登录</div>
</Field>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.validateUrl" let:uniqueId>
    <label for={uniqueId}>Validate URL</label>
    <input type="url" id={uniqueId} bind:value={config.validateUrl} required={isRequired} />

    <div class="help-block">这个是给本平台后端调用的地址</div>
</Field>

<Field class="form-field {isRequired ? 'required' : ''}" name="{key}.realm" let:uniqueId>
    <label for={uniqueId}>邮箱域名</label>
    <input type="text" id={uniqueId} bind:value={config.realm} required={isRequired} />

    <div class="help-block">本平台的用户的电子邮件与此域名相同才是有效用户</div>
</Field>

<Field class="form-field" name="{key}.logoutUrl" let:uniqueId>
    <label for={uniqueId}>Logout URL</label>
    <input type="url" id={uniqueId} bind:value={config.logoutUrl} />

    <div class="help-block">如果填写此选项，登出本平台的时候将会激活单点登出！</div>
</Field>

<Field class="form-field" name="{key}.onlyCasLogin" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={config.onlyCasLogin} />
    <label for={uniqueId}>仅使用此登录方法</label>

    <div class="help-block">如果启用此选项，登录界面会隐藏普通登录方式！</div>
</Field>

<Field class="form-field" name="{key}.createNewUser" let:uniqueId>
    <input type="checkbox" id={uniqueId} bind:checked={config.createNewUser} />
    <label for={uniqueId}>是否创建不存在的用户</label>

    <div class="help-block">如果启用此选项，通过CAS登录成功后，本平台中不存在的用户，自动创建此用户。</div>
</Field>
