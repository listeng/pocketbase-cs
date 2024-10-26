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
                      "message": "缺少必需的值。"
                    }
                  }
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">确认验证 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>确认 <strong>{collection.name}</strong> 账户验证请求。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').confirmVerification('TOKEN');

        // 可选：使用最新的记录更改刷新之前的 authStore 状态
        await pb.collection('${collection?.name}').authRefresh();
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').confirmVerification('TOKEN');

        // 可选：使用最新的记录更改刷新之前的 authStore 状态
        await pb.collection('${collection?.name}').authRefresh();
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/confirm-verification
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
                    <span class="label label-success">必需</span>
                    <span>token</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>来自验证请求电子邮件的令牌。</td>
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