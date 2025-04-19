<script>
    import PasswordResetApiConfirmDocs from "@/components/collections/docs/PasswordResetApiConfirmDocs.svelte";
    import PasswordResetApiRequestDocs from "@/components/collections/docs/PasswordResetApiRequestDocs.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    const apiTabs = [
        { title: "请求密码重置", component: PasswordResetApiRequestDocs },
        { title: "确认密码重置", component: PasswordResetApiConfirmDocs },
    ];

    let activeApiTab = 0;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);
</script>

<h3 class="m-b-sm">密码重置 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>发送<strong>{collection.name}</strong>密码重置邮件请求。</p>
    <p>
        成功重置密码后，该记录所有先前颁发的认证令牌将自动失效。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').requestPasswordReset('test@example.com');

        // ---
        // (可选) 在自定义确认页面:
        // ---

        // 注意: 此调用后所有先前颁发的认证令牌将失效
        await pb.collection('${collection?.name}').confirmPasswordReset(
            'RESET_TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').requestPasswordReset('test@example.com');

        // ---
        // (可选) 在自定义确认页面:
        // ---

        // 注意: 此调用后所有先前颁发的认证令牌将失效
        await pb.collection('${collection?.name}').confirmPasswordReset(
          'RESET_TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `}
/>

<h6 class="m-b-xs">API详情</h6>
<div class="tabs">
    <div class="tabs-header compact">
        {#each apiTabs as tab, i}
            <button class="tab-item" class:active={activeApiTab == i} on:click={() => (activeApiTab = i)}>
                <div class="txt">{tab.title}</div>
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each apiTabs as tab, i}
            <div class="tab-item" class:active={activeApiTab == i}>
                <svelte:component this={tab.component} {collection} />
            </div>
        {/each}
    </div>
</div>