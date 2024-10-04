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

    $: allowEmail = collection?.options?.allowEmailAuth;

    $: allowUsername = collection?.options?.allowUsernameAuth;

    $: exampleIdentityLabel =
        allowUsername && allowEmail
            ? "您的用户名或邮箱"
            : allowUsername
            ? "您的用户名"
            : "您的邮箱";

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
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "认证失败。",
                  "data": {
                    "identity": {
                      "code": "validation_required",
                      "message": "缺少必需的值。"
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
        通过
        <strong>
            {#if allowUsername && allowEmail}
                用户名/邮箱
            {:else if allowUsername}
                用户名
            {:else if allowEmail}
                邮箱
            {/if}
        </strong>
        和 <strong>密码</strong> 返回新的认证令牌和账户数据。
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

        // 在上述之后，您还可以从 authStore 访问认证数据
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "注销"最后一个认证的账户
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

        // 在上述之后，您还可以从 authStore 访问认证数据
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "注销"最后一个认证的账户
        pb.authStore.clear();
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/auth-with-password
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
                    <span>identity</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                要认证的记录的
                {#if allowUsername}
                    <strong>username</strong>
                {/if}
                {#if allowUsername && allowEmail}
                    或
                {/if}
                {#if allowEmail}
                    <strong>email</strong>
                {/if}
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
                自动扩展记录关系。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多 6 层嵌套关系的扩展。 <br />
                扩展的关系将附加到记录下的
                <code>expand</code> 属性（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。
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