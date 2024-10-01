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

    $: isAuth = collection?.type === "auth";

    $: adminsOnly = collection?.updateRule === null;

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
                  "message": "更新记录失败。",
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
        {
            code: 404,
            body: `
                {
                  "code": 404,
                  "message": "请求的资源未找到。",
                  "data": {}
                }
            `,
        },
    ];

    $: if (collection.type === "auth") {
        baseData = {
            username: "test_username_update",
            emailVisibility: false,
            password: "87654321",
            passwordConfirm: "87654321",
            oldPassword: "12345678",
        };
    } else {
        baseData = {};
    }
</script>

<h3 class="m-b-sm">更新 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>更新单个 <strong>{collection.name}</strong> 记录。</p>
    <p>
        请求体参数可以以 <code>application/json</code> 或
        <code>multipart/form-data</code> 格式发送。
    </p>
    <p>
        仅支持通过 <code>multipart/form-data</code> 进行文件上传。
        <br />
        更多信息和示例可以查看详细的
        <a href={import.meta.env.PB_FILE_UPLOAD_DOCS} target="_blank" rel="noopener noreferrer">
            文件上传和处理文档
        </a>。
    </p>
    {#if isAuth}
        <p>
            <em>
                请注意，如果更改密码，当前记录的所有先前发出的令牌将自动失效，如果希望用户保持登录状态，您需要在更新调用后手动重新验证。
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
const data = ${JSON.stringify(Object.assign({}, baseData, CommonHelper.dummyCollectionSchemaData(collection)), null, 4)};

const record = await pb.collection('${collection?.name}').update('RECORD_ID', data);
    `}
    dart={`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${backendAbsUrl}');

...

// 示例更新请求体
final body = <String, dynamic>${JSON.stringify(Object.assign({}, baseData, CommonHelper.dummyCollectionSchemaData(collection)), null, 2)};

final record = await pb.collection('${collection?.name}').update('RECORD_ID', body: body);
    `}
/>

<h6 class="m-b-xs">API 详细信息</h6>
<div class="alert alert-warning">
    <strong class="label label-primary">PATCH</strong>
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
            <td>要更新的记录的 ID。</td>
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
                <td colspan="3" class="txt-hint">身份验证字段</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>用户名</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>身份验证记录的用户名。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>邮箱</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>
                    身份验证记录的电子邮件地址。
                    <br />
                    该字段只能由管理员或具有“管理”访问权限的身份验证记录更新。
                    <br />
                    普通账户可以通过调用“请求更改电子邮件”来更新其电子邮件。
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>邮箱可见性</span>
                    </div>
                </td>
                <td>
                    <span class="label">布尔值</span>
                </td>
                <td>在获取记录数据时是否显示/隐藏身份验证记录的电子邮件。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>旧密码</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>
                    旧的身份验证记录密码。
                    <br />
                    该字段仅在更改记录密码时需要。管理员和具有“管理”访问权限的身份验证记录可以跳过此字段。
                </td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>密码</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>新的身份验证记录密码。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>密码确认</span>
                    </div>
                </td>
                <td>
                    <span class="label">字符串</span>
                </td>
                <td>新的身份验证记录密码确认。</td>
            </tr>
            <tr>
                <td>
                    <div class="inline-flex">
                        <span class="label label-warning">可选</span>
                        <span>已验证</span>
                    </div>
                </td>
                <td>
                    <span class="label">布尔值</span>
                </td>
                <td>
                    指示身份验证记录是否已验证。
                    <br />
                    该字段只能由管理员或具有“管理”访问权限的身份验证记录设置。
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
                        关联记录 {field.options?.maxSelect > 1 ? "ids" : "id"}。
                    {:else if field.type === "user"}
                        用户 {field.options?.maxSelect > 1 ? "ids" : "id"}。
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
                在返回更新的记录时自动扩展关系。例如：
                <CodeBlock content={`?expand=relField1,relField2.subRelField21`} />
                支持最多 6 层嵌套关系扩展。<br />
                扩展的关系将附加到记录下的
                <code>expand</code> 属性（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。只有用户有权限<strong>查看</strong>的关系才会被扩展。
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