<script>
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import { setErrors } from "@/stores/errors";
    import { addErrorToast, addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher, tick } from "svelte";

    const dispatch = createEventDispatcher();

    const formId = "email_test_" + CommonHelper.randomString(5);
    const emailStorageKey = "last_email_test";
    const testRequestKey = "email_test_request";

    const templateOptions = [
        { label: "验证邮件", value: "verification" },
        { label: "密码重置", value: "password-reset" },
        { label: "邮箱变更确认", value: "email-change" },
        { label: "一次性密码", value: "otp" },
        { label: "登录提醒", value: "login-alert" },
    ];

    let panel;
    let collectionIdOrName = "";
    let email = localStorage.getItem(emailStorageKey);
    let template = templateOptions[0].value;
    let isSubmitting = false;
    let testTimeoutId = null;
    let authCollections = [];
    let isAuthCollectionsLoading = false;
    let showAuthCollections = false;

    $: canSubmit = !!email && !!template && !!collectionIdOrName;

    export function show(collectionArg = "", emailArg = "", templateArg = "") {
        setErrors({}); // reset any previous errors

        showAuthCollections = false;

        collectionIdOrName = collectionArg || "";
        if (!collectionIdOrName) {
            loadAuthCollections();
        }

        email = emailArg || localStorage.getItem(emailStorageKey);

        template = templateArg || templateOptions[0].value;

        panel?.show();
    }

    export function hide() {
        clearTimeout(testTimeoutId);
        return panel?.hide();
    }

    async function submit() {
        if (!canSubmit || isSubmitting || !collectionIdOrName) {
            return;
        }

        isSubmitting = true;

        // store in local storage for later use
        localStorage?.setItem(emailStorageKey, email);

        // auto cancel the test request after 30sec
        clearTimeout(testTimeoutId);
        testTimeoutId = setTimeout(() => {
            ApiClient.cancelRequest(testRequestKey);
            addErrorToast("测试邮件发送超时");
        }, 30000);

        try {
            await ApiClient.settings.testEmail(collectionIdOrName, email, template, {
                $cancelKey: testRequestKey,
            });

            addSuccessToast("测试邮件发送成功");
            dispatch("submit");
            isSubmitting = false;

            await tick();

            hide();
        } catch (err) {
            isSubmitting = false;
            ApiClient.error(err);
        }

        clearTimeout(testTimeoutId);
    }

    async function loadAuthCollections() {
        showAuthCollections = true;
        isAuthCollectionsLoading = true;

        try {
            authCollections = await ApiClient.collections.getFullList({
                filter: "type='auth'",
                sort: "+name",
                requestKey: formId + "_collections_loading",
            });

            collectionIdOrName = authCollections[0]?.id || "";
            isAuthCollectionsLoading = false;
        } catch (err) {
            if (!err.isAbort) {
                isAuthCollectionsLoading = false;
                ApiClient.error(err);
            }
        }
    }
</script>

<OverlayPanel
    bind:this={panel}
    class="overlay-panel-sm email-test-popup"
    overlayClose={!isSubmitting}
    escClose={!isSubmitting}
    beforeHide={() => !isSubmitting}
    popup
    on:show
    on:hide
>
    <svelte:fragment slot="header">
        <h4 class="center txt-break">发送测试邮件</h4>
    </svelte:fragment>

    <form id={formId} autocomplete="off" on:submit|preventDefault={() => submit()}>
        <Field class="form-field required" name="template" let:uniqueId>
            {#each templateOptions as option (option.value)}
                <div class="form-field-block">
                    <input
                        type="radio"
                        name="template"
                        id={uniqueId + option.value}
                        value={option.value}
                        bind:group={template}
                    />
                    <label for={uniqueId + option.value}>{option.label}</label>
                </div>
            {/each}
        </Field>

        {#if showAuthCollections}
            <Field class="form-field required" name="collection" let:uniqueId>
                <label for={uniqueId}>认证集合</label>
                <ObjectSelect
                    id={uniqueId}
                    selectPlaceholder={isAuthCollectionsLoading
                        ? "正在加载认证集合..."
                        : "选择认证集合"}
                    noOptionsText={"未找到认证集合"}
                    selectionKey="id"
                    items={authCollections}
                    bind:keyOfSelected={collectionIdOrName}
                />
            </Field>
        {/if}

        <Field class="form-field required m-0" name="email" let:uniqueId>
            <label for={uniqueId}>收件邮箱</label>
            <!-- svelte-ignore a11y-autofocus -->
            <input type="email" id={uniqueId} autofocus required bind:value={email} />
        </Field>
    </form>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-transparent" on:click={hide} disabled={isSubmitting}
            >关闭</button
        >
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSubmitting}
            disabled={!canSubmit || isSubmitting}
        >
            <i class="ri-mail-send-line" />
            <span class="txt">发送</span>
        </button>
    </svelte:fragment>
</OverlayPanel>