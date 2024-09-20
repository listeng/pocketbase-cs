<script>
    import { slide } from "svelte/transition";
    import tooltip from "@/actions/tooltip";
    import CommonHelper from "@/utils/CommonHelper";
    import RuleField from "@/components/collections/RuleField.svelte";

    export let collection;

    $: fields = CommonHelper.getAllCollectionIdentifiers(collection);

    let showFiltersInfo = false;
</script>

<div class="block m-b-sm handle">
    <div class="flex txt-sm txt-hint m-b-5">
        <p>
            所有规则都基于
            <a href={import.meta.env.PB_RULES_SYNTAX_DOCS} target="_blank" rel="noopener noreferrer">
                PocketBase 过滤器语法和操作符
            </a>.
        </p>
        <button
            type="button"
            class="expand-handle txt-sm txt-bold txt-nowrap link-hint"
            on:click={() => (showFiltersInfo = !showFiltersInfo)}
        >
            {showFiltersInfo ? "隐藏可用字段" : "显示可用字段"}
        </button>
    </div>

    {#if showFiltersInfo}
        <div transition:slide={{ duration: 150 }}>
            <div class="alert alert-warning m-0">
                <div class="content">
                    <p class="m-b-0">下面这些记录的字段可用:</p>
                    <div class="inline-flex flex-gap-5">
                        {#each fields as name}
                            <code>{name}</code>
                        {/each}
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        请求的字段可用特殊的过滤器 <em>@request</em> 来访问:
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@request.headers.*</code>
                        <code>@request.query.*</code>
                        <code>@request.data.*</code>
                        <code>@request.auth.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        您还可以使用 <em>@collection</em> 过滤器添加约束并查询其他集合：
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@collection.ANY_COLLECTION_NAME.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p>
                        示例规则:
                        <br />
                        <code>@request.auth.id != "" && created > "2022-01-01 00:00:00"</code>
                    </p>
                </div>
            </div>
        </div>
    {/if}
</div>

<RuleField label="列表/搜索规则" formKey="listRule" {collection} bind:rule={collection.listRule} />

<RuleField label="查看规则" formKey="viewRule" {collection} bind:rule={collection.viewRule} />

{#if collection?.type !== "view"}
    <RuleField label="创建规则" formKey="createRule" {collection} bind:rule={collection.createRule}>
        <svelte:fragment slot="afterLabel" let:isAdminOnly>
            {#if !isAdminOnly}
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `创建规则在提交的数据进行“预保存”后执行，您可以像在其他规则中一样访问主要记录字段。`,
                        position: "top",
                    }}
                />
            {/if}
        </svelte:fragment>
    </RuleField>

    <RuleField label="更新规则" formKey="updateRule" {collection} bind:rule={collection.updateRule} />

    <RuleField label="删除规则" formKey="deleteRule" {collection} bind:rule={collection.deleteRule} />
{/if}

{#if collection?.type === "auth"}
    <RuleField
        label="管理规则"
        formKey="options.manageRule"
        placeholder=""
        required={collection.options.manageRule !== null}
        {collection}
        bind:rule={collection.options.manageRule}
    >
        <svelte:fragment>
            <p>
                此 API 规则授予类似管理员的权限，允许完全管理身份认证记录，例如无需输入旧密码即可更改密码，直接更新已验证状态或电子邮件等。
            </p>
            <p>
                此规则与 <code>create</code> 和 <code>update</code> API 规则一起执行。
            </p>
        </svelte:fragment>
    </RuleField>
{/if}
