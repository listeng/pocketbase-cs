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
                    token: "JWT_AUTH_TOKEN",
                    record: CommonHelper.dummyCollectionRecord(collection),
                    meta: {
                        id: "abc123",
                        name: "John Doe",
                        username: "john.doe",
                        email: "test@example.com",
                        avatarUrl: "https://example.com/avatar.png",
                        accessToken: "...",
                        refreshToken: "...",
                        rawUser: {},
                    },
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
                  "message": "提交表单时发生错误。",
                  "data": {
                    "provider": {
                      "code": "validation_required",
                      "message": "缺少必填值。"
                    }
                  }
                }
            `,
        },
    ];
</script>

<h3 class="m-b-sm">使用 OAuth2 认证 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>通过 OAuth2 提供者进行身份验证，并返回新的身份验证令牌和记录数据。</p>
    <p>
        有关更多详细信息，请查看
        <a href={import.meta.env.PB_OAUTH2_EXAMPLE} target="_blank" rel="noopener noreferrer">
            OAuth2 集成文档
        </a>。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        // 通过单个实时调用进行 OAuth2 认证。
        //
        // 确保注册 ${backendAbsUrl}/api/oauth2-redirect 作为重定向 URL。
        const authData = await pb.collection('${collection.name}').authWithOAuth2({ provider: 'google' });

        // 或者通过手动 OAuth2 代码交换进行身份验证
        // const authData = await pb.collection('${collection.name}').authWithOAuth2Code(...);

        // 在上述操作后，您还可以从 authStore 访问身份验证数据
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "登出"最后一次认证的模型
        pb.authStore.clear();
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';
        import 'package:url_launcher/url_launcher.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        // 通过单个实时调用进行 OAuth2 认证。
        //
        // 确保注册 ${backendAbsUrl}/api/oauth2-redirect 作为重定向 URL。
        final authData = await pb.collection('${collection.name}').authWithOAuth2('google', (url) async {
          await launchUrl(url);
        });

        // 或者通过手动 OAuth2 代码交换进行身份验证
        // final authData = await pb.collection('${collection.name}').authWithOAuth2Code(...);

        // 在上述操作后，您还可以从 authStore 访问身份验证数据
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "登出"最后一次认证的模型
        pb.authStore.clear();
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/auth-with-oauth2
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
                    <span>provider</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>OAuth2 客户端提供者的名称（例如 "google"）。</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>code</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>从初始请求返回的授权码。</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>codeVerifier</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>作为 code_challenge 的一部分与初始请求一起发送的代码验证器。</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>redirectUrl</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>与初始请求一起发送的重定向 URL。</td>
        </tr>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-warning">可选</span>
                    <span>createData</span>
                </div>
            </td>
            <td>
                <span class="label">对象</span>
            </td>
            <td>
                <p>在 OAuth2 注册时用于创建身份验证记录的可选数据。</p>
                <p>
                    创建的身份验证记录必须符合常规 <strong>create</strong> 操作中的相同要求和验证。
                    <br />
                    <em>
                        数据只能以 <code>json</code> 格式提交，即 <code>multipart/form-data</code>，并且在 OAuth2 注册过程中不支持文件上传。
                    </em>
                </p>
            </td>
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
                支持最多 6 层嵌套关系的扩展。 <br />
                扩展的关系将附加到记录的 <code>expand</code> 属性下（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。
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