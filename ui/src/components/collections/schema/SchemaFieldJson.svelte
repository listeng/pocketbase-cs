<script>
    import { slide } from "svelte/transition";
    import CommonHelper from "@/utils/CommonHelper";
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";

    export let field;
    export let key = "";

    let showInfo = false;

    $: if (CommonHelper.isEmpty(field.options)) {
        loadDefaults();
    }

    function loadDefaults() {
        field.options = {
            maxSize: 2000000,
        };
    }
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <Field class="form-field required m-b-sm" name="schema.{key}.options.maxSize" let:uniqueId>
            <label for={uniqueId}>最大长度 <small>(字节)</small></label>
            <input type="number" id={uniqueId} step="1" min="0" bind:value={field.options.maxSize} />
        </Field>

        <button
            type="button"
            class="btn btn-sm {showInfo ? 'btn-secondary' : 'btn-hint btn-transparent'}"
            on:click={() => {
                showInfo = !showInfo;
            }}
        >
            <strong class="txt">字符串标准化</strong>
            {#if showInfo}
                <i class="ri-arrow-up-s-line txt-sm" />
            {:else}
                <i class="ri-arrow-down-s-line txt-sm" />
            {/if}
        </button>
        {#if showInfo}
            <div class="block" transition:slide={{ duration: 150 }}>
                <div class="alert alert-warning m-b-0 m-t-10">
                    <div class="content">
                        为了无缝支持 <code>application/json</code> 和 <code>multipart/form-data</code> 请求，如果 <code>json</code> 字段是一个<strong>普通字符串</strong>，将应用以下标准化规则：
                        <ul>
                            <li>"true" 将转换为 json 的 <code>true</code></li>
                            <li>"false" 将转换为 json 的 <code>false</code></li>
                            <li>"null" 将转换为 json 的 <code>null</code></li>
                            <li>"[1,2,3]" 将转换为 json 的 <code>[1,2,3]</code></li>
                            <li>
                                {'"{"a":1,"b":2}"'} 将转换为 json 的 <code>{'{"a":1,"b":2}'}</code>
                            </li>
                            <li>数字字符串将转换为 json 的数字</li>
                            <li>双引号包裹的字符串将保持不变（即不进行标准化）</li>
                            <li>任何其他字符串（包括空字符串）将被加上双引号</li>
                        </ul>
                        另外，如果您想避免字符串值的标准化，可以将数据包装在一个对象内，例如 <code>{'{"data": anything}'}</code>。
                    </div>
                </div>
            </div>
        {/if}
    </svelte:fragment>
</SchemaField>
