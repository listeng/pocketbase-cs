<script>
    import { createEventDispatcher, tick } from "svelte";
    import { slide } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import { ClientResponseError } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { setErrors } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import ModelDateIcon from "@/components/base/ModelDateIcon.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import AuthFields from "@/components/records/fields/AuthFields.svelte";
    import TextField from "@/components/records/fields/TextField.svelte";
    import NumberField from "@/components/records/fields/NumberField.svelte";
    import BoolField from "@/components/records/fields/BoolField.svelte";
    import EmailField from "@/components/records/fields/EmailField.svelte";
    import UrlField from "@/components/records/fields/UrlField.svelte";
    import DateField from "@/components/records/fields/DateField.svelte";
    import SelectField from "@/components/records/fields/SelectField.svelte";
    import JsonField from "@/components/records/fields/JsonField.svelte";
    import FileField from "@/components/records/fields/FileField.svelte";
    import RelationField from "@/components/records/fields/RelationField.svelte";
    import EditorField from "@/components/records/fields/EditorField.svelte";
    import ExternalAuthsList from "@/components/records/ExternalAuthsList.svelte";

    const dispatch = createEventDispatcher();
    const formId = "record_" + CommonHelper.randomString(5);
    const tabFormKey = "form";
    const tabProviderKey = "providers";

    export let collection;

    let recordPanel;
    let original = {};
    let record = {};
    let initialDraft = null;
    let isSaving = false;
    let confirmHide = false; // prevent close recursion
    let uploadedFilesMap = {}; // eg.: {"field1":[File1, File2], ...}
    let deletedFileNamesMap = {}; // eg.: {"field1":[0, 1], ...}
    let originalSerializedData = JSON.stringify(original);
    let serializedData = originalSerializedData;
    let activeTab = tabFormKey;
    let isNew = true;
    let isLoading = true;
    let initialCollection = collection;

    $: isAuthCollection = collection?.type === "auth";

    $: hasEditorField = !!collection?.schema?.find((f) => f.type === "editor");

    $: hasFileChanges =
        CommonHelper.hasNonEmptyProps(uploadedFilesMap) || CommonHelper.hasNonEmptyProps(deletedFileNamesMap);

    $: serializedData = JSON.stringify(record);

    $: hasChanges = hasFileChanges || originalSerializedData != serializedData;

    $: isNew = !original || !original.id;

    $: canSave = !isLoading && (isNew || hasChanges);

    $: if (!isLoading) {
        updateDraft(serializedData);
    }

    $: if (collection && initialCollection?.id != collection?.id) {
        onCollectionChange();
    }

    export function show(model) {
        load(model);

        confirmHide = true;

        activeTab = tabFormKey;

        return recordPanel?.show();
    }

    export function hide() {
        return recordPanel?.hide();
    }

    function forceHide() {
        confirmHide = false;
        hide();
    }

    function onCollectionChange() {
        initialCollection = collection;

        if (!recordPanel?.isActive()) {
            return;
        }

        updateDraft(JSON.stringify(record));

        forceHide();
    }

    async function resolveModel(model) {
        if (model && typeof model === "string") {
            // load from id
            try {
                return await ApiClient.collection(collection.id).getOne(model);
            } catch (err) {
                if (!err.isAbort) {
                    forceHide();
                    console.warn("resolveModel:", err);
                    addErrorToast(`无法加载 ID 为 "${model}" 的模型的记录`);
                }
            }

            return null;
        }

        return model;
    }

    async function load(model) {
        isLoading = true;

        // resets
        setErrors({});
        uploadedFilesMap = {};
        deletedFileNamesMap = {};

        // load the minimum model data if possible to minimize layout shifts
        original =
            typeof model === "string"
                ? { id: model, collectionId: collection?.id, collectionName: collection?.name }
                : model || {};
        record = structuredClone(original);

        // resolve the complete model
        original = (await resolveModel(model)) || {};
        record = structuredClone(original);

        // wait to populate the fields to get the normalized values
        await tick();

        initialDraft = getDraft();
        if (!initialDraft || areRecordsEqual(record, initialDraft)) {
            initialDraft = null;
        } else {
            delete initialDraft.password;
            delete initialDraft.passwordConfirm;
        }

        originalSerializedData = JSON.stringify(record);

        isLoading = false;
    }

    async function replaceOriginal(newOriginal) {
        setErrors({}); // reset errors
        original = newOriginal || {};
        uploadedFilesMap = {};
        deletedFileNamesMap = {};

        // to avoid layout shifts we replace only the file and non-schema fields
        const skipFields = collection?.schema?.filter((f) => f.type != "file")?.map((f) => f.name) || [];
        for (let k in newOriginal) {
            if (skipFields.includes(k)) {
                continue;
            }
            record[k] = newOriginal[k];
        }

        // wait to populate the fields to get the normalized values
        await tick();

        originalSerializedData = JSON.stringify(record);

        deleteDraft();
    }

    function draftKey() {
        return "record_draft_" + (collection?.id || "") + "_" + (original?.id || "");
    }

    function getDraft(fallbackRecord) {
        try {
            const raw = window.localStorage.getItem(draftKey());
            if (raw) {
                return JSON.parse(raw);
            }
        } catch (_) {}

        return fallbackRecord;
    }

    function updateDraft(newSerializedData) {
        try {
            window.localStorage.setItem(draftKey(), newSerializedData);
        } catch (e) {
            // ignore local storage errors in case the serialized data
            // exceed the browser localStorage single value quota
            console.warn("updateDraft failure:", e);
            window.localStorage.removeItem(draftKey());
        }
    }

    function restoreDraft() {
        if (initialDraft) {
            record = initialDraft;
            initialDraft = null;
        }
    }

    function areRecordsEqual(recordA, recordB) {
        const cloneA = structuredClone(recordA || {});
        const cloneB = structuredClone(recordB || {});

        const fileFields = collection?.schema?.filter((f) => f.type === "file");
        for (let field of fileFields) {
            delete cloneA[field.name];
            delete cloneB[field.name];
        }

        // props to exclude from the checks
        const excludeProps = ["expand", "password", "passwordConfirm"];
        for (let prop of excludeProps) {
            delete cloneA[prop];
            delete cloneB[prop];
        }

        return JSON.stringify(cloneA) == JSON.stringify(cloneB);
    }

    function deleteDraft() {
        initialDraft = null;
        window.localStorage.removeItem(draftKey());
    }

    async function save(hidePanel = true) {
        if (isSaving || !canSave || !collection?.id) {
            return;
        }

        isSaving = true;

        try {
            const data = exportFormData();

            let result;
            if (isNew) {
                result = await ApiClient.collection(collection.id).create(data);
            } else {
                result = await ApiClient.collection(collection.id).update(record.id, data);
            }

            addSuccessToast(isNew ? "创建成功" : "更新成功");

            deleteDraft();

            if (hidePanel) {
                forceHide();
            } else {
                replaceOriginal(result);
            }

            dispatch("save", {
                isNew: isNew,
                record: result,
            });
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function deleteConfirm() {
        if (!original?.id) {
            return; // nothing to delete
        }

        confirm(`确定要删除选中的记录吗?`, () => {
            return ApiClient.collection(original.collectionId)
                .delete(original.id)
                .then(() => {
                    hide();
                    addSuccessToast("删除成功");
                    dispatch("delete", original);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function exportFormData() {
        const data = structuredClone(record || {});
        const formData = new FormData();

        const exportableFields = {
            id: data.id,
        };

        const jsonFields = {};

        for (const field of collection?.schema || []) {
            exportableFields[field.name] = true;

            if (field.type == "json") {
                jsonFields[field.name] = true;
            }
        }

        if (isAuthCollection) {
            exportableFields["username"] = true;
            exportableFields["email"] = true;
            exportableFields["emailVisibility"] = true;
            exportableFields["password"] = true;
            exportableFields["passwordConfirm"] = true;
            exportableFields["verified"] = true;
        }

        // export base fields
        for (const key in data) {
            // skip non-schema fields
            if (!exportableFields[key]) {
                continue;
            }

            // normalize nullable values
            if (typeof data[key] === "undefined") {
                data[key] = null;
            }

            // "validate" json fields
            if (jsonFields[key] && data[key] !== "") {
                try {
                    JSON.parse(data[key]);
                } catch (err) {
                    const fieldErr = {};
                    fieldErr[key] = {
                        code: "invalid_json",
                        message: err.toString(),
                    };
                    // emulate server error
                    throw new ClientResponseError({
                        status: 400,
                        response: {
                            data: fieldErr,
                        },
                    });
                }
            }

            CommonHelper.addValueToFormData(formData, key, data[key]);
        }

        // add uploaded files  (if any)
        for (const key in uploadedFilesMap) {
            const files = CommonHelper.toArray(uploadedFilesMap[key]);
            for (const file of files) {
                formData.append(key, file);
            }
        }

        // unset deleted files (if any)
        for (const key in deletedFileNamesMap) {
            const names = CommonHelper.toArray(deletedFileNamesMap[key]);
            for (const name of names) {
                formData.append(key + "." + name, "");
            }
        }

        return formData;
    }

    function sendVerificationEmail() {
        if (!collection?.id || !original?.email) {
            return;
        }

        confirm(`确定要发送验证邮件给 ${original.email}?`, () => {
            return ApiClient.collection(collection.id)
                .requestVerification(original.email)
                .then(() => {
                    addSuccessToast(`成功发送验证邮件给 ${original.email}.`);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function sendPasswordResetEmail() {
        if (!collection?.id || !original?.email) {
            return;
        }

        confirm(`确定要发送密码重置邮件给 ${original.email} 吗?`, () => {
            return ApiClient.collection(collection.id)
                .requestPasswordReset(original.email)
                .then(() => {
                    addSuccessToast(`成功发送密码重置邮件给 ${original.email}.`);
                })
                .catch((err) => {
                    ApiClient.error(err);
                });
        });
    }

    function duplicateConfirm() {
        if (hasChanges) {
            confirm("您还没有保存。确定要放弃吗?", () => {
                duplicate();
            });
        } else {
            duplicate();
        }
    }

    async function duplicate() {
        let clone = original ? structuredClone(original) : null;

        if (clone) {
            clone.id = "";
            clone.created = "";
            clone.updated = "";

            // reset file fields
            const fields = collection?.schema || [];
            for (const field of fields) {
                if (field.type === "file") {
                    delete clone[field.name];
                }
            }
        }

        deleteDraft();
        show(clone);

        await tick();

        originalSerializedData = "";
    }

    function handleFormKeydown(e) {
        if ((e.ctrlKey || e.metaKey) && e.code == "KeyS") {
            e.preventDefault();
            e.stopPropagation();
            save(false);
        }
    }
</script>

<OverlayPanel
    bind:this={recordPanel}
    class="
        record-panel
        {hasEditorField ? 'overlay-panel-xl' : 'overlay-panel-lg'}
        {isAuthCollection && !isNew ? 'colored-header' : ''}
    "
    btnClose={!isLoading}
    escClose={!isLoading}
    overlayClose={!isLoading}
    beforeHide={() => {
        if (hasChanges && confirmHide) {
            confirm("您还没有保存。确定要关闭吗?", () => {
                forceHide();
            });

            return false;
        }

        setErrors({});
        deleteDraft();

        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        {#if isLoading}
            <span class="loader loader-sm" />
            <h4 class="panel-title txt-hint">加载中...</h4>
        {:else}
            <h4 class="panel-title">
                {isNew ? "创建" : "编辑"}
                <strong>{collection?.name}</strong> 的记录
            </h4>

            {#if !isNew}
                <div class="flex-fill" />
                <div
                    tabindex="0"
                    role="button"
                    aria-label="More record options"
                    class="btn btn-sm btn-circle btn-transparent flex-gap-0"
                >
                    <i class="ri-more-line" aria-hidden="true" />
                    <Toggler class="dropdown dropdown-right dropdown-nowrap">
                        {#if isAuthCollection && !original.verified && original.email}
                            <button
                                type="button"
                                class="dropdown-item closable"
                                role="menuitem"
                                on:click={() => sendVerificationEmail()}
                            >
                                <i class="ri-mail-check-line" />
                                <span class="txt">发送验证邮件</span>
                            </button>
                        {/if}
                        {#if isAuthCollection && original.email}
                            <button
                                type="button"
                                class="dropdown-item closable"
                                role="menuitem"
                                on:click={() => sendPasswordResetEmail()}
                            >
                                <i class="ri-mail-lock-line" />
                                <span class="txt">发送重置密码邮件</span>
                            </button>
                        {/if}
                        <button
                            type="button"
                            class="dropdown-item closable"
                            role="menuitem"
                            on:click={() => duplicateConfirm()}
                        >
                            <i class="ri-file-copy-line" />
                            <span class="txt">复制</span>
                        </button>
                        <button
                            type="button"
                            class="dropdown-item txt-danger closable"
                            role="menuitem"
                            on:click|preventDefault|stopPropagation={() => deleteConfirm()}
                        >
                            <i class="ri-delete-bin-7-line" />
                            <span class="txt">删除</span>
                        </button>
                    </Toggler>
                </div>
            {/if}
        {/if}

        {#if isAuthCollection && !isNew}
            <div class="tabs-header stretched">
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === tabFormKey}
                    on:click={() => (activeTab = tabFormKey)}
                >
                    账号
                </button>
                <button
                    type="button"
                    class="tab-item"
                    class:active={activeTab === tabProviderKey}
                    on:click={() => (activeTab = tabProviderKey)}
                >
                    授权提供者
                </button>
            </div>
        {/if}
    </svelte:fragment>

    <div class="tabs-content no-animations">
        <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
        <form
            id={formId}
            class="tab-item"
            class:no-pointer-events={isLoading}
            class:active={activeTab === tabFormKey}
            on:submit|preventDefault={save}
            on:keydown={handleFormKeydown}
        >
            {#if !hasChanges && initialDraft && !isLoading}
                <div class="block" out:slide={{ duration: 150 }}>
                    <div class="alert alert-info m-0">
                        <div class="icon">
                            <i class="ri-information-line" />
                        </div>
                        <div class="flex flex-gap-xs">
                            这条记录之前有未保存的草稿.
                            <button
                                type="button"
                                class="btn btn-sm btn-secondary"
                                on:click={() => restoreDraft()}
                            >
                                恢复草稿
                            </button>
                        </div>
                        <button
                            type="button"
                            class="close"
                            aria-label="Discard draft"
                            use:tooltip={"Discard draft"}
                            on:click|preventDefault={() => deleteDraft()}
                        >
                            <i class="ri-close-line" />
                        </button>
                    </div>
                    <div class="clearfix p-b-base" />
                </div>
            {/if}

            <Field class="form-field {!isNew ? 'readonly' : ''}" name="id" let:uniqueId>
                <label for={uniqueId}>
                    <i class={CommonHelper.getFieldTypeIcon("primary")} />
                    <span class="txt">id</span>
                    <span class="flex-fill" />
                </label>
                {#if !isNew}
                    <div class="form-field-addon">
                        <ModelDateIcon model={record} />
                    </div>
                {/if}
                <input
                    type="text"
                    id={uniqueId}
                    placeholder={!isLoading ? "留空的话自动生成..." : ""}
                    minlength="15"
                    readonly={!isNew}
                    bind:value={record.id}
                />
            </Field>

            {#if isAuthCollection}
                <AuthFields bind:record {isNew} {collection} />

                {#if collection?.schema?.length}
                    <hr />
                {/if}
            {/if}

            {#each collection?.schema || [] as field (field.name)}
                {#if field.type === "text"}
                    <TextField {field} bind:value={record[field.name]} />
                {:else if field.type === "number"}
                    <NumberField {field} bind:value={record[field.name]} />
                {:else if field.type === "bool"}
                    <BoolField {field} bind:value={record[field.name]} />
                {:else if field.type === "email"}
                    <EmailField {field} bind:value={record[field.name]} />
                {:else if field.type === "url"}
                    <UrlField {field} bind:value={record[field.name]} />
                {:else if field.type === "editor"}
                    <EditorField {field} bind:value={record[field.name]} />
                {:else if field.type === "date"}
                    <DateField {field} bind:value={record[field.name]} />
                {:else if field.type === "select"}
                    <SelectField {field} bind:value={record[field.name]} />
                {:else if field.type === "json"}
                    <JsonField {field} bind:value={record[field.name]} />
                {:else if field.type === "file"}
                    <FileField
                        {field}
                        {record}
                        bind:value={record[field.name]}
                        bind:uploadedFiles={uploadedFilesMap[field.name]}
                        bind:deletedFileNames={deletedFileNamesMap[field.name]}
                    />
                {:else if field.type === "relation"}
                    <RelationField {field} bind:value={record[field.name]} />
                {/if}
            {/each}
        </form>

        {#if isAuthCollection && !isNew}
            <div class="tab-item" class:active={activeTab === tabProviderKey}>
                <ExternalAuthsList {record} />
            </div>
        {/if}
    </div>

    <svelte:fragment slot="footer">
        <button
            type="button"
            class="btn btn-transparent"
            disabled={isSaving || isLoading}
            on:click={() => hide()}
        >
            <span class="txt">取消</span>
        </button>

        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSaving || isLoading}
            disabled={!canSave || isSaving}
        >
            <span class="txt">{isNew ? "创建" : "保存"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<style>
    .panel-title {
        line-height: var(--smBtnHeight);
    }
</style>
