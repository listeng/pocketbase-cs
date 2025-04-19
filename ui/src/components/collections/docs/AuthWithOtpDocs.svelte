<script>
    import AuthWithOtpApiAuthDocs from "@/components/collections/docs/AuthWithOtpApiAuthDocs.svelte";
    import AuthWithOtpApiRequestDocs from "@/components/collections/docs/AuthWithOtpApiRequestDocs.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    const apiTabs = [
        { title: "OTP请求", component: AuthWithOtpApiRequestDocs },
        { title: "OTP认证", component: AuthWithOtpApiAuthDocs },
    ];

    let activeApiTab = 0;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);
</script>

<h3 class="m-b-sm">使用OTP认证 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>使用一次性密码(OTP)进行认证。</p>
    <p>
        请注意，在请求OTP时，即使提供的邮箱不存在用户记录，我们也会返回一个<code>otpId</code>作为基本的枚举保护措施。
    </p>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        // 向指定认证记录发送OTP邮件
        const req = await pb.collection('${collection?.name}').requestOTP('test@example.com');

        // ... 显示一个界面/弹窗让用户输入邮件中的密码 ...

        // 使用请求的OTP ID和邮件密码进行认证
        const authData = await pb.collection('${collection?.name}').authWithOTP(
            req.otpId,
            "你的OTP",
        );

        // 认证后可以从authStore访问认证数据
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "登出"
        pb.authStore.clear();
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        // 向指定认证记录发送OTP邮件
        final req = await pb.collection('${collection?.name}').requestOTP('test@example.com');

        // ... 显示一个界面/弹窗让用户输入邮件中的密码 ...

        // 使用请求的OTP ID和邮件密码进行认证
        final authData = await pb.collection('${collection?.name}').authWithOTP(
            req.otpId,
            "你的OTP",
        );

        // 认证后可以从authStore访问认证数据
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "登出"
        pb.authStore.clear();
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