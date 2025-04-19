<script>
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import { hideControls } from "@/stores/app";
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
</script>

<PageSidebar class="settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">系统</div>
        <a href="/settings" class="sidebar-list-item" use:active={{ path: "/settings" }} use:link>
            <i class="ri-home-gear-line" aria-hidden="true" />
            <span class="txt">应用设置</span>
        </a>
        <a
            href="/settings/mail"
            class="sidebar-list-item"
            use:active={{ path: "/settings/mail/?.*" }}
            use:link
        >
            <i class="ri-send-plane-2-line" aria-hidden="true" />
            <span class="txt">邮件设置</span>
        </a>
        <a
            href="/settings/storage"
            class="sidebar-list-item"
            use:active={{ path: "/settings/storage/?.*" }}
            use:link
        >
            <i class="ri-archive-drawer-line" aria-hidden="true" />
            <span class="txt">文件存储</span>
        </a>
        <a
            href="/settings/backups"
            class="sidebar-list-item"
            use:active={{ path: "/settings/backups/?.*" }}
            use:link
        >
            <i class="ri-archive-line" aria-hidden="true" />
            <span class="txt">备份</span>
        </a>
        <a
            href="/settings/crons"
            class="sidebar-list-item"
            use:active={{ path: "/settings/crons/?.*" }}
            use:link
        >
            <i class="ri-time-line" aria-hidden="true"></i>
            <span class="txt">定时任务</span>
        </a>

        {#if !$hideControls}
            <div class="sidebar-title">
                <span class="txt">同步</span>
            </div>
            <a
                href="/settings/export-collections"
                class="sidebar-list-item"
                use:active={{ path: "/settings/export-collections/?.*" }}
                use:link
            >
                <i class="ri-uninstall-line" aria-hidden="true" />
                <span class="txt">导出集合</span>
            </a>
            <a
                href="/settings/import-collections"
                class="sidebar-list-item"
                use:active={{ path: "/settings/import-collections/?.*" }}
                use:link
            >
                <i class="ri-install-line" aria-hidden="true" />
                <span class="txt">导入集合</span>
            </a>
        {/if}
    </div>
</PageSidebar>