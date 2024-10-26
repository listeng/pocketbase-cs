<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";

    export let collection;

    let responseTab = 204;
    let responses = [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: responses = [
        {
            code: 204,
            body: "null",
        },
        {
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "认证失败。",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "缺少必填值。"
                    }
                  }
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">确认密码重置（{collection.name}）</h3>
<div class="content txt-lg m-b-sm">
    <p>确认<strong>{collection.name}</strong>的密码重置请求并设置新密码。</p>
    <p>
        在此请求之后，所有之前发出的特定记录的令牌将自动失效。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        let oldAuth = pb.authStore.model;

        await pb.collection('${collection?.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );

        // 如有需要，则重新认证
        // （在上述调用之后，所有之前发出的令牌将失效）
        await pb.collection('${collection?.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        final oldAuth = pb.authStore.model;

        await pb.collection('${collection?.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );

        // 如有需要，则重新认证
        // （在上述调用之后，所有之前发出的令牌将失效）
        await pb.collection('${collection?.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/confirm-password-reset
        </p>
    </div>
</div>

<div class="section-title">请求体参数</div>
<table class="table-compact table-border m-b-base">
    <thead>
        <tr>
            <th>参数</th>
            <th>类型</th>
            <th width="50%">描述</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>token</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>来自密码重置请求邮件的令牌。</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>password</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>要设置的新密码。</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>passwordConfirm</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>新密码的确认。</td>
        </tr>
    </tbody>
</table>

<div class="section-title">响应</div>
<div class="tabs">
    <div class="tabs-header compact combined left">
        {#each responses as response (response.code)}
            <button
                class="tab-item"
                class:active={responseTab === response.code}
                on:click={() => (responseTab = response.code)}
            >
                {response.code}
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each responses as response (response.code)}
            <div class="tab-item" class:active={responseTab === response.code}>
                <CodeBlock content={response.body} />
            </div>
        {/each}
    </div>
</div>