<script>
    import providersList from "@/providers.js";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let record;

    let externalAuths = [];
    let isLoading = false;

    function getProviderConfig(provider) {
        return providersList.find((p) => p.key == provider) || {};
    }

    function getProviderTitle(provider) {
        return getProviderConfig(provider)?.title || CommonHelper.sentenize(provider, false);
    }

    async function loadExternalAuths() {
        if (!record?.id) {
            externalAuths = [];
            isLoading = false;
            return;
        }

        isLoading = true;

        try {
            externalAuths = await ApiClient.collection("_externalAuths").getFullList({
                filter: ApiClient.filter("collectionRef = {:collectionId} && recordRef = {:recordId}", {
                    collectionId: record.collectionId,
                    recordId: record.id,
                }),
            });
        } catch (err) {
            ApiClient.error(err);
        }

        isLoading = false;
    }

    function unlinkExternalAuth(externalAuth) {
        if (!record?.id || !externalAuth) {
            return; // nothing to unlink
        }

        confirm(
            `确定要解除${getProviderTitle(externalAuth.provider)}提供商的关联吗？`,
            () => {
                return ApiClient.collection("_externalAuths")
                    .delete(externalAuth.id)
                    .then(() => {
                        addSuccessToast(
                            `成功解除${getProviderTitle(externalAuth.provider)}提供商的关联`,
                        );
                        dispatch("unlink", externalAuth.provider);
                        loadExternalAuths(); // reload list
                    })
                    .catch((err) => {
                        ApiClient.error(err);
                    });
            },
        );
    }

    loadExternalAuths();
</script>

{#if isLoading}
    <div class="block txt-center">
        <span class="loader" />
    </div>
{:else if record?.id && externalAuths.length}
    <div class="list">
        {#each externalAuths as auth}
            <div class="list-item">
                <figure class="provider-logo">
                    <img
                        src="{import.meta.env.BASE_URL}images/oauth2/{getProviderConfig(auth.provider)?.logo}"
                        alt="提供商标志"
                    />
                </figure>
                <span class="txt">{getProviderTitle(auth.provider)}</span>
                <div class="txt-hint">ID: {auth.providerId}</div>
                <button
                    type="button"
                    class="btn btn-transparent link-hint btn-circle btn-sm m-l-auto"
                    on:click={() => unlinkExternalAuth(auth)}
                >
                    <i class="ri-close-line" />
                </button>
            </div>
        {/each}
    </div>
{:else}
    <h6 class="txt-hint txt-center m-t-sm m-b-sm">没有关联的OAuth2提供商。</h6>
{/if}