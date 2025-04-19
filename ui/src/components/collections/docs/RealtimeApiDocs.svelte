<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import SdkTabs from "@/components/base/SdkTabs.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    $: backendAbsUrl = CommonHelper.getApiExampleUrl(ApiClient.baseURL);
</script>

<h3 class="m-b-sm">实时更新 ({collection.name})</h3>
<div class="content txt-lg m-b-sm">
    <p>通过服务器推送事件(SSE)订阅实时变更。</p>
    <p>
        当记录发生<strong>新增</strong>、<strong>更新</strong>
        或<strong>删除</strong>操作时会触发事件（参见下方的"事件数据格式"部分）。
    </p>
</div>
<div class="alert alert-info m-t-10 m-b-sm">
    <div class="icon">
        <i class="ri-information-line" />
    </div>
    <div class="contet">
        <p>
            <strong>您可以订阅单个记录或整个集合。</strong>
        </p>
        <p>
            当订阅<strong>单个记录</strong>时，系统会使用集合的
            <strong>查看规则(ViewRule)</strong>来判断订阅者是否有权限接收事件消息。
        </p>
        <p>
            当订阅<strong>整个集合</strong>时，系统会使用集合的
            <strong>列表规则(ListRule)</strong>来判断订阅者是否有权限接收事件消息。
        </p>
    </div>
</div>

<SdkTabs
    js={`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${backendAbsUrl}');

        ...

        // (可选)身份验证
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // 订阅${collection?.name}集合中所有记录的变更
        pb.collection('${collection?.name}').subscribe('*', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* 其他选项如：filter, expand, custom headers等 */ });

        // 仅订阅指定记录的变更
        pb.collection('${collection?.name}').subscribe('RECORD_ID', function (e) {
            console.log(e.action);
            console.log(e.record);
        }, { /* 其他选项如：filter, expand, custom headers等 */ });

        // 取消订阅
        pb.collection('${collection?.name}').unsubscribe('RECORD_ID'); // 移除所有'RECORD_ID'订阅
        pb.collection('${collection?.name}').unsubscribe('*'); // 移除所有'*'主题订阅
        pb.collection('${collection?.name}').unsubscribe(); // 移除该集合所有订阅
    `}
    dart={`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${backendAbsUrl}');

        ...

        // (可选)身份验证
        await pb.collection('users').authWithPassword('test@example.com', '123456');

        // 订阅${collection?.name}集合中所有记录的变更
        pb.collection('${collection?.name}').subscribe('*', (e) {
            print(e.action);
            print(e.record);
        }, /* 其他选项如：filter, expand, custom headers等 */);

        // 仅订阅指定记录的变更
        pb.collection('${collection?.name}').subscribe('RECORD_ID', (e) {
            print(e.action);
            print(e.record);
        }, /* 其他选项如：filter, expand, custom headers等 */);

        // 取消订阅
        pb.collection('${collection?.name}').unsubscribe('RECORD_ID'); // 移除所有'RECORD_ID'订阅
        pb.collection('${collection?.name}').unsubscribe('*'); // 移除所有'*'主题订阅
        pb.collection('${collection?.name}').unsubscribe(); // 移除该集合所有订阅
    `}
/>

<h6 class="m-b-xs">API详情</h6>
<div class="alert">
    <strong class="label label-primary">SSE</strong>
    <div class="content">
        <p>/api/realtime</p>
    </div>
</div>

<div class="section-title">事件数据格式</div>
<CodeBlock
    content={JSON.stringify(
        {
            action: "create",
            record: CommonHelper.dummyCollectionRecord(collection),
        },
        null,
        2,
    ).replace('"action": "create"', '"action": "create" // 新增、更新或删除')}
/>