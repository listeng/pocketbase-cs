<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: isAuth = collection?.type === "auth";

    $: superusersOnly = collection?.updateRule === null;

    $: excludedTableFields = isAuth ? ["id", "password", "verified", "email", "emailVisibility"] : ["id"];

    $: tableFields =
        collection?.fields?.filter((f) => {
            return !f.hidden && f.type != "autodate" && !excludedTableFields.includes(f.name);
        }) || [];

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(CommonHelper.dummyCollectionRecord(collection), null, 2),
        },
        {
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "Failed to update record.",
                  "data": {
                    "${collection?.fields?.[0]?.name}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `,
        },
        {
            code: 403,
            body: `
                {
                  "status": 403,
                  "message": "You are not allowed to perform this request.",
                  "data": {}
                }
            `,
        },
        {
            code: 404,
            body: `
                {
                  "status": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `,
        },
    ];

    function getPayload(collection) {
        let payload = CommonHelper.dummyCollectionSchemaData(collection, true);

        if (isAuth) {
            payload.oldPassword = "12345678";
            payload.password = "87654321";
            payload.passwordConfirm = "87654321";
            delete payload.verified;
            delete payload.email;
        }

        return payload;
    }
</script>

<h3 class="m-b-sm">更新 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>更新单条 <strong>{collection.name}</strong> 记录。</p>
    <p>
        请求体参数可以以 <code>application/json</code> 或
        <code>multipart/form-data</code> 格式发送。
    </p>
    <p>
        仅支持通过 <code>multipart/form-data</code> 上传文件。
        <br />
        更多信息和示例可查看详细的
        <a href={import.meta.env.PB_FILE_UPLOAD_DOCS} target="_blank" rel="noopener noreferrer">
            文件上传和处理文档
        </a>。
    </p>
    {#if isAuth}
        <p>
            <em>
                注意：如果更改密码，当前记录之前颁发的所有令牌将自动失效，如果您希望用户保持登录状态，需要在更新调用后手动重新认证。
            </em>
        </p>
    {/if}
</div>

<!-- prettier-ignore -->
<SdkTabs
    js={`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${backendAbsUrl}');

...

// 示例更新数据
const data = ${JSON.stringify(getPayload(collection), null, 4)};

const record = await pb.collection('${collection?.name}').update('RECORD_ID', data);
    `}
    dart={`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${backendAbsUrl}');

...

// 示例更新请求体
final body = <String, dynamic>${JSON.stringify(getPayload(collection), null, 2)};

final record = await pb.collection('${collection?.name}').update('RECORD_ID', body: body);
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-warning">
    <strong class="label label-primary">PATCH</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records/<strong>:id</strong>
        </p>
    </div>
    {#if superusersOnly}
        <p class="txt-hint txt-sm txt-right">需要超级用户 <code>Authorization:TOKEN</code> 请求头</p>
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
            <td>要更新的记录ID。</td>
        </tr>
    </tbody>
</table>

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
        {#if isAuth}
            <tr>
                <td colspan="3" class="txt-hint txt-bold">认证特定字段</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>email</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>
                    认证记录的邮箱地址。
                    <br />
                    此字段只能由超级用户或具有"管理"权限的认证记录更新。
                    <br />
                    普通账户可以通过调用"请求更改邮箱"来更新邮箱。
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if collection?.fields?.find((f) => f.name == "emailVisibility")?.required}
                            <span class="label label-success">必填</span>
                        {:else}
                            <span class="label label-warning">可选</span>
                        {/if}
                        <span>emailVisibility</span>
                    </div>
                </td>
                <td>
                    <span class="label">布尔值</span>
                </td>
                <td>获取记录数据时是否显示/隐藏认证记录的邮箱。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>oldPassword</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>
                    旧认证记录密码。
                    <br />
                    仅在更改记录密码时需要此字段。超级用户和具有"管理"权限的认证记录可以跳过此字段。
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>password</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>新认证记录密码。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>passwordConfirm</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>新认证记录密码确认。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>verified</span>
                    </div>
                </td>
                <td>
                    <span class="label">布尔值</span>
                </td>
                <td>
                    表示认证记录是否已验证。
                    <br />
                    此字段只能由超级用户或具有"管理"权限的认证记录设置。
                </td>
            </tr>
            <tr>
                <td colspan="3" class="txt-hint txt-bold">其他字段</td>
            </tr>
        {/if}

        {#each tableFields as field (field.name)}
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if field.required}
                            <span class="label label-success">必填</span>
                        {:else}
                            <span class="label label-warning">可选</span>
                        {/if}
                        <span>{field.name}</span>
                    </div>
                </td>
                <td>
                    <span class="label">{CommonHelper.getFieldValueType(field)}</span>
                </td>
                <td>
                    {#if field.type === "text"}
                        纯文本值。
                    {:else if field.type === "number"}
                        数字值。
                    {:else if field.type === "json"}
                        JSON数组或对象。
                    {:else if field.type === "email"}
                        邮箱地址。
                    {:else if field.type === "url"}
                        URL地址。
                    {:else if field.type === "file"}
                        文件对象。<br />
                        设置为 <code>null</code> 可删除已上传的文件。
                    {:else if field.type === "relation"}
                        关联记录 {field.maxSelect == 1 ? "ID" : "IDs"}。
                    {/if}
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<div class="section-title">查询参数</div>
<table class="table-compact table-border m-b-lg">
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
                返回更新记录时自动展开关联。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField21`} />
                支持最多6层深度的嵌套关联展开。<br />
                展开的关联将附加到记录的 <code>expand</code> 属性下（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。只有用户有<strong>查看</strong>权限的关联才会被展开。
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