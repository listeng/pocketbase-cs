<script>
    import { tick, createEventDispatcher } from "svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const dispatch = createEventDispatcher();

    let panel;
    let oldCollection;
    let newCollection;

    $: isCollectionRenamed = oldCollection?.name != newCollection?.name;

    $: isNewCollectionView = newCollection?.type === "view";

    $: renamedFields =
        newCollection?.schema?.filter(
            (field) => field.id && !field.toDelete && field.originalName != field.name
        ) || [];

    $: deletedFields = newCollection?.schema?.filter((field) => field.id && field.toDelete) || [];

    $: multipleToSingleFields =
        newCollection?.schema?.filter((field) => {
            const old = oldCollection?.schema?.find((f) => f.id == field.id);
            if (!old) {
                return false;
            }
            return old.options?.maxSelect != 1 && field.options?.maxSelect == 1;
        }) || [];

    $: showChanges = !isNewCollectionView || isCollectionRenamed;

    export async function show(original, changed) {
        oldCollection = original;
        newCollection = changed;

        await tick();

        if (
            isCollectionRenamed ||
            renamedFields.length ||
            deletedFields.length ||
            multipleToSingleFields.length
        ) {
            panel?.show();
        } else {
            // no changes to review -> confirm directly
            confirm();
        }
    }

    export function hide() {
        panel?.hide();
    }

    function confirm() {
        hide();
        dispatch("confirm");
    }
</script>

<OverlayPanel bind:this={panel} class="confirm-changes-panel" popup on:hide on:show>
    <svelte:fragment slot="header">
        <h4>修改数据集确认</h4>
    </svelte:fragment>

    <div class="alert alert-warning">
        <div class="icon">
            <i class="ri-error-warning-line" />
        </div>
        <div class="content txt-bold">
            <p>
                如果任何数据集的更改属于另一个数据集规则、过滤器或视图查询，您将需要手动更新它！
            </p>
            {#if deletedFields.length}
                <p>与被删除字段相关的所有数据将被永久删除！</p>
            {/if}
        </div>
    </div>

    {#if showChanges}
        <h6>修改:</h6>
        <ul class="changes-list">
            {#if isCollectionRenamed}
                <li>
                    <div class="inline-flex">
                        重命名数据集
                        <strong class="txt-strikethrough txt-hint">{oldCollection?.name}</strong>
                        <i class="ri-arrow-right-line txt-sm" />
                        <strong class="txt"> {newCollection?.name}</strong>
                    </div>
                </li>
            {/if}

            {#if !isNewCollectionView}
                {#each multipleToSingleFields as field}
                    <li>
                        字段的多值转换为单值
                        <strong>{field.name}</strong>
                        <em class="txt-sm">(将只保留最后一个数组项)</em>
                    </li>
                {/each}

                {#each renamedFields as field}
                    <li>
                        <div class="inline-flex">
                            重命名字段
                            <strong class="txt-strikethrough txt-hint">{field.originalName}</strong>
                            <i class="ri-arrow-right-line txt-sm" />
                            <strong class="txt"> {field.name}</strong>
                        </div>
                    </li>
                {/each}

                {#each deletedFields as field}
                    <li class="txt-danger">
                        删除字段 <span class="txt-bold">{field.name}</span>
                    </li>
                {/each}
            {/if}
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
