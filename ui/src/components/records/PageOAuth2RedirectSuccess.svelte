<script>
    import { onMount } from "svelte";
    import { pageTitle } from "@/stores/app";

    $pageTitle = "OAuth2 认证完成";

    onMount(() => {
        window.close();
    });
</script>

<div class="content txt-hint txt-center p-base">
    <h3 class="m-b-sm">认证完成。</h3>
    <h5>您可以关闭此窗口并返回应用。</h5>
</div>