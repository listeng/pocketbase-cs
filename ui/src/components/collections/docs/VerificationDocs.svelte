<script>
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import VerificationApiConfirmDocs from "@/components/collections/docs/VerificationApiConfirmDocs.svelte";
    import VerificationApiRequestDocs from "@/components/collections/docs/VerificationApiRequestDocs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    const apiTabs = [
        { title: "请求验证", component: VerificationApiRequestDocs },
        { title: "确认验证", component: VerificationApiConfirmDocs },
    ];

    let activeApiTab = 0;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);
</script>

<h3 class="m-b-sm">账户验证 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>发送 <strong>{collection.name}</strong> 账户验证请求。</p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').requestVerification('test@example.com');

        // ---
        // (可选) 在自定义确认页面:
        // ---

        await pb.collection('${collection?.name}').confirmVerification('VERIFICATION_TOKEN');
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        await pb.collection('${collection?.name}').requestVerification('test@example.com');

        // ---
        // (可选) 在自定义确认页面:
        // ---

        await pb.collection('${collection?.name}').confirmVerification('VERIFICATION_TOKEN');
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