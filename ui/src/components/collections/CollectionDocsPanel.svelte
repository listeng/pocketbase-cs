<script>
    import tooltip from "@/actions/tooltip";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";

    const baseTabs = {
        list: {
            label: "列表/搜索",
            component: import("@/components/collections/docs/ListApiDocs.svelte"),
        },
        view: {
            label: "查看",
            component: import("@/components/collections/docs/ViewApiDocs.svelte"),
        },
        create: {
            label: "创建",
            component: import("@/components/collections/docs/CreateApiDocs.svelte"),
        },
        update: {
            label: "更新",
            component: import("@/components/collections/docs/UpdateApiDocs.svelte"),
        },
        delete: {
            label: "删除",
            component: import("@/components/collections/docs/DeleteApiDocs.svelte"),
        },
        realtime: {
            label: "实时",
            component: import("@/components/collections/docs/RealtimeApiDocs.svelte"),
        },
        batch: {
            label: "批量",
            component: import("@/components/collections/docs/BatchApiDocs.svelte"),
        },
    };

    const authTabs = {
        "list-auth-methods": {
            label: "列出认证方法",
            component: import("@/components/collections/docs/AuthMethodsDocs.svelte"),
        },
        "auth-with-password": {
            label: "密码认证",
            component: import("@/components/collections/docs/AuthWithPasswordDocs.svelte"),
        },
        "auth-with-oauth2": {
            label: "OAuth2认证",
            component: import("@/components/collections/docs/AuthWithOAuth2Docs.svelte"),
        },
        "auth-with-otp": {
            label: "OTP认证",
            component: import("@/components/collections/docs/AuthWithOtpDocs.svelte"),
        },
        refresh: {
            label: "认证刷新",
            component: import("@/components/collections/docs/AuthRefreshDocs.svelte"),
        },
        verification: {
            label: "验证",
            component: import("@/components/collections/docs/VerificationDocs.svelte"),
        },
        "password-reset": {
            label: "密码重置",
            component: import("@/components/collections/docs/PasswordResetDocs.svelte"),
        },
        "email-change": {
            label: "邮箱变更",
            component: import("@/components/collections/docs/EmailChangeDocs.svelte"),
        },
    };

    let docsPanel;
    let collection = {};
    let activeTab;
    let tabs = [];

    $: if (collection.type === "auth") {
        tabs = Object.assign({}, baseTabs, authTabs);
        tabs["auth-with-password"].disabled = !collection.passwordAuth.enabled;
        tabs["auth-with-oauth2"].disabled = !collection.oauth2.enabled;
        tabs["auth-with-otp"].disabled = !collection.otp.enabled;
    } else if (collection.type === "view") {
        tabs = Object.assign({}, baseTabs);
        delete tabs.create;
        delete tabs.update;
        delete tabs.delete;
        delete tabs.realtime;
        delete tabs.batch;
    } else {
        tabs = Object.assign({}, baseTabs);
    }

    // reset active tab on tabs list change
    if (tabs.length) {
        activeTab = Object.keys(tabs)[0];
    }

    export function show(model) {
        collection = model;

        changeTab(Object.keys(tabs)[0]);

        return docsPanel?.show();
    }

    export function hide() {
        return docsPanel?.hide();
    }

    export function changeTab(newTab) {
        activeTab = newTab;
    }
</script>

<OverlayPanel bind:this={docsPanel} on:hide on:show class="docs-panel">
    <div class="docs-content-wrapper">
        <aside class="docs-sidebar" class:compact={collection?.type === "auth"}>
            <nav class="sidebar-content">
                {#each Object.entries(tabs) as [key, tab], i (key)}
                    <!-- add a separator before the first auth tab -->
                    {#if i === Object.keys(baseTabs).length}
                        <hr class="m-t-sm m-b-sm" />
                    {/if}

                    {#if tab.disabled}
                        <div
                            class="sidebar-item disabled"
                            use:tooltip={{ position: "left", text: "该集合未启用此功能" }}
                        >
                            {tab.label}
                        </div>
                    {:else}
                        <button
                            type="button"
                            class="sidebar-item"
                            class:active={activeTab === key}
                            on:click={() => changeTab(key)}
                        >
                            {tab.label}
                        </button>
                    {/if}
                {/each}
            </nav>
        </aside>

        <div class="docs-content">
            {#each Object.entries(tabs) as [key, tab] (key)}
                {#if activeTab === key}
                    {#await tab.component then { default: TabComponent }}
                        <TabComponent {collection} />
                    {/await}
                {/if}
            {/each}
        </div>
    </div>

    <!-- visible only on small screens -->
    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={() => hide()}>
            <span class="txt">关闭</span>
        </button>
    </svelte:fragment>
</OverlayPanel>