<script>
    import tooltip from "@/actions/tooltip";
    import RuleField from "@/components/collections/RuleField.svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import { slide } from "svelte/transition";

    export let collection;

    $: fieldNames = CommonHelper.getAllCollectionIdentifiers(collection);

    $: hiddenFieldNames = collection.fields?.filter((f) => f.hidden).map((f) => f.name);

    let showFiltersInfo = false;

    let showExtraRules = collection.manageRule !== null || collection.authRule !== "";
</script>

<div class="block m-b-sm handle">
    <div class="flex txt-sm txt-hint m-b-5">
        <p>
            所有规则遵循
            <a href={import.meta.env.PB_RULES_SYNTAX_DOCS} target="_blank" rel="noopener noreferrer">
                PocketBase过滤语法和操作符
            </a>。
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
                    <p class="m-b-0">以下记录字段可用：</p>
                    <div class="inline-flex flex-gap-5">
                        {#each fieldNames as name}
                            {#if !hiddenFieldNames.includes(name)}
                                <code>{name}</code>
                            {/if}
                        {/each}
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        请求字段可以通过特殊的 <em>@request</em> 过滤器访问：
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@request.headers.*</code>
                        <code>@request.query.*</code>
                        <code>@request.body.*</code>
                        <code>@request.auth.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p class="m-b-0">
                        您还可以使用 <em>@collection</em> 过滤器添加约束和查询其他集合：
                    </p>
                    <div class="inline-flex flex-gap-5">
                        <code>@collection.ANY_COLLECTION_NAME.*</code>
                    </div>

                    <hr class="m-t-10 m-b-5" />

                    <p>
                        示例规则：
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
        <svelte:fragment slot="afterLabel" let:isSuperuserOnly>
            {#if !isSuperuserOnly}
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `主记录字段包含将要插入数据库的值。`,
                        position: "top",
                    }}
                />
            {/if}
        </svelte:fragment>
    </RuleField>

    <RuleField label="更新规则" formKey="updateRule" {collection} bind:rule={collection.updateRule}>
        <svelte:fragment slot="afterLabel" let:isSuperuserOnly>
            {#if !isSuperuserOnly}
                <i
                    class="ri-information-line link-hint"
                    use:tooltip={{
                        text: `主记录字段代表旧的/现有的记录字段值。\n要针对新提交的值，可以使用 @request.body.*`,
                        position: "top",
                    }}
                />
            {/if}
        </svelte:fragment>
    </RuleField>

    <RuleField label="删除规则" formKey="deleteRule" {collection} bind:rule={collection.deleteRule} />
{/if}

{#if collection?.type === "auth"}
    <hr />

    <button
        type="button"
        class="btn btn-sm m-b-sm {showExtraRules ? 'btn-secondary' : 'btn-hint btn-transparent'}"
        on:click={() => {
            showExtraRules = !showExtraRules;
        }}
    >
        <strong class="txt">额外的认证集合规则</strong>
        {#if showExtraRules}
            <i class="ri-arrow-up-s-line txt-sm" />
        {:else}
            <i class="ri-arrow-down-s-line txt-sm" />
        {/if}
    </button>

    {#if showExtraRules}
        <div class="block" transition:slide={{ duration: 150 }}>
            <RuleField
                label="认证规则"
                formKey="authRule"
                placeholder=""
                {collection}
                bind:rule={collection.authRule}
            >
                <svelte:fragment>
                    <p>
                        此规则在每次认证前执行，允许您限制谁可以进行认证。
                    </p>
                    <p>
                        例如，仅允许已验证用户，可以设置为 <code>verified = true</code>。
                    </p>
                    <p>留空则允许任何有账户的用户认证。</p>
                    <p>要完全禁用认证，可以改为"仅限超级用户"。</p>
                </svelte:fragment>
            </RuleField>

            <RuleField
                label="管理规则"
                formKey="manageRule"
                placeholder=""
                required={collection.manageRule !== null}
                {collection}
                bind:rule={collection.manageRule}
            >
                <svelte:fragment>
                    <p>
                        此规则在 <code>create</code> 和 <code>update</code> API 规则之外额外执行。
                    </p>
                    <p>
                        它提供类似超级用户的权限，允许完全管理认证记录，例如无需输入旧密码即可更改密码、直接更新验证状态或电子邮件等。
                    </p>
                </svelte:fragment>
            </RuleField>
        </div>
    {/if}
{/if}