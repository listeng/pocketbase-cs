<script>
    import CodeBlock from "@/components/base/CodeBlock.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let collection;

    let responseTab = 200;
    let responses = [];

    $: responses = [
        {
            code: 200,
            body: JSON.stringify(
                {
                    otpId: CommonHelper.randomString(15),
                },
                null,
                2,
            ),
        },
        {
            code: 400,
            body: `
                {
                  "status": 400,
                  "message": "验证提交数据时发生错误。",
                  "data": {
                    "email": {
                      "code": "validation_is_email",
                      "message": "必须是有效的电子邮件地址。"
                    }
                  }
                }
            `,
        },
        {
            code: 429,
            body: `
                {
                  "status": 429,
                  "message": "您发送了太多OTP请求，请稍后再试。",
                  "data": {}
                }
            `,
        },
    ];
</script>

<div class="alert alert-success">
    <strong class="label label-primary">POST</strong>
    <div class="content">
        <p>
            /api/collections/<strong>{collection.name}</strong>/request-otp
        </p>
    </div>
</div>

<div class="section-title">请求参数</div>
<table class="table-compact table-border m-b-base">
    <thead>
        <tr>
            <th>参数</th>
            <th>类型</th>
            <th width="50%">描述</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>
                <div class="inline-flex">
                    <span class="label label-success">必填</span>
                    <span>email</span>
                </div>
            </td>
            <td>
                <span class="label">字符串</span>
            </td>
            <td>要发送OTP请求的认证记录电子邮件地址（如果存在）。</td>
        </tr>
    </tbody>
</table>

<div class="section-title">响应</div>
<div class="tabs">
    <div class="tabs-header compact combined left">
        {#each responses as response (response.code)}
            <button
                class="tab-item"
                class:active={responseTab === response.code}
                on:click={() => (responseTab = response.code)}
            >
                {response.code}
            </button>
        {/each}
    </div>
    <div class="tabs-content">
        {#each responses as response (response.code)}
            <div class="tab-item" class:active={responseTab === response.code}>
                <CodeBlock content={response.body} />
            </div>
        {/each}
    </div>
</div>