<script>
    import { tick } from "svelte";
    import { replace } from "svelte-spa-router";
    import { getTokenPayload } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import { addInfoToast, addErrorToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";
    import Field from "@/components/base/Field.svelte";
    import FullPage from "@/components/base/FullPage.svelte";

    export let params;

    let email = "";
    let password = "";
    let passwordConfirm = "";
    let isLoading = false;
    let isUploading = false;

    let emailInput;
    let backupFileInput;

    $: isBusy = isLoading || isUploading;

    checkToken();

    async function checkToken() {
        if (!params?.token) {
            return replace("/");
        }

        isLoading = true;

        try {
            const payload = getTokenPayload(params?.token);

            await ApiClient.collection("_superusers").getOne(payload.id, {
                requestKey: "installer_token_check",
                headers: { Authorization: params?.token },
            });
        } catch (err) {
            if (!err?.isAbort) {
                addErrorToast("安装令牌无效或已过期。");

                replace("/");
            }
        }

        isLoading = false;

        await tick();

        emailInput?.focus();
    }

    async function submit() {
        if (isBusy) {
            return;
        }

        isLoading = true;

        try {
            await ApiClient.collection("_superusers").create(
                {
                    email,
                    password,
                    passwordConfirm,
                },
                {
                    headers: { Authorization: params?.token },
                },
            );

            await ApiClient.collection("_superusers").authWithPassword(email, password);

            replace("/");
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    function resetSelectedBackupFile() {
        if (backupFileInput) {
            backupFileInput.value = "";
        }
    }

    function uploadBackupConfirm(file) {
        if (!file) {
            return;
        }

        confirm(
            `请注意，我们不会对上传的备份文件进行验证。请谨慎操作，仅在上传可信文件源时继续。\n\n` +
                `确定要上传并初始化"${file.name}"吗？`,
            () => {
                uploadBackup(file);
            },
            () => {
                resetSelectedBackupFile();
            },
        );
    }

    async function uploadBackup(file) {
        if (!file || isBusy) {
            return;
        }

        isUploading = true;

        try {
            await ApiClient.backups.upload(
                { file: file },
                {
                    headers: { Authorization: params?.token },
                },
            );

            await ApiClient.backups.restore(file.name, {
                headers: { Authorization: params?.token },
            });

            addInfoToast("正在解压上传的存档文件，请稍候！");

            // optimistic restore completion
            await new Promise((r) => setTimeout(r, 2000));

            replace("/");
        } catch (err) {
            ApiClient.error(err);
        }

        resetSelectedBackupFile();

        isUploading = false;
    }
</script>

<FullPage>
    <form class="block" autocomplete="off" on:submit|preventDefault={submit}>
        <div class="content txt-center m-b-base">
            <h4>创建您的第一个超级用户账户以继续</h4>
        </div>

        <Field class="form-field required" name="email" let:uniqueId>
            <label for={uniqueId}>邮箱</label>
            <input
                bind:this={emailInput}
                type="email"
                autocomplete="off"
                id={uniqueId}
                disabled={isBusy}
                bind:value={email}
                required
            />
        </Field>

        <Field class="form-field required" name="password" let:uniqueId>
            <label for={uniqueId}>密码</label>
            <input
                type="password"
                autocomplete="new-password"
                minlength="10"
                id={uniqueId}
                disabled={isBusy}
                bind:value={password}
                required
            />
            <div class="help-block">建议至少10个字符。</div>
        </Field>

        <Field class="form-field required" name="passwordConfirm" let:uniqueId>
            <label for={uniqueId}>确认密码</label>
            <input
                type="password"
                minlength="10"
                id={uniqueId}
                disabled={isBusy}
                bind:value={passwordConfirm}
                required
            />
        </Field>

        <button
            type="submit"
            class="btn btn-lg btn-block btn-next"
            class:btn-disabled={isBusy}
            class:btn-loading={isLoading}
        >
            <span class="txt">创建超级用户并登录</span>
            <i class="ri-arrow-right-line" />
        </button>
    </form>

    <hr />

    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <label
        for="backupFileInput"
        class="btn btn-lg btn-hint btn-transparent btn-block"
        class:btn-disabled={isBusy}
        class:btn-loading={isUploading}
    >
        <i class="ri-upload-cloud-line" />
        <span class="txt">或从备份初始化</span>
    </label>
    <input
        bind:this={backupFileInput}
        id="backupFileInput"
        type="file"
        class="hidden"
        accept=".zip"
        on:change={(e) => {
            uploadBackupConfirm(e.target?.files?.[0]);
        }}
    />
</FullPage>