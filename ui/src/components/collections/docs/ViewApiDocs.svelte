<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: adminsOnly = collection?.viewRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: if (collection?.id) {
        responses.push({
            code: 200,
            body: JSON.stringify(CommonHelper.dummyCollectionRecord(collection), null, 2),
        });

        if (adminsOnly) {
            responses.push({
                code: 403,
                body: `
                    {
                      "code": 403,
                      "message": "仅管理员可以访问此操作。",
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

<h3 class="m-b-sm">查看 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>获取单个 <strong>{collection.name}</strong> 记录。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        const record = await pb.collection('${collection?.name}').getOne('RECORD_ID', {
            expand: 'relField1,relField2.subRelField',
        });
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        final record = await pb.collection('${collection?.name}').getOne('RECORD_ID',
          expand: 'relField1,relField2.subRelField',
        );
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-info">
    <strong class="label label-primary">GET</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong>
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
            <td>要查看的记录的 ID。</td>
        </tr>
    </tbody>
</table>

<div class="section-title">查询参数</div>
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
            <td>expand</td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                自动展开记录关系。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多 6 层嵌套关系的展开。<br />
                扩展的关系将附加到记录下的
                <code>expand</code> 属性中（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。
                <br />
                只有请求用户有权限<strong>查看</strong>的关系才会被展开。
            </td>
        </tr>
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