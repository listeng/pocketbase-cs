<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: isAuth = collection.type === "auth";

    $: superusersOnly = collection?.createRule === null;

    $: excludedTableFields = isAuth ? ["password", "verified", "email", "emailVisibility"] : [];

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
                  "message": "创建记录失败。",
                  "data": {
                    "${collection?.fields?.[0]?.name}": {
                      "code": "validation_required",
                      "message": "缺少必填值。"
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
                  "message": "您无权执行此请求。",
                  "data": {}
                }
            `,
        },
    ];

    function getPayload(collection) {
        let payload = CommonHelper.dummyCollectionSchemaData(collection, true);

        if (isAuth) {
            payload.password = "12345678";
            payload.passwordConfirm = "12345678";
            delete payload.verified;
        }

        return payload;
    }
</script>

<h3 class="m-b-sm">创建({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>创建一个新的<strong>{collection.name}</strong>记录。</p>
    <p>
        请求体参数可以以<code>application/json</code>或
        <code>multipart/form-data</code>格式发送。
    </p>
    <p>
        文件上传仅支持<code>multipart/form-data</code>格式。
        <br />
        如需更多信息和示例，请查看详细的
        <a href={import.meta.env.PB_FILE_UPLOAD_DOCS} target="_blank" rel="noopener noreferrer">
            文件上传和处理文档
        </a>。
    </p>
</div>

<!-- prettier-ignore -->
<SdkTabs
    js={`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${backendAbsUrl}');

...

// 示例创建数据
const data = ${JSON.stringify(getPayload(collection), null, 4)};

const record = await pb.collection('${collection?.name}').create(data);
` + (isAuth ?
`
// (可选)发送邮箱验证请求
await pb.collection('${collection?.name}').requestVerification('test@example.com');
` : ""
)}
    dart={`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${backendAbsUrl}');

...

// 示例创建请求体
final body = <String, dynamic>${JSON.stringify(getPayload(collection), null, 2)};

final record = await pb.collection('${collection?.name}').create(body: body);
` + (isAuth ?
`
// (可选)发送邮箱验证请求
await pb.collection('${collection?.name}').requestVerification('test@example.com');
` : ""
)}
/>

<h6 class="m-b-xs">API详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if superusersOnly}
        <p class="txt-hint txt-sm txt-right">需要超级用户<code>Authorization:TOKEN</code>请求头</p>
    {/if}
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
        {#if isAuth}
            <tr>
                <td colspan="3" class="txt-hint txt-bold">认证专用字段</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if collection?.fields?.find((f) => f.name == "email")?.required}
                            <span class="label label-success">必填</span>
                        {:else}
                            <span class="label label-warning">可选</span>
                        {/if}
                        <span>email</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>认证记录邮箱地址。</td>
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
                <td>获取记录数据时是否显示/隐藏认证记录邮箱。</td>
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
                <td>认证记录密码。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-success">必填</span>
                        <span>passwordConfirm</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>认证记录密码确认。</td>
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
                        {#if !field.required || (field.type == "text" && field.autogeneratePattern)}
                            <span class="label label-warning">可选</span>
                        {:else}
                            <span class="label label-success">必填</span>
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
                        {#if field.autogeneratePattern}
                            如果未设置则会自动生成。
                        {/if}
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
                        设置为空值(<code>null</code>、<code>""</code>或<code>[]</code>)以删除已上传的文件。
                    {:else if field.type === "relation"}
                        关联记录{field.maxSelect === 1 ? "ID" : "IDs"}。
                    {/if}
                </td>
            </tr>
        {/each}
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
                返回创建记录时自动展开关联关系。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多6层深度的嵌套关联展开。<br />
                展开的关联关系将被附加到记录的<code>expand</code>属性下(例如<code>{`"expand": {"relField1": {...}, ...}`}</code>)。
                <br />
                只有请求用户有<strong>查看</strong>权限的关联关系才会被展开。
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