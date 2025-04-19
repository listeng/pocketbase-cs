<script>
    import Field from "@/components/base/Field.svelte";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import ImportPopup from "@/components/settings/ImportPopup.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addErrorToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { tick } from "svelte";

    $pageTitle = "导入集合";

    let fileInput;
    let importPopup;

    let schemas = "";
    let isLoadingFile = false;
    let newCollections = [];
    let oldCollections = [];
    let deleteMissing = true;
    let collectionsToUpdate = [];
    let isLoadingOldCollections = false;
    let mergeWithOldCollections = false; // 替代默认deleteMissing选项的另一种方式

    $: if (typeof schemas !== "undefined" && mergeWithOldCollections !== null) {
        loadNewCollections(schemas);
    }

    $: isValid =
        !!schemas &&
        newCollections.length &&
        newCollections.length === newCollections.filter((item) => !!item.id && !!item.name).length;

    $: collectionsToDelete = oldCollections.filter((collection) => {
        return (
            isValid &&
            !mergeWithOldCollections &&
            deleteMissing &&
            !CommonHelper.findByKey(newCollections, "id", collection.id)
        );
    });

    $: collectionsToAdd = newCollections.filter((collection) => {
        return isValid && !CommonHelper.findByKey(oldCollections, "id", collection.id);
    });

    $: if (typeof newCollections !== "undefined" || typeof deleteMissing !== "undefined") {
        loadCollectionsToUpdate();
    }

    $: hasChanges =
        !!schemas && (collectionsToDelete.length || collectionsToAdd.length || collectionsToUpdate.length);

    $: canImport = !isLoadingOldCollections && isValid && hasChanges;

    $: idReplacableCollections = newCollections.filter((collection) => {
        let old =
            CommonHelper.findByKey(oldCollections, "name", collection.name) ||
            CommonHelper.findByKey(oldCollections, "id", collection.id);

        if (!old) {
            return false; // 新集合
        }

        if (old.id != collection.id) {
            return true;
        }

        // 检查匹配的字段
        const oldFields = Array.isArray(old.fields) ? old.fields : [];
        const newFields = Array.isArray(collection.fields) ? collection.fields : [];
        for (const field of newFields) {
            const oldFieldById = CommonHelper.findByKey(oldFields, "id", field.id);
            if (oldFieldById) {
                continue; // 无需替换
            }

            const oldFieldByName = CommonHelper.findByKey(oldFields, "name", field.name);
            if (oldFieldByName && field.id != oldFieldByName.id) {
                return true;
            }
        }

        return false;
    });

    loadOldCollections();

    async function loadOldCollections() {
        isLoadingOldCollections = true;

        try {
            oldCollections = await ApiClient.collections.getFullList(200);
            for (let collection of oldCollections) {
                // 删除时间戳
                delete collection.created;
                delete collection.updated;

                // 清空OAuth2提供者
                delete collection.oauth2?.providers;
            }
        } catch (err) {
            ApiClient.error(err);
        }

        isLoadingOldCollections = false;
    }

    function loadCollectionsToUpdate() {
        collectionsToUpdate = [];

        if (!isValid) {
            return;
        }

        for (let newCollection of newCollections) {
            const oldCollection = CommonHelper.findByKey(oldCollections, "id", newCollection.id);
            if (
                // 没有旧集合
                !oldCollection?.id ||
                // 没有变化
                !CommonHelper.hasCollectionChanges(oldCollection, newCollection, deleteMissing)
            ) {
                continue;
            }

            collectionsToUpdate.push({
                new: newCollection,
                old: oldCollection,
            });
        }
    }

    function loadNewCollections() {
        newCollections = [];

        try {
            newCollections = JSON.parse(schemas);
        } catch (_) {}

        if (!Array.isArray(newCollections)) {
            newCollections = [];
        } else {
            newCollections = CommonHelper.filterDuplicatesByKey(newCollections);
        }

        // 标准化
        for (let collection of newCollections) {
            // 删除时间戳
            delete collection.created;
            delete collection.updated;

            // 合并重复字段
            collection.fields = CommonHelper.filterDuplicatesByKey(collection.fields);
        }
    }

    function replaceIds() {
        for (let collection of newCollections) {
            const old =
                CommonHelper.findByKey(oldCollections, "name", collection.name) ||
                CommonHelper.findByKey(oldCollections, "id", collection.id);

            if (!old) {
                continue;
            }

            const originalId = collection.id;
            const replacedId = old.id;
            collection.id = replacedId;

            // 替换字段ID
            const oldFields = Array.isArray(old.fields) ? old.fields : [];
            const newFields = Array.isArray(collection.fields) ? collection.fields : [];
            for (const field of newFields) {
                const oldField = CommonHelper.findByKey(oldFields, "name", field.name);
                if (oldField && oldField.id) {
                    field.id = oldField.id;
                }
            }

            // 更新引用
            for (let ref of newCollections) {
                if (!Array.isArray(ref.fields)) {
                    continue;
                }
                for (let field of ref.fields) {
                    if (field.collectionId && field.collectionId === originalId) {
                        field.collectionId = replacedId;
                    }
                }
            }

            // 更新包含集合ID的索引名
            for (let i = 0; i < collection.indexes?.length; i++) {
                collection.indexes[i] = collection.indexes[i].replace(
                    /create\s+(?:unique\s+)?\s*index\s*(?:if\s+not\s+exists\s+)?(\S*)\s+on/gim,
                    (v) => v.replace(originalId, replacedId),
                );
            }
        }

        schemas = JSON.stringify(newCollections, null, 4);
    }

    function loadFile(file) {
        isLoadingFile = true;

        const reader = new FileReader();

        reader.onload = async (event) => {
            isLoadingFile = false;
            fileInput.value = ""; // 重置

            schemas = event.target.result;

            await tick();

            if (!newCollections.length) {
                addErrorToast("无效的集合配置。");
                clear();
            }
        };

        reader.onerror = (err) => {
            console.warn(err);
            addErrorToast("加载JSON文件失败。");

            isLoadingFile = false;
            fileInput.value = ""; // 重置
        };

        reader.readAsText(file);
    }

    function clear() {
        schemas = "";
        fileInput.value = "";
        setErrors({});
    }

    function review() {
        const collectionsToImport = !mergeWithOldCollections
            ? newCollections
            : CommonHelper.filterDuplicatesByKey(oldCollections.concat(newCollections));

        importPopup?.show(oldCollections, collectionsToImport, deleteMissing);
    }
</script>

<SettingsSidebar />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">设置</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>

    <div class="wrapper">
        <div class="panel">
            {#if isLoadingOldCollections}
                <div class="loader" />
            {:else}
                <input
                    bind:this={fileInput}
                    type="file"
                    class="hidden"
                    accept=".json"
                    on:change={() => {
                        if (fileInput.files.length) {
                            loadFile(fileInput.files[0]);
                        }
                    }}
                />

                <div class="content txt-xl m-b-base">
                    <p>
                        在下方粘贴要导入的集合配置或
                        <button
                            class="btn btn-outline btn-sm m-l-5"
                            class:btn-loading={isLoadingFile}
                            on:click={() => {
                                fileInput.click();
                            }}
                        >
                            <span class="txt">从JSON文件加载</span>
                        </button>
                    </p>
                </div>

                <Field class="form-field {!isValid ? 'field-error' : ''}" name="collections" let:uniqueId>
                    <label for={uniqueId} class="p-b-10">集合</label>
                    <textarea
                        id={uniqueId}
                        class="code"
                        spellcheck="false"
                        rows="15"
                        required
                        bind:value={schemas}
                    />

                    {#if !!schemas && !isValid}
                        <div class="help-block help-block-error">无效的集合配置。</div>
                    {/if}
                </Field>

                {#if newCollections.length}
                    <Field class="form-field form-field-toggle" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            bind:checked={mergeWithOldCollections}
                            disabled={!isValid}
                        />
                        <label for={uniqueId}>与现有集合合并</label>
                    </Field>
                {/if}

                {#if false}
                    <!-- 暂时隐藏显式删除控制，根据用户反馈决定是否启用/移除 -->
                    <Field class="form-field form-field-toggle" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            bind:checked={deleteMissing}
                            disabled={!isValid}
                        />
                        <label for={uniqueId}>删除缺失的集合和字段</label>
                    </Field>
                {/if}

                {#if isValid && newCollections.length && !hasChanges}
                    <div class="alert alert-info">
                        <div class="icon">
                            <i class="ri-information-line" />
                        </div>
                        <div class="content">
                            <string>您的集合配置已是最新！</string>
                        </div>
                    </div>
                {/if}

                {#if isValid && newCollections.length && hasChanges}
                    <h5 class="section-title">检测到的变更</h5>

                    <div class="list">
                        {#if collectionsToDelete.length}
                            {#each collectionsToDelete as collection (collection.id)}
                                <div class="list-item">
                                    <span class="label label-danger list-label">已删除</span>
                                    <div class="inline-flex flex-gap-5">
                                        <strong>{collection.name}</strong>
                                        {#if collection.id}
                                            <small class="txt-hint">{collection.id}</small>
                                        {/if}
                                    </div>
                                </div>
                            {/each}
                        {/if}

                        {#if collectionsToUpdate.length}
                            {#each collectionsToUpdate as pair (pair.old.id + pair.new.id)}
                                <div class="list-item">
                                    <span class="label label-warning list-label">已更改</span>
                                    <div class="inline-flex flex-gap-5">
                                        {#if pair.old.name !== pair.new.name}
                                            <strong class="txt-strikethrough txt-hint">
                                                {pair.old.name}
                                            </strong>
                                            <i class="ri-arrow-right-line txt-sm" />
                                        {/if}
                                        <strong>{pair.new.name}</strong>
                                        {#if pair.new.id}
                                            <small class="txt-hint">{pair.new.id}</small>
                                        {/if}
                                    </div>
                                </div>
                            {/each}
                        {/if}

                        {#if collectionsToAdd.length}
                            {#each collectionsToAdd as collection (collection.id)}
                                <div class="list-item">
                                    <span class="label label-success list-label">已添加</span>
                                    <div class="inline-flex flex-gap-5">
                                        <strong>{collection.name}</strong>
                                        {#if collection.id}
                                            <small class="txt-hint">{collection.id}</small>
                                        {/if}
                                    </div>
                                </div>
                            {/each}
                        {/if}
                    </div>
                {/if}

                {#if idReplacableCollections.length}
                    <div class="alert alert-warning m-t-base">
                        <div class="icon">
                            <i class="ri-error-warning-line" />
                        </div>
                        <div class="content">
                            <string>
                                部分导入的集合名称和/或字段相同但使用了不同的ID。您可以选择替换它们。
                            </string>
                        </div>
                        <button
                            type="button"
                            class="btn btn-warning btn-sm btn-outline"
                            on:click={() => replaceIds()}
                        >
                            <span class="txt">替换为原始ID</span>
                        </button>
                    </div>
                {/if}

                <div class="flex m-t-base">
                    {#if !!schemas}
                        <button type="button" class="btn btn-transparent link-hint" on:click={() => clear()}>
                            <span class="txt">清空</span>
                        </button>
                    {/if}
                    <div class="flex-fill" />
                    <button
                        type="button"
                        class="btn btn-expanded btn-warning m-l-auto"
                        disabled={!canImport}
                        on:click={review}
                    >
                        <span class="txt">审核</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</PageWrapper>

<ImportPopup
    bind:this={importPopup}
    on:submit={() => {
        clear();
        loadOldCollections();
    }}
/>

<style>
    .list-label {
        min-width: 65px;
    }
</style>