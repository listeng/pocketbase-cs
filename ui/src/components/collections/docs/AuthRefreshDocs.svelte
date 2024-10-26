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
            body: JSON.stringify(
                {
                    token: "JWT_TOKEN",
                    record: CommonHelper.dummyCollectionRecord(collection),
                },
                null,
                2
            ),
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
                  "message": "缺少授权记录上下文。",
                  "data": {}
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">身份验证刷新 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>
        为<strong>已经验证的记录</strong>返回新的身份验证响应（令牌和记录数据）。
    </p>
    <p>
        <em>
            此方法通常在页面/屏幕重新加载时由用户调用，以确保之前存储在<code>pb.authStore</code>中的数据仍然有效且是最新的。
        </em>
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        const authData = await pb.collection('${collection?.name}').authRefresh();

        // 以上之后，您还可以从authStore访问刷新后的身份验证数据
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        final authData = await pb.collection('${collection?.name}').authRefresh();

        // 以上之后，您还可以从authStore访问刷新后的身份验证数据
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/auth-refresh
        </p>
    </div>
    <p class="txt-hint txt-sm txt-right">需要记录<code>Authorization:TOKEN</code>头</p>
</div>

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
                自动扩展记录关系。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多6级嵌套关系的扩展。 <br />
                扩展的关系将附加到记录的<code>expand</code>属性下（例如<code>{`"expand": {"relField1": {...}, ...}`}</code>）。
                <br />
                只有请求用户有权限<strong>查看</strong>的关系才会被扩展。
            </td>
        </tr>
        <FieldsQueryParam prefix="record." />
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