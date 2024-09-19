<script>
    import { onMount } from 'svelte';
    import { link, replace, querystring } from "svelte-spa-router";
    import FullPage from "@/components/base/FullPage.svelte";
    import ApiClient from "@/utils/ApiClient";
    import Field from "@/components/base/Field.svelte";
    import { addErrorToast, removeAllToasts } from "@/stores/toasts";

    const queryParams = new URLSearchParams($querystring);

    let email = queryParams.get("demoEmail") || "";
    let password = queryParams.get("demoPassword") || "";
    let isLoading = false;
    let isSSOLoading = false;
    let ssoCfg = null;
    let showSSOLoginButton = false;

    onMount(() => {
        initLogin();
    });

    function login() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        return ApiClient.admins
            .authWithPassword(email, password)
            .then(() => {
                removeAllToasts();
                replace("/");
            })
            .catch(() => {
                addErrorToast("Invalid login credentials.");
            })
            .finally(() => {
                isLoading = false;
            });
    }

    async function initLogin() {
        const config_data = await ApiClient.collection('config').getList(1, 20, {
            filter: 'code = "cas-config"'
        });

        if (config_data.totalItems > 0) {
            const cfg = config_data.items[0].config;
            if (cfg['cas-url'] && cfg['valid-url'] && cfg['callback-url'] && cfg['realm']) {
                showSSOLoginButton = true;
                ssoCfg = cfg;

                const ps = getUrlParams(cleanUrl(window.location.href).url);
                if (ps['iscas'] == 1 && ps['ticket']) {
                    isSSOLoading = true;
                    validUserinfo(ps['ticket'])
                }
            }
        }
    }

    async function validUserinfo(ticket) {
        try {
            const service = localStorage.getItem('cas-service');
            const response = await fetch(ssoCfg['callback-url'] + '?ticket=' + 
                ticket + '&service=' + encodeURIComponent(service));
            
            if (response.status === 200) {
                const result = await response.json();

                localStorage.setItem('pb_admin_auth', JSON.stringify(result));
                removeAllToasts();

                const targetUrlData = localStorage.getItem('cas-service');
                const pps = getUrlParams(targetUrlData);
                let targetUrl = targetUrlData.split('?')[0];

                if (!targetUrl.endsWith('/')) {
                    targetUrl += '/';
                }

                window.location.href = targetUrl + '#' + pps['hashurl'];
            }

        } catch (error) {
            console.error('Error fetching or processing the request:', error);
            return false;
        }
    }

    function cleanUrl(url) {
        const urlObject = new URL(url);
        let hash = urlObject.hash;
        let cleanHash = hash.split('?')[0];

        if (hash.includes('?')) {
            const hashQueryParams = hash.split('?')[1];

            if (hashQueryParams) {
                urlObject.search = urlObject.search ? urlObject.search + '&' + hashQueryParams : '?' + hashQueryParams;
            }
        }

        urlObject.hash = '';

        return {
            url: urlObject.toString(),
            hash: cleanHash ? cleanHash.slice(1) : ''
        };
    }

    function getUrlParams(cleanedUrl) {
        const urlObject = new URL(cleanedUrl);
        const params = new URLSearchParams(urlObject.search);
        let paramsObj = {};
        for (let [key, value] of params.entries()) {
            paramsObj[key] = value;
        }

        return paramsObj;
    }

    function ssologin() {
        let result = cleanUrl(window.location.href);
        let currentUrl = result.url;
        let hash = result.hash;
        if (currentUrl.indexOf('?') >= 0) {
            currentUrl += '&iscas=1&hashurl=' + hash;
        } else {
            currentUrl += '?iscas=1&hashurl=' + hash;
        }
        const loginUrl = ssoCfg['cas-url'] + '/login?service=' + encodeURIComponent(currentUrl);
        localStorage.setItem('cas-service', currentUrl);

        window.location.href = loginUrl;
    }
</script>

<FullPage>
    {#if isSSOLoading}
    <div class="content txt-center m-b-base">
        <h4>正在登录……</h4>
    </div>
    {:else}
    <form class="block" on:submit|preventDefault={login} >
        <div class="content txt-center m-b-base">
            <h4>管理员登录</h4>
        </div>

        <Field class="form-field required" name="identity" let:uniqueId>
            <label for={uniqueId}>邮箱</label>
            <!-- svelte-ignore a11y-autofocus -->
            <input type="email" id={uniqueId} bind:value={email} required autofocus />
        </Field>

        <Field class="form-field required" name="password" let:uniqueId>
            <label for={uniqueId}>密码</label>
            <input type="password" id={uniqueId} bind:value={password} required />
            <div class="help-block">
                <a href="/request-password-reset" class="link-hint" use:link>忘记密码?</a>
            </div>
        </Field>

        <button
            type="submit"
            class="btn btn-lg btn-block btn-next"
            class:btn-disabled={isLoading}
            class:btn-loading={isLoading}
            
        >
            <span class="txt">登录</span>
            <i class="ri-arrow-right-line" />
        </button>

        {#if showSSOLoginButton}
        <button
            on:click={ssologin}
            class="btn btn-lg btn-block btn-next"
            style="margin-top: 20px; background-color: #4BABD3"
        >
            <span class="txt">单点登录</span>
            <i class="ri-arrow-right-line" />
        </button>
        {/if}
    </form>
    {/if}
</FullPage>
