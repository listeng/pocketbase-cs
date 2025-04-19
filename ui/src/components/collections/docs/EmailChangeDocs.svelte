<script>
    import EmailChangeApiConfirmDocs from "@/components/collections/docs/EmailChangeApiConfirmDocs.svelte";
    import EmailChangeApiRequestDocs from "@/components/collections/docs/EmailChangeApiRequestDocs.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    const apiTabs = [
        { title: "请求更改邮箱", component: EmailChangeApiRequestDocs },
        { title: "确认邮箱更改", component: EmailChangeApiConfirmDocs },
    ];

    let activeApiTab = 0;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);
</script>

<h3 class="m-b-sm">邮箱更改 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>发送 <strong>{collection.name}</strong> 邮箱更改请求。</p>
    <p>
        成功更改邮箱后，该记录之前颁发的所有认证令牌将自动失效。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${collection?.name}').requestEmailChange('new@example.com');

        // ---
        // (可选) 在您的自定义确认页面中:
        // ---

        // 注意: 此调用后，之前颁发的所有认证令牌将失效
        await pb.collection('${collection?.name}').confirmEmailChange(
            'EMAIL_CHANGE_TOKEN',
            'YOUR_PASSWORD',
        );
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${collection?.name}').requestEmailChange('new@example.com');

        ...

        // ---
        // (可选) 在您的自定义确认页面中:
        // ---

        // 注意: 此调用后，之前颁发的所有认证令牌将失效
        await pb.collection('${collection?.name}').confirmEmailChange(
          'EMAIL_CHANGE_TOKEN',
          'YOUR_PASSWORD',
        );
    `}
/>

<h6 class="m-b-xs">API 详情</h6>
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