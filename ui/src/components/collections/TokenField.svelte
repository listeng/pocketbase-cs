<script>
    import Field from "@/components/base/Field.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let key;
    export let label;
    export let duration;
    export let secret;
</script>

<Field class="form-field required" name="{key}.duration" let:uniqueId>
    <label for={uniqueId}>{label} 有效期（秒）</label>
    <input type="number" id={uniqueId} required bind:value={duration} placeholder="保持不变" />
    <div class="help-block">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <!-- svelte-ignore a11y-no-static-element-interactions -->
        <span
            class="link-primary"
            class:txt-success={!!secret}
            on:click={() => {
                // toggle
                if (secret) {
                    secret = undefined;
                } else {
                    secret = CommonHelper.randomSecret(50);
                }
            }}
        >
            使之前颁发的所有令牌失效
        </span>
    </div>
</Field>