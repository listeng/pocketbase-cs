<script>
    import { createEventDispatcher, onDestroy } from "svelte";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";

    const dispatch = createEventDispatcher();
    const backupRequestKey = "upload_backup";

    let classes = "";
    export { classes as class };

    let fileInput;
    let isUploading = false;

    function resetSelectedFile() {
        if (fileInput) {
            fileInput.value = "";
        }
    }

    function uploadConfirm(file) {
        if (!file) {
            return;
        }

        confirm(
            `请注意，我们不会对上传的备份文件进行验证。请谨慎操作，仅在上传来源可信时继续。\n\n` +
                `确定要上传"${file.name}"吗？`,
            () => {
                upload(file);
            },
            () => {
                resetSelectedFile();
            },
        );
    }

    async function upload(file) {
        if (isUploading || !file) {
            return;
        }

        isUploading = true;

        const data = new FormData();
        data.set("file", file);

        try {
            await ApiClient.backups.upload(data, { requestKey: backupRequestKey });
            isUploading = false;
            dispatch("success");
            addSuccessToast("备份文件上传成功");
        } catch (err) {
            if (!err.isAbort) {
                isUploading = false;
                if (err.response?.data?.file?.message) {
                    addErrorToast(err.response.data.file.message);
                } else {
                    ApiClient.error(err);
                }
            }
        }

        resetSelectedFile();
    }

    onDestroy(() => {
        ApiClient.cancelRequest(backupRequestKey);
    });
</script>

<button
    type="button"
    class="btn btn-circle btn-transparent {classes}"
    class:btn-loading={isUploading}
    class:btn-disabled={isUploading}
    aria-label="上传备份"
    use:tooltip={"上传备份"}
    on:click={() => fileInput?.click()}
>
    <i class="ri-upload-cloud-line" />
</button>

<input
    bind:this={fileInput}
    type="file"
    accept="application/zip"
    class="hidden"
    on:change={(e) => {
        uploadConfirm(e?.target?.files?.[0]);
    }}
/>