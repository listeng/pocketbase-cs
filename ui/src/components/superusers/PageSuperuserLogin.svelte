<script>
    import { link, replace, querystring } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import FullPage from "@/components/base/FullPage.svelte";
    import Field from "@/components/base/Field.svelte";
    import { setErrors } from "@/stores/errors";
    import { addErrorToast, removeAllToasts } from "@/stores/toasts";

    const queryParams = new URLSearchParams($querystring);

    let identity = queryParams.get("demoEmail") || "";
    let password = queryParams.get("demoPassword") || "";

    let authMethods = {};
    let currentStep = 1;
    let totalSteps = 1;

    let passwordAuthSubmitting = false;
    let otpRequestSubmitting = false;
    let otpAuthSubmitting = false;
    let isLoading = false;

    let mfaId = "";
    let otpId = "";
    let lastOTPId = "";
    let otpEmail = "";
    let otpPassword = "";

    $: {
        totalSteps = 1;
        currentStep = 1;

        if (authMethods?.mfa?.enabled) {
            totalSteps++;
        }

        if (authMethods?.otp?.enabled) {
            totalSteps++;
        }

        if (mfaId != "") {
            currentStep++;
        }

        if (otpId != "") {
            currentStep++;
        }
    }

    loadAuthMethods();

    async function loadAuthMethods() {
        if (isLoading) {
            return;
        }

        isLoading = true;

        try {
            authMethods = await ApiClient.collection("_superusers").listAuthMethods();
        } catch (err) {
            ApiClient.error(err);
        }

        if (authMethods.cas.enabled) {
            const ps = getUrlParams(cleanUrl(window.location.href).url);
            if (ps['iscas'] == 1 && ps['ticket']) {
                validUserinfo(ps['ticket'])
            }
        }

        isLoading = false;
    }

    async function authWithPassword() {
        if (passwordAuthSubmitting) {
            return;
        }

        passwordAuthSubmitting = true;

        try {
            await ApiClient.collection("_superusers").authWithPassword(identity, password);
            removeAllToasts();
            setErrors({});
            replace("/");
        } catch (err) {
            if (err.status == 401) {
                mfaId = err.response.mfaId;

                // show the otp forms
                if (
                    // if the identity field is just the email use it to directly send an otp request
                    authMethods?.password?.identityFields?.length == 1 &&
                    authMethods.password.identityFields[0] == "email"
                ) {
                    // prefill and request
                    otpEmail = identity;
                    await requestOTP();
                } else if (/^[^@\s]+@[^@\s]+$/.test(identity)) {
                    // only prefill
                    otpEmail = identity;
                }
            } else if (err.status != 400) {
                ApiClient.error(err);
            } else {
                addErrorToast("无效的登录凭证。");
            }
        }

        passwordAuthSubmitting = false;
    }

    async function requestOTP() {
        if (otpRequestSubmitting) {
            return;
        }

        otpRequestSubmitting = true;

        try {
            const result = await ApiClient.collection("_superusers").requestOTP(otpEmail);
            otpId = result.otpId;
            lastOTPId = otpId;
            removeAllToasts();
            setErrors({});
        } catch (err) {
            // reset the form
            if (err.status == 429) {
                otpId = lastOTPId;
            }

            ApiClient.error(err);
        }

        otpRequestSubmitting = false;
    }

    async function authWithOTP() {
        if (otpAuthSubmitting) {
            return;
        }

        otpAuthSubmitting = true;

        try {
            await ApiClient.collection("_superusers").authWithOTP(otpId || lastOTPId, otpPassword, { mfaId });
            removeAllToasts();
            setErrors({});
            replace("/");
        } catch (err) {
            ApiClient.error(err);
        }

        otpAuthSubmitting = false;
    }

    async function validUserinfo(ticket) {
        try {
            const service = localStorage.getItem('cas-service');
            const response = await fetch(authMethods.cas.callbackUrl + '?ticket=' + 
                ticket + '&service=' + encodeURIComponent(service));
            
            if (response.status === 200) {
                const result = await response.json();

                localStorage.setItem('__pb_superuser_auth__', JSON.stringify(result));
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
        const loginUrl = authMethods.cas.loginUrl + '?service=' + encodeURIComponent(currentUrl);
        localStorage.setItem('cas-service', currentUrl);

        window.location.href = loginUrl;
    }
</script>

<FullPage>
    <div class="content txt-center m-b-base">
        <h4>
            超级用户登录
            {#if totalSteps > 1}
                ({currentStep}/{totalSteps})
            {/if}
        </h4>
    </div>

    {#if isLoading}
        <div class="block txt-center">
            <span class="loader" />
        </div>
    {:else if authMethods.password.enabled && !mfaId}
        <!-- auth with password -->

        {#if !authMethods.cas.onlyCasLogin}
        <form class="block" on:submit|preventDefault={authWithPassword}>
            <Field class="form-field required" name="identity" let:uniqueId>
                <label for={uniqueId}>
                    {CommonHelper.sentenize(authMethods.password.identityFields.join(" 或 "), false)}
                </label>
                <!-- svelte-ignore a11y-autofocus -->
                <input
                    id={uniqueId}
                    type={authMethods.password.identityFields.length == 1 &&
                    authMethods.password.identityFields[0] == "email"
                        ? "email"
                        : "text"}
                    value={identity}
                    on:input={(e) => {
                        identity = e.target.value;
                    }}
                    required
                    autofocus
                />
            </Field>

            <Field class="form-field required" name="password" let:uniqueId>
                <label for={uniqueId}>密码</label>
                <input type="password" id={uniqueId} bind:value={password} required />
                <div class="help-block">
                    <a href="/request-password-reset" class="link-hint" use:link>忘记密码？</a>
                </div>
            </Field>

            <button
                type="submit"
                class="btn btn-lg btn-block btn-next"
                class:btn-disabled={passwordAuthSubmitting}
                class:btn-loading={passwordAuthSubmitting}
            >
                <span class="txt">{totalSteps > 1 ? "下一步" : "登录"}</span>
                <i class="ri-arrow-right-line" />
            </button>
        </form>
        {/if}

        {#if authMethods.cas.enabled}
            <button
                on:click|preventDefault|stopPropagation={ssologin}
                class="btn btn-lg btn-block btn-next"
                style="margin-top: 20px; background-color: #4BABD3"
            >
                <span class="txt">{authMethods.cas.displayName}</span>
                <i class="ri-arrow-right-line" />
            </button>
        {/if}
    {:else if authMethods.otp.enabled}
        {#if !otpId}
            <!-- request otp -->
            <form class="block" on:submit|preventDefault={requestOTP}>
                <Field class="form-field required" name="email" let:uniqueId>
                    <label for={uniqueId}>邮箱</label>
                    <input type="email" id={uniqueId} bind:value={otpEmail} required />
                </Field>

                <button
                    type="submit"
                    class="btn btn-lg btn-block btn-next"
                    class:btn-disabled={otpRequestSubmitting}
                    class:btn-loading={otpRequestSubmitting}
                >
                    <i class="ri-mail-send-line" />
                    <span class="txt">发送OTP</span>
                </button>
            </form>
        {:else}
            {#if otpEmail}
                <div class="content txt-center m-b-sm">
                    <p>
                        请检查您的 <strong>{otpEmail}</strong> 邮箱，并在下方输入收到的一次性密码(OTP)。
                    </p>
                </div>
            {/if}

            <!-- auth with otp -->
            <form class="block" on:submit|preventDefault={authWithOTP}>
                <Field class="form-field required" name="otpId" let:uniqueId>
                    <label for={uniqueId}>ID</label>
                    <input
                        type="text"
                        id={uniqueId}
                        value={otpId}
                        placeholder={lastOTPId}
                        on:change={(e) => {
                            otpId = e.target.value || lastOTPId;
                            e.target.value = otpId;
                        }}
                        required
                    />
                </Field>

                <Field class="form-field required" name="password" let:uniqueId>
                    <label for={uniqueId}>一次性密码</label>
                    <!-- svelte-ignore a11y-autofocus -->
                    <input type="password" id={uniqueId} bind:value={otpPassword} required autofocus />
                </Field>

                <button
                    type="submit"
                    class="btn btn-lg btn-block btn-next"
                    class:btn-disabled={otpAuthSubmitting}
                    class:btn-loading={otpAuthSubmitting}
                >
                    <span class="txt">登录</span>
                    <i class="ri-arrow-right-line" />
                </button>
            </form>

            <div class="content txt-center m-t-sm">
                <button
                    type="button"
                    class="link-hint"
                    disabled={otpAuthSubmitting}
                    on:click={() => {
                        otpId = "";
                    }}
                >
                    重新请求OTP
                </button>
            </div>
        {/if}
    {/if}
</FullPage>