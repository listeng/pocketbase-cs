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
            code: 401,
            body: `
                {
                  "code": 401,
                  "message": "请求需要设置有效的记录授权令牌。",
                  "data": {}
                }
            `,
        },
        {
            code: 403,
            body: `
                {
                  "code": 403,
                  "message": "授权的记录模型不允许执行此操作。",
                  "data": {}
                }
            `,
        },
        {
            code: 404,
            body: `
                {
                  "code": 404,
                  "message": "请求的资源未找到。",
                  "data": {}
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">解除 OAuth2 账户关联 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>
        从 <strong>{collection.name}</strong> 记录中解除单个外部 OAuth2 提供者的关联。
    </p>
    <p>只有管理员和账户拥有者可以访问此操作。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${collection?.name}').unlinkExternalAuth(
            pb.authStore.model.id,
            'google'
        );
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '123456');

        await pb.collection('${collection?.name}').unlinkExternalAuth(
          pb.authStore.model.id,
          'google',
        );
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-danger">
    <strong class="label label-primary">DELETE</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong
            >/external-auths/<strong>:provider</strong>
        </p>
    </div>
    <p class="txt-hint txt-sm txt-right">需要 <code>Authorization:TOKEN</code> 头</p>
</div>

<div class="section-title">路径参数</div>
<table class="table-compact table-border m-b-base">
    <thead>
        <tr>
            <th>参数</th>
            <th>类型</th>
            <th width="60%">描述</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>id</td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>授权记录的 ID。</td>
        </tr>
        <tr>
            <td>provider</td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                要解除关联的授权提供者的名称，例如 <code>google</code>、<code>twitter</code>、
                <code>github</code> 等。
            </td>
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