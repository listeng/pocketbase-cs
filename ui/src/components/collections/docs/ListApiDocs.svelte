<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import FilterSyntax from "@/components/collections/docs/FilterSyntax.svelte";
    import SdkTabs from "@/components/collections/docs/SdkTabs.svelte";
    import FieldsQueryParam from "@/components/collections/docs/FieldsQueryParam.svelte";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: fieldNames = CommonHelper.getAllCollectionIdentifiers(collection);

    $: adminsOnly = collection?.listRule === null;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseUrl);

    $: if (collection?.id) {
        responses.push({
            code: 200,
            body: JSON.stringify(
                {
                    page: 1,
                    perPage: 30,
                    totalPages: 1,
                    totalItems: 2,
                    items: [
                        CommonHelper.dummyCollectionRecord(collection),
                        CommonHelper.dummyCollectionRecord(collection),
                    ],
                },
                null,
                2
            ),
        });

        responses.push({
            code: 400,
            body: `
                {
                  "code": 400,
                  "message": "处理您的请求时发生错误。无效的过滤器。",
                  "data": {}
                }
            `,
        });

        if (adminsOnly) {
            responses.push({
                code: 403,
                body: `
                    {
                      "code": 403,
                      "message": "只有管理员能操作",
                      "data": {}
                    }
                `,
            });
        }
    }
</script>

<h3 class="m-b-sm">列表/搜索 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>
        获取分页的 <strong>{collection.name}</strong> 记录列表，支持排序和过滤。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        // 获取分页记录列表
        const resultList = await pb.collection('${collection?.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // 你也可以通过 getFullList 一次性获取所有记录
        const records = await pb.collection('${collection?.name}').getFullList({
            sort: '-created',
        });

        // 或者仅获取第一个符合指定过滤条件的记录
        const record = await pb.collection('${collection?.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        // 获取分页记录列表
        final resultList = await pb.collection('${collection?.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // 你也可以通过 getFullList 一次性获取所有记录
        final records = await pb.collection('${collection?.name}').getFullList(
          sort: '-created',
        );

        // 或者仅获取第一个符合指定过滤条件的记录
        final record = await pb.collection('${collection?.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}
/>

<h6 class="m-b-xs">API 详细信息</h6>
<div class="alert alert-info">
    <strong class="label label-primary">GET</strong>
    <div class="content">
        <p>
            /pb-proxy/api/collections/<strong>{collection.name}</strong>/records
        </p>
    </div>
    {#if adminsOnly}
        <p class="txt-hint txt-sm txt-right">需要管理员 <code>Authorization:TOKEN</code> 头</p>
    {/if}
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
            <td>page</td>
            <td>
                <span class="label">数字</span>
            </td>
            <td>分页列表的页码（默认值为 1）。</td>
        </tr>
        <tr>
            <td>perPage</td>
            <td>
                <span class="label">数字</span>
            </td>
            <td>指定每页返回的最大记录数（默认值为 30）。</td>
        </tr>
        <tr>
            <td>sort</td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                指定记录的排序属性。 <br />
                在属性前添加 <code>-</code> / <code>+</code>（默认）以指定降序/升序。
                例如：
                <CodeBlock
                    content={`
                        // 按创建时间降序，按 ID 升序
                        ?sort=-created,id
                    `}
                />
                <p>
                    <strong>支持的记录排序字段：</strong> <br />
                    <code>@random</code>,
                    {#each fieldNames as name, i}
                        <code>{name}</code>{i < fieldNames.length - 1 ? ", " : ""}
                    {/each}
                </p>
            </td>
        </tr>
        <tr>
            <td>filter</td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                过滤返回的记录。示例：
                <CodeBlock
                    content={`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}
                />
                <FilterSyntax />
            </td>
        </tr>
        <tr>
            <td>expand</td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>
                自动扩展记录关系。示例：
                <CodeBlock content={`?expand=relField1,relField2.subRelField`} />
                支持最多 6 层深的嵌套关系扩展。 <br />
                扩展的关系将附加到每个单独的记录下的
                <code>expand</code> 属性中（例如 <code>{`"expand": {"relField1": {...}, ...}`}</code>）。
                <br />
                仅扩展请求用户有权限<strong>查看</strong>的关系。
            </td>
        </tr>
        <FieldsQueryParam />
        <tr>
            <td id="query-page">skipTotal</td>
            <td>
                <span class="label">布尔值</span>
            </td>
            <td>
                如果设置，则总计数查询将被跳过，响应字段
                <code>totalItems</code> 和 <code>totalPages</code> 将有 <code>-1</code> 的值。
                <br />
                当不需要总计数或使用基于游标的分页时，这可以大幅加快搜索查询速度。
                <br />
                出于优化目的，默认情况下在
                <code>getFirstListItem()</code>
                和
                <code>getFullList()</code> SDK 方法中设置。
            </td>
        </tr>
    </tbody>
</table>

<div class="section-title">响应</div>
<div class="tabs">
    <div class="tabs-header compact combined left">
        {#each responses as response (response.code)}
            <button
                type="button"
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