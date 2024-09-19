<script>
    import PocketBase, { getTokenPayload } from "pocketbase";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";

    export let params;

    let password = "";
    let isLoading = false;
    let success = false;

    $: newEmail = CommonHelper.getJWTPayload(params?.token).newEmail || "";

    async function submit() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        // init a custom client to avoid interfering with the admin state
        const client = new PocketBase(import.meta.env.PB_BACKEND_URL);

        try {
            const payload = getTokenPayload(params?.token);
            await client.collection(payload.collectionId).confirmEmailChange(params?.token, password);
            success = true;
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }
</script>

<FullPage nobranding>
    {#if success}
        <div class="alert alert-success">
            <div class="icon"><i class="ri-checkbox-circle-line" /></div>
            <div class="content txt-bold">
                <p>修改邮箱地址成功。</p>
                <p>现在你可以用新邮箱登录了。</p>
            </div>
        </div>

        <button type="button" class="btn btn-transparent btn-block" on:click={() => window.close()}>
            关闭
        </button>
    {:else}
        <form on:submit|preventDefault={submit}>
            <div class="content txt-center m-b-base">
                <h5>
                    输入密码以确认修改邮箱
                    {#if newEmail}
                        to <strong class="txt-nowrap">{newEmail}</strong>
                    {/if}
                </h5>
            </div>

            <Field class="form-field required" name="password" let:uniqueId>
                <label for={uniqueId}>密码</label>
                <!-- svelte-ignore a11y-autofocus -->
                <input type="password" id={uniqueId} required autofocus bind:value={password} />
            </Field>

            <button
                type="submit"
                class="btn btn-lg btn-block"
                class:btn-loading={isLoading}
                disabled={isLoading}
            >
                <span class="txt">确认新邮箱</span>
            </button>
        </form>
    {/if}
</FullPage>
