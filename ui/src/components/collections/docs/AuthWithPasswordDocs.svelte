<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: identityFields = collection?.passwordAuth?.identityFields || [];

    $: exampleIdentityLabel =
        identityFields.length == 0 ? "无" : "您的_" + identityFields.join("_或_").toUpperCase();

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(
                {
                    token: "JWT令牌",
                    record: CommonHelper.dummyCollectionRecord(collection),
                },
                null,
                2,
            ),
        },
        {
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "认证失败。",
                  "data": {
                    "identity": {
                      "code": "validation_required",
                      "message": "缺少必填值。"
                    }
                  }
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">使用密码认证 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>
        使用 <strong>{identityFields.join("/")}</strong> 和 <strong>密码</strong> 组合进行认证。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        const authData = await pb.collection('${collection?.name}').authWithPassword(
            '${exampleIdentityLabel}',
            '您的密码',
        );

        // 认证后，您还可以从authStore访问认证数据
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "登出"
        pb.authStore.clear();
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        final authData = await pb.collection('${collection?.name}').authWithPassword(
          '${exampleIdentityLabel}',
          '您的密码',
        );

        // 认证后，您还可以从authStore访问认证数据
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "登出"
        pb.authStore.clear();
    `}
/>

<h6 class="m-b-xs">API详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/auth-with-password
        </p>
    </div>
</div>

<div class="section-title">请求参数</div>
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
                    <span>identity</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                {#each identityFields as name, i}
                    {#if i > 0}或{/if}
                    <strong>{name}</strong>
                {/each}
                用于认证的记录字段。
            </td>
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
            <td>认证记录的密码。</td>
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
                自动展开记录关联。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多6层深度的嵌套关联展开。<br />
                展开的关联将被附加到记录的 <code>expand</code> 属性下（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。
                <br />
                只有请求用户有<strong>查看</strong>权限的关联才会被展开。
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