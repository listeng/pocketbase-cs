<script>
    import { createEventDispatcher, tick } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollection;
    let newCollection;
    let hideAfterSave;
    let conflictingOIDCs = [];
    let changedRules = [];

    $: isCollectionRenamed = oldCollection?.name != newCollection?.name;

    $: isNewCollectionView = newCollection?.type === "view";

    $: isNewCollectionAuth = newCollection?.type === "auth";

    $: renamedFields =
        (!isNewCollectionView &&
            newCollection?.fields?.filter(
                (field) => field.id && !field._toDelete && field._originalName != field.name,
            )) ||
        [];

    $: deletedFields =
        (!isNewCollectionView && newCollection?.fields?.filter((field) => field.id && field._toDelete)) || [];

    $: multipleToSingleFields =
        newCollection?.fields?.filter((field) => {
            const old = oldCollection?.fields?.find((f) => f.id == field.id);
            if (!old) {
                return false;
            }
            return old.maxSelect != 1 && field.maxSelect == 1;
        }) || [];

    $: showChanges = !isNewCollectionView || isCollectionRenamed || changedRules.length;

    export async function show(original, changed, hideAfterSaveArg = true) {
        oldCollection = original;
        newCollection = changed;
        hideAfterSave = hideAfterSaveArg;

        await detectConflictingOIDCs();

        detectRulesChange();

        await tick();

        if (
            isCollectionRenamed ||
            renamedFields.length ||
            deletedFields.length ||
            multipleToSingleFields.length ||
            conflictingOIDCs.length ||
            changedRules.length
        ) {
            panel?.show();
        } else {
            // 没有需要审查的变更 -> 直接确认
            confirm();
        }
    }

    export function hide() {
        panel?.hide();
    }

    function confirm() {
        hide();
        dispatch("confirm", hideAfterSave);
    }

    const oidcProviders = ["oidc", "oidc2", "oidc3"];

    async function detectConflictingOIDCs() {
        conflictingOIDCs = [];

        for (let name of oidcProviders) {
            let oldProvider = oldCollection?.oauth2?.providers?.find((p) => p.name == name);
            let newProvider = newCollection?.oauth2?.providers?.find((p) => p.name == name);

            if (!oldProvider || !newProvider) {
                continue;
            }

            let oldHost = new URL(oldProvider.authURL).host;
            let newHost = new URL(newProvider.authURL).host;
            if (oldHost == newHost) {
                continue;
            }

            // 检查是否存在现有的外部认证
            if (await haveExternalAuths(name)) {
                conflictingOIDCs.push({ name, oldHost, newHost });
            }
        }
    }

    async function haveExternalAuths(provider) {
        try {
            await ApiClient.collection("_externalAuths").getFirstListItem(
                ApiClient.filter("collectionRef={:collectionId} && provider={:provider}", {
                    collectionId: newCollection?.id,
                    provider: provider,
                }),
            );
            return true;
        } catch {}

        return false;
    }

    function getExternalAuthsFilterLink(provider) {
        return `#/collections?collection=_externalAuths&filter=collectionRef%3D%22${newCollection?.id}%22+%26%26+provider%3D%22${provider}%22`;
    }

    function detectRulesChange() {
        changedRules = [];

        // 目前仅对"生产"环境启用
        if (window.location.protocol != "https:") {
            return;
        }

        const ruleProps = ["listRule", "viewRule"];
        if (!isNewCollectionView) {
            ruleProps.push("createRule", "updateRule", "deleteRule");
        }
        if (isNewCollectionAuth) {
            ruleProps.push("manageRule", "authRule");
        }

        let oldRule, newRule;
        for (let prop of ruleProps) {
            oldRule = oldCollection?.[prop];
            newRule = newCollection?.[prop];
            if (oldRule === newRule) {
                continue;
            }

            changedRules.push({ prop, oldRule, newRule });
        }
    }
</script>

<OverlayPanel bind:this={panel} class="confirm-changes-panel" popup on:hide on:show>
    <svelte:fragment slot="header">
        <h4>确认集合变更</h4>
    </svelte:fragment>

    {#if isCollectionRenamed || deletedFields.length || renamedFields.length}
        <div class="alert alert-warning">
            <div class="icon">
                <i class="ri-error-warning-line" />
            </div>
            <div class="content txt-bold">
                <p>
                    如果任何集合变更是其他集合规则、过滤器或视图查询的一部分，
                    您需要手动更新它们！
                </p>
                {#if deletedFields.length}
                    <p>所有与被删除字段关联的数据将被永久删除！</p>
                {/if}
            </div>
        </div>
    {/if}

    {#if showChanges}
        <h6>变更内容：</h6>
        <ul class="changes-list">
            {#if isCollectionRenamed}
                <li>
                    <div class="inline-flex">
                        重命名集合
                        <strong class="txt-strikethrough txt-hint">{oldCollection?.name}</strong>
                        <i class="ri-arrow-right-line txt-sm" />
                        <strong class="txt"> {newCollection?.name}</strong>
                    </div>
                </li>
            {/if}

            {#if !isNewCollectionView}
                {#each multipleToSingleFields as field}
                    <li>
                        将字段 <strong>{field.name}</strong> 从多值转换为单值
                        <em class="txt-sm">(将只保留最后一个数组项)</em>
                    </li>
                {/each}

                {#each renamedFields as field}
                    <li>
                        <div class="inline-flex">
                            重命名字段
                            <strong class="txt-strikethrough txt-hint">{field._originalName}</strong>
                            <i class="ri-arrow-right-line txt-sm" />
                            <strong class="txt">{field.name}</strong>
                        </div>
                    </li>
                {/each}

                {#each deletedFields as field}
                    <li class="txt-danger">
                        删除字段 <span class="txt-bold">{field.name}</span>
                    </li>
                {/each}
            {/if}

            {#each changedRules as ruleChange}
                <li>
                    更改API规则 <code class="txt-sm">{ruleChange.prop}</code>:
                    <br />
                    <small class="txt-mono txt-hint">
                        <strong>旧值</strong>:
                        <span class="txt-preline">
                            {ruleChange.oldRule === null
                                ? "null (仅超级用户)"
                                : ruleChange.oldRule || '""'}
                        </span>
                    </small>
                    <br />
                    <small class="txt-mono txt-success">
                        <strong>新值</strong>:
                        <span class="txt-preline">
                            {ruleChange.newRule === null
                                ? "null (仅超级用户)"
                                : ruleChange.newRule || '""'}
                        </span>
                    </small>
                </li>
            {/each}

            {#each conflictingOIDCs as oidc}
                <li>
                    更改 <code>{oidc.name}</code> 主机
                    <div class="inline-flex m-l-5">
                        <strong class="txt-strikethrough txt-hint">{oidc.oldHost}</strong>
                        <i class="ri-arrow-right-line txt-sm" />
                        <strong class="txt">{oidc.newHost}</strong>
                    </div>
                    <br />
                    <em class="txt-hint">
                        如果新旧OIDC配置不是针对同一提供者，请考虑删除
                        当前集合和提供者关联的所有旧 <code class="txt-sm">_externalAuths</code> 记录，
                        否则可能导致账户链接错误。
                        <a href={getExternalAuthsFilterLink(oidc.name)} target="_blank">
                            查看现有的 <code class="txt-sm">_externalAuths</code> 记录
                            <i class="ri-external-link-line txt-sm"></i>
                        </a>.
                    </em>
                </li>
            {/each}
        </ul>
    {/if}

    <svelte:fragment slot="footer">
        <!-- svelte-ignore a11y-autofocus -->
        <button autofocus type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">取消</span>
        </button>
        <button type="button" class="btn btn-expanded" on:click={() => confirm()}>
            <span class="txt">确认</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style lang="scss">
    .changes-list {
        word-break: break-word;
        line-height: var(--smLineHeight);
        li {
            margin-top: 10px;
            margin-bottom: 10px;
        }
    }
</style>