<script>
    import { onMount } from "svelte";
    import { slide } from "svelte/transition";
    import ApiClient from "@/utils/ApiClient";
    import tooltip from "@/actions/tooltip";
    import { removeError } from "@/stores/errors";
    import Field from "@/components/base/Field.svelte";
    import RedactedPasswordInput from "@/components/base/RedactedPasswordInput.svelte";

    const testRequestKey = "s3_test_request";

    export let originalConfig = {};
    export let config = {};
    export let configKey = "s3";
    export let toggleLabel = "启用S3";
    export let testFilesystem = "storage"; // storage or backups
    export let testError = null;
    export let isTesting = false;

    let testTimeoutId = null;
    let testDebounceId = null;
    let maskSecret = false;

    $: if (originalConfig?.enabled) {
        refreshMaskSecret();
        testConnectionWithDebounce(100);
    }

    // clear s3 errors on disable
    $: if (!config.enabled) {
        removeError(configKey);
    }

    function refreshMaskSecret() {
        maskSecret = !!originalConfig?.accessKey;
    }

    function testConnectionWithDebounce(timeout) {
        isTesting = true;
        clearTimeout(testDebounceId);
        testDebounceId = setTimeout(() => {
            testConnection();
        }, timeout);
    }

    async function testConnection() {
        testError = null;

        if (!config.enabled) {
            isTesting = false;
            return testError; // nothing to test
        }

        // auto cancel the test request after 30sec
        ApiClient.cancelRequest(testRequestKey);
        clearTimeout(testTimeoutId);
        testTimeoutId = setTimeout(() => {
            ApiClient.cancelRequest(testRequestKey);
            testError = new Error("S3连接测试超时。");
            isTesting = false;
        }, 30000);

        isTesting = true;

        let err;

        try {
            await ApiClient.settings.testS3(testFilesystem, {
                $cancelKey: testRequestKey,
            });
        } catch (e) {
            err = e;
        }

        if (!err?.isAbort) {
            testError = err;
            isTesting = false;
            clearTimeout(testTimeoutId);
        }

        return testError;
    }

    onMount(() => {
        return () => {
            clearTimeout(testTimeoutId);
            clearTimeout(testDebounceId);
        };
    });
</script>

<Field class="form-field form-field-toggle" let:uniqueId>
    <input type="checkbox" id={uniqueId} required bind:checked={config.enabled} />
    <label for={uniqueId}>{toggleLabel}</label>
</Field>

<slot {isTesting} {testError} enabled={config.enabled} />

{#if config.enabled}
    <div class="grid" transition:slide={{ duration: 150 }}>
        <div class="col-lg-6">
            <Field class="form-field required" name="{configKey}.endpoint" let:uniqueId>
                <label for={uniqueId}>端点</label>
                <input type="text" id={uniqueId} required bind:value={config.endpoint} />
            </Field>
        </div>
        <div class="col-lg-3">
            <Field class="form-field required" name="{configKey}.bucket" let:uniqueId>
                <label for={uniqueId}>存储桶</label>
                <input type="text" id={uniqueId} required bind:value={config.bucket} />
            </Field>
        </div>
        <div class="col-lg-3">
            <Field class="form-field required" name="{configKey}.region" let:uniqueId>
                <label for={uniqueId}>区域</label>
                <input type="text" id={uniqueId} required bind:value={config.region} />
            </Field>
        </div>
        <div class="col-lg-6">
            <Field class="form-field required" name="{configKey}.accessKey" let:uniqueId>
                <label for={uniqueId}>访问密钥</label>
                <input type="text" id={uniqueId} required bind:value={config.accessKey} />
            </Field>
        </div>
        <div class="col-lg-6">
            <Field class="form-field required" name="{configKey}.secret" let:uniqueId>
                <label for={uniqueId}>密钥</label>
                <RedactedPasswordInput
                    required
                    id={uniqueId}
                    bind:mask={maskSecret}
                    bind:value={config.secret}
                />
            </Field>
        </div>
        <div class="col-lg-12">
            <Field class="form-field" name="{configKey}.forcePathStyle" let:uniqueId>
                <input type="checkbox" id={uniqueId} bind:checked={config.forcePathStyle} />
                <label for={uniqueId}>
                    <span class="txt">强制路径样式寻址</span>
                    <i
                        class="ri-information-line link-hint"
                        use:tooltip={{
                            text: '强制请求使用路径样式寻址，例如"https://s3.amazonaws.com/BUCKET/KEY"而不是默认的"https://BUCKET.s3.amazonaws.com/KEY"。',
                            position: "top",
                        }}
                    />
                </label>
            </Field>
        </div>
        <!-- margin helper -->
        <div class="col-lg-12" />
    </div>
{/if}