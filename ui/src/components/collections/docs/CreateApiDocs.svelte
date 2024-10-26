<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";

    export let collection;

    let responseTab = 200;
    let responses = [];
    let baseData = {};

    $: isAuth = collection.type === "auth";

    $: adminsOnly = collection?.createRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(CommonHelper.dummyCollectionRecord(collection), null, 2),
        },
        {
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "创建记录失败。",
                  "data": {
                    "${collection?.schema?.[0]?.name}": {
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
                  "code": 403,
                  "message": "您没有权限执行此请求。",
                  "data": {}
                }
            `,
        },
    ];

    $: if (isAuth) {
        baseData = {
            username: "test_username",
            email: "test@example.com",
            emailVisibility: true,
            password: "12345678",
            passwordConfirm: "12345678",
        };
    } else {
        baseData = {};
    }
</script>

<h3 class="m-b-sm">创建 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>创建一个新的 <strong>{collection.name}</strong> 记录。</p>
    <p>
        请求体参数可以作为 <code>application/json</code> 或
        <code>multipart/form-data</code> 发送。
    </p>
    <p>
        文件上传仅支持 <code>multipart/form-data</code>。
        <br />
        有关更多信息和示例，请查看详细的
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
const data = ${JSON.stringify(Object.assign({}, baseData, CommonHelper.dummyCollectionSchemaData(collection)), null, 4)};

const record = await pb.collection('${collection?.name}').create(data);
` + (isAuth ?
`
// （可选）发送电子邮件验证请求
await pb.collection('${collection?.name}').requestVerification('test@example.com');
` : ""
)}
    dart={`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${backendAbsUrl}');

...

// 示例创建请求体
final body = <String, dynamic>${JSON.stringify(Object.assign({}, baseData, CommonHelper.dummyCollectionSchemaData(collection)), null, 2)};

final record = await pb.collection('${collection?.name}').create(body: body);
` + (isAuth ?
`
// （可选）发送电子邮件验证请求
await pb.collection('${collection?.name}').requestVerification('test@example.com');
` : ""
)}
/>

<h6 class="m-b-xs">API 详情</h6>
<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">需要管理员 <code>Authorization:TOKEN</code> 头</p>
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
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-warning">可选</span>
                    <span>id</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                <strong>15 个字符的字符串</strong> 用于存储记录 ID。
                <br />
                如果未设置，将自动生成。
            </td>
        </tr>

        {#if isAuth}
            <tr>
                <td colspan="3" class="txt-hint">认证字段</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>username</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>
                    认证记录的用户名。
                    <br />
                    如果未设置，将自动生成。
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        {#if collection?.options?.requireEmail}
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
                <td>认证记录的电子邮件地址。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>emailVisibility</span>
                    </div>
                </td>
                <td>
                    <span class="label">布尔</span>
                </td>
                <td>在获取记录数据时是否显示/隐藏认证记录的电子邮件。</td>
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
                <td>认证记录的密码确认。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>verified</span>
                    </div>
                </td>
                <td>
                    <span class="label">布尔</span>
                </td>
                <td>
                    指示认证记录是否已验证。
                    <br />
                    此字段只能由管理员或具有“管理”访问权限的认证记录设置。
                </td>
            </tr>
            <tr>
                <td colspan="3" class="txt-hint">模式字段</td>
            </tr>
        {/if}

        {#each collection?.schema as field (field.name)}
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
                        普通文本值。
                    {:else if field.type === "number"}
                        数值。
                    {:else if field.type === "json"}
                        JSON 数组或对象。
                    {:else if field.type === "email"}
                        电子邮件地址。
                    {:else if field.type === "url"}
                        URL 地址。
                    {:else if field.type === "file"}
                        文件对象。<br />
                        设置为 <code>null</code> 以删除已上传的文件。
                    {:else if field.type === "relation"}
                        关系记录 {field.options?.maxSelect === 1 ? "id" : "ids"}。
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
                返回创建的记录时自动展开关系。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多 6 层嵌套关系的展开。 <br />
                展开的关系将附加到记录的 <code>expand</code> 属性下（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。
                <br />
                仅展开请求用户有权限<strong>查看</strong>的关系。
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