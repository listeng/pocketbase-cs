<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";

    export let collection;

    let responseTab = 204;
    let responses = [];

    $: adminsOnly = collection?.deleteRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: if (collection?.id) {
        responses.push({
            code: 204,
            body: `
                null
            `,
        });

        responses.push({
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "删除记录失败。确保该记录不是必需关系引用的一部分。",
                  "data": {}
                }
            `,
        });

        if (adminsOnly) {
            responses.push({
                code: 403,
                body: `
                    {
                      "code": 403,
                      "message": "只有管理员可以访问此操作。",
                      "data": {}
                    }
                `,
            });
        }

        responses.push({
            code: 404,
            body: `
                {
                  "code": 404,
                  "message": "请求的资源未找到。",
                  "data": {}
                }
            `,
        });
    }
</script>

<h3 class="m-b-sm">删除 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>删除单个 <strong>{collection.name}</strong> 记录。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').delete('RECORD_ID');
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').delete('RECORD_ID');
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-danger">
    <strong class="label label-primary">DELETE</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong>
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">需要管理员 <code>Authorization:TOKEN</code> 头</p>
    {/if}
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
            <td>要删除的记录的 ID。</td>
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