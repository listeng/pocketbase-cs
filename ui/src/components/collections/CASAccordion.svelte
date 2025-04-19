<script>
    import { scale } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import tooltip from "@/actions/tooltip";
    import { errors } from "@/stores/errors";
    import Accordion from "@/components/base/Accordion.svelte";
    import Field from "@/components/base/Field.svelte";
    import ApiClient from "@/utils/ApiClient";

    export let collection;

    $: hasErrors = !CommonHelper.isEmpty($errors?.cas);
    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);
</script>

<Accordion single>
    <svelte:fragment slot="header">
        <div class="inline-flex">
            <i class="ri-key-line"></i>
            <span class="txt"> CAS单点登录 </span>
        </div>

        <div class="flex-fill" />

        {#if collection.cas.enabled}
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

    <div class="grid">
        <Field class="form-field form-field-toggle" name="cas.enabled" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.cas.enabled} />
            <label for={uniqueId}>
                <span class="txt">启用</span>
            </label>
        </Field>

        <Field class="form-field required" name="cas.displayName" let:uniqueId>
            <label for={uniqueId}>显示名称</label>
            <input type="text" id={uniqueId} bind:value={collection.cas.displayName} required />
        </Field>
        
        <div class="section-title">接口端点</div>
        
        <Field class="form-field required" name="cas.loginUrl" let:uniqueId>
            <label for={uniqueId}>Login URL</label>
            <input type="url" id={uniqueId} bind:value={collection.cas.loginUrl} required />
        </Field>
        
        <Field class="form-field required" name="cas.callbackUrl" let:uniqueId>
            <label for={uniqueId}>Callback URL</label>
            <input type="url" id={uniqueId} bind:value={collection.cas.callbackUrl} required />
        
            <div class="help-block">一般是：{backendAbsUrl}/api/collections/(:users)/auth-with-cas</div>
            <div class="help-block">users 可以是 users、_superusers或者某个用户数据集</div>
        </Field>
        
        <Field class="form-field required" name="cas.validateUrl" let:uniqueId>
            <label for={uniqueId}>Validate URL</label>
            <input type="url" id={uniqueId} bind:value={collection.cas.validateUrl} required />
        
            <div class="help-block">这个是给本平台后端调用的地址</div>
        </Field>
        
        <Field class="form-field required" name="cas.realm" let:uniqueId>
            <label for={uniqueId}>邮箱域名</label>
            <input type="text" id={uniqueId} bind:value={collection.cas.realm} required />
        
            <div class="help-block">本平台的用户的电子邮件与此域名相同才是有效用户</div>
        </Field>
        
        <Field class="form-field" name="cas.adminRole" let:uniqueId>
            <label for={uniqueId}>管理员角色名称</label>
            <input type="text" id={uniqueId} bind:value={collection.cas.adminRole} />
        
            <div class="help-block">如果用户是这个角色，则登录的时候设置成管理员</div>
        </Field>
        
        <Field class="form-field" name="cas.logoutUrl" let:uniqueId>
            <label for={uniqueId}>Logout URL</label>
            <input type="url" id={uniqueId} bind:value={collection.cas.logoutUrl} />
        
            <div class="help-block">如果填写此选项，登出本平台的时候将会激活单点登出！</div>
        </Field>
        
        <Field class="form-field" name="cas.onlyCasLogin" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.cas.onlyCasLogin} />
            <label for={uniqueId}>仅使用此登录方法</label>
        
            <div class="help-block">如果启用此选项，登录界面会隐藏普通登录方式！</div>
        </Field>
        
        <Field class="form-field" name="cas.createNewUser" let:uniqueId>
            <input type="checkbox" id={uniqueId} bind:checked={collection.cas.createNewUser} />
            <label for={uniqueId}>是否创建不存在的用户</label>
        
            <div class="help-block">如果启用此选项，通过CAS登录成功后，本平台中不存在的用户，自动创建此用户。</div>
        </Field>
    </div>
</Accordion>