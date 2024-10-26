<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: responses = [
        {
            code: 200,
            body: `
                [
                    {
                      "id": "8171022dc95a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${collection.id}",
                      "provider": "google",
                      "providerId": "2da15468800514p",
                    },
                    {
                      "id": "171022dc895a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${collection.id}",
                      "provider": "twitter",
                      "providerId": "720688005140514",
                    }
                ]
            `,
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
                  "message": "未找到请求的资源。",
                  "data": {}
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">列出 OAuth2 账户（{collection.name}）</h3>
<div class="content txt-lg m-b-sm">
    <p>
        返回与单个 <strong>{collection.name}</strong> 关联的所有 OAuth2 提供者的列表。
    </p>
    <p>只有管理员和账户拥有者可以访问此操作。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${collection?.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${collection?.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-info">
    <strong class="label label-primary">GET</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong>/external-auths
        </p>
    </div>
    <p class="txt-hint txt-sm txt-right">需要 <code>Authorization:TOKEN</code> 头部</p>
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
    </tbody>
</table>

<div class="section-title">查询参数</div>
<table class="table-compact table-border m-b-base">
    <thead>
        <tr>
            <th>参数</th>
            <th>类型</th>
            <th width="50%">描述</th>
        </tr>
    </thead>
    <tbody>
        <FieldsQueryParam />
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