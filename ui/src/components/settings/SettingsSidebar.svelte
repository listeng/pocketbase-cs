<script>
    import { link } from "svelte-spa-router";
    import active from "svelte-spa-router/active";
    import PageSidebar from "@/components/base/PageSidebar.svelte";
    import { hideControls } from "@/stores/app";
</script>

<PageSidebar class="settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">系统</div>
        <a href="/settings" class="sidebar-list-item" use:active={{ path: "/settings" }} use:link>
            <i class="ri-home-gear-line" aria-hidden="true" />
            <span class="txt">应用</span>
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
                <span class="txt">导出数据集</span>
            </a>
            <a
                href="/settings/import-collections"
                class="sidebar-list-item"
                use:active={{ path: "/settings/import-collections/?.*" }}
                use:link
            >
                <i class="ri-install-line" aria-hidden="true" />
                <span class="txt">导入数据集</span>
            </a>
        {/if}

        <div class="sidebar-title">授权</div>
        <a
            href="/settings/auth-providers"
            class="sidebar-list-item"
            use:active={{ path: "/settings/auth-providers/?.*" }}
            use:link
        >
            <i class="ri-lock-password-line" aria-hidden="true" />
            <span class="txt">第三方登录</span>
        </a>
        <a
            href="/settings/tokens"
            class="sidebar-list-item"
            use:active={{ path: "/settings/tokens/?.*" }}
            use:link
        >
            <i class="ri-key-line" aria-hidden="true" />
            <span class="txt">令牌设置</span>
        </a>
        <a
            href="/settings/admins"
            class="sidebar-list-item"
            use:active={{ path: "/settings/admins/?.*" }}
            use:link
        >
            <i class="ri-shield-user-line" aria-hidden="true" />
            <span class="txt">管理员</span>
        </a>
    </div>
</PageSidebar>
