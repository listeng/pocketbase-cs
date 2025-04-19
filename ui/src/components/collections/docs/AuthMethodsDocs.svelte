<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];
    let authMethods = {};
    let isLoading = false;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: responses = [
        {
            code: 200,
            body: isLoading ? "..." : JSON.stringify(authMethods, null, 2),
        },
        {
            code: 404,
            body: `
                {
                  "status": 404,
                  "message": "Missing collection context.",
                  "data": {}
                }
            `,
        },
    ];

    listAuthMethods();

    async function listAuthMethods() {
        isLoading = true;
        try {
            authMethods = await ApiClient.collection(collection.name).listAuthMethods();
        } catch (err) {
            ApiClient.error(err);
        }
        isLoading = false;
    }
</script>

<h3 class="m-b-sm">列出认证方式 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>返回一个包含所有允许的<strong>{collection.name}</strong>认证方式的公开列表。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        const result = await pb.collection('${collection?.name}').listAuthMethods();
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        final result = await pb.collection('${collection?.name}').listAuthMethods();
    `}
/>

<h6 class="m-b-xs">API详情</h6>
<div class="alert alert-info">
    <strong class="label label-primary">GET</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/auth-methods
        </p>
    </div>
</div>

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