<script>
    import Field from "@/components/base/Field.svelte";
    import SchemaField from "@/components/collections/schema/SchemaField.svelte";
    import { slide } from "svelte/transition";

    export let field;
    export let key = "";

    let showInfo = false;
</script>

<SchemaField bind:field {key} on:rename on:remove on:duplicate {...$$restProps}>
    <svelte:fragment slot="options">
        <Field class="form-field m-b-sm" name="fields.{key}.maxSize" let:uniqueId>
            <label for={uniqueId}>最大大小 <small>(字节)</small></label>
            <input
                type="number"
                id={uniqueId}
                step="1"
                min="0"
                max={Number.MAX_SAFE_INTEGER}
                value={field.maxSize || ""}
                on:input={(e) => (field.maxSize = parseInt(e.target.value, 10))}
                placeholder="默认约5MB"
            />
        </Field>

        <button
            type="button"
            class="btn btn-sm {showInfo ? 'btn-secondary' : 'btn-hint btn-transparent'}"
            on:click={() => {
                showInfo = !showInfo;
            }}
        >
            <strong class="txt">字符串值规范化</strong>
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
                        为了同时无缝支持<code>application/json</code>和
                        <code>multipart/form-data</code>
                        请求，当<code>json</code>字段是
                        <strong>纯字符串</strong>时，会应用以下规范化规则：
                        <ul>
                            <li>"true"会被转换为json的<code>true</code></li>
                            <li>"false"会被转换为json的<code>false</code></li>
                            <li>"null"会被转换为json的<code>null</code></li>
                            <li>"[1,2,3]"会被转换为json的<code>[1,2,3]</code></li>
                            <li>
                                {'"{"a":1,"b":2}"'}会被转换为json的<code>{'{"a":1,"b":2}'}</code>
                            </li>
                            <li>数字字符串会被转换为json数字</li>
                            <li>双引号字符串会保持原样（即不进行规范化）</li>
                            <li>任何其他字符串（包括空字符串）都会被加上双引号</li>
                        </ul>
                        或者，如果你想避免字符串值规范化，可以将数据包装在对象中，例如<code>{'{"data": 任意内容}'}</code>
                    </div>
                </div>
            </div>
        {/if}
    </svelte:fragment>
</SchemaField>