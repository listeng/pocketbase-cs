<script>
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import { addSuccessToast } from "@/stores/toasts";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import TokenField from "@/components/settings/TokenField.svelte";

    const recordTokensList = [
        { key: "recordAuthToken", label: "记录认证令牌" },
        { key: "recordVerificationToken", label: "记录邮箱验证令牌" },
        { key: "recordPasswordResetToken", label: "记录密码重置令牌" },
        { key: "recordEmailChangeToken", label: "记录邮箱更改令牌" },
        { key: "recordFileToken", label: "记录受保护文件访问令牌" },
    ];

    const adminTokensList = [
        { key: "adminAuthToken", label: "管理员认证令牌" },
        { key: "adminPasswordResetToken", label: "管理员密码重置令牌" },
        { key: "adminFileToken", label: "管理员受保护文件访问令牌" },
    ];

    $pageTitle = "令牌设置";

    let originalFormSettings = {};
    let formSettings = {};
    let isLoading = false;
    let isSaving = false;

    $: initialHash = JSON.stringify(originalFormSettings);

    $: hasChanges = initialHash != JSON.stringify(formSettings);

    loadSettings();

    async function loadSettings() {
        isLoading = true;

        try {
            const result = (await ApiClient.settings.getAll()) || {};
            initSettings(result);
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    async function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        try {
            const result = await ApiClient.settings.update(CommonHelper.filterRedactedProps(formSettings));
            initSettings(result);
            addSuccessToast("Successfully saved tokens options.");
        } catch (err) {
            ApiClient.error(err);
        }

        isSaving = false;
    }

    function initSettings(data) {
        data = data || {};
        formSettings = {};

        const tokensList = recordTokensList.concat(adminTokensList);

        for (const listItem of tokensList) {
            formSettings[listItem.key] = {
                duration: data[listItem.key]?.duration || 0,
            };
        }

        originalFormSettings = JSON.parse(JSON.stringify(formSettings));
    }

    function reset() {
        formSettings = JSON.parse(JSON.stringify(originalFormSettings || {}));
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
        <form class="panel" autocomplete="off" on:submit|preventDefault={save}>
            <div class="content m-b-sm txt-xl">
                <p>设置令牌参数</p>
            </div>

            {#if isLoading}
                <div class="loader" />
            {:else}
                <h3 class="section-title">记录令牌</h3>
                {#each recordTokensList as token (token.key)}
                    <TokenField
                        key={token.key}
                        label={token.label}
                        bind:duration={formSettings[token.key].duration}
                        bind:secret={formSettings[token.key].secret}
                    />
                {/each}

                <hr />

                <h3 class="section-title">管理令牌</h3>
                {#each adminTokensList as token (token.key)}
                    <TokenField
                        key={token.key}
                        label={token.label}
                        bind:duration={formSettings[token.key].duration}
                        bind:secret={formSettings[token.key].secret}
                    />
                {/each}

                <div class="flex">
                    <div class="flex-fill" />
                    {#if hasChanges}
                        <button
                            type="button"
                            class="btn btn-transparent btn-hint"
                            disabled={isSaving}
                            on:click={() => reset()}
                        >
                            <span class="txt">取消</span>
                        </button>
                    {/if}
                    <button
                        type="submit"
                        class="btn btn-expanded"
                        class:btn-loading={isSaving}
                        disabled={!hasChanges || isSaving}
                        on:click={() => save()}
                    >
                        <span class="txt">保存</span>
                    </button>
                </div>
            {/if}
        </form>
    </div>
</PageWrapper>
