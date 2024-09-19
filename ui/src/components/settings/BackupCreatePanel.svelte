<script>
    import { createEventDispatcher, onDestroy } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { setErrors } from "@/stores/errors";
    import { addInfoToast, addSuccessToast } from "@/stores/toasts";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Field from "@/components/base/Field.svelte";

    const dispatch = createEventDispatcher();

    const formId = "backup_create_" + CommonHelper.randomString(5);

    let panel;
    let name = "";
    let isSubmitting = false;
    let submitTimeoutId;

    export function show(newName) {
        setErrors({});
        isSubmitting = false;
        name = newName || "";
        panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    async function submit() {
        if (isSubmitting) {
            return;
        }

        isSubmitting = true;

        clearTimeout(submitTimeoutId);
        submitTimeoutId = setTimeout(() => {
            hide();
        }, 1500);

        try {
            await ApiClient.backups.create(name, { $cancelKey: formId });

            isSubmitting = false;

            hide();
            dispatch("submit");
            addSuccessToast("生成新的备份成功");
        } catch (err) {
            if (!err.isAbort) {
                ApiClient.error(err);
            }
        }

        clearTimeout(submitTimeoutId);
        isSubmitting = false;
    }

    onDestroy(() => {
        clearTimeout(submitTimeoutId);
    });
</script>

<OverlayPanel
    bind:this={panel}
    class="backup-create-panel"
    beforeOpen={() => {
        if (isSubmitting) {
            addInfoToast("正在备份，请稍等");
            return false;
        }

        return true;
    }}
    beforeHide={() => {
        if (isSubmitting) {
            addInfoToast(
                "备份一起动，需要耗费一些时间。您可以稍后再来查看。",
                4500
            );
        }

        return true;
    }}
    popup
    on:show
    on:hide
>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">初始化新的备份</h4>
    </svelte:fragment>

    <div class="alert alert-info">
        <div class="icon">
            <i class="ri-information-line" />
        </div>
        <div class="content">
            <p>
                请注意，在备份过程中，其他并发的写请求可能会失败，因为数据库会暂时“锁定”（这种情况通常只在生成 ZIP 文件时发生）。
            </p>
            <p class="txt-bold">
                如果您使用 S3 存储来上传文件（在数据集中存在上传文件字段），您必须单独备份这些文件，因为它们不存储在本地，且不会包含在最终备份中！
            </p>
        </div>
    </div>

    <form id={formId} autocomplete="off" on:submit|preventDefault={submit}>
        <Field class="form-field m-0" name="name" let:uniqueId>
            <label for={uniqueId}>备份名称</label>
            <input
                type="text"
                id={uniqueId}
                placeholder={"留空的话自动生成"}
                pattern="^[a-z0-9_-]+\.zip$"
                bind:value={name}
            />
            <em class="help-block">格式为 [a-z0-9_-].zip</em>
        </Field>
    </form>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide} disabled={isSubmitting}>
            <span class="txt">取消</span>
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSubmitting}
            disabled={isSubmitting}
        >
            <span class="txt">开始备份</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
