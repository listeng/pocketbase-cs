<script>
    let expanded = false;

    function toggle() {
        expanded = !expanded;
    }
</script>

<button class="btn btn-sm btn-secondary m-t-10" on:click={toggle}>
    {#if expanded}
        <span class="txt">隐藏详情</span>
        <i class="ri-arrow-up-s-line" />
    {:else}
        <span class="txt">显示详情</span>
        <i class="ri-arrow-down-s-line" />
    {/if}
</button>

{#if expanded}
    <p>
        语法基本遵循以下格式
        <code>
            <span class="txt-success">操作数</span>
            <span class="txt-danger">操作符</span>
            <span class="txt-success">操作数</span></code
        >，其中：
    </p>
    <ul>
        <li>
            <code class="txt-success">操作数</code> - 可以是上述任何字段字面量、字符串（单引号或双引号）、数字、null、true、false
        </li>
        <li>
            <code class="txt-danger">操作符</code> - 包括：
            <br />
            <ul>
                <li>
                    <code class="filter-op">{"="}</code>
                    <span class="txt">等于</span>
                </li>
                <li>
                    <code class="filter-op">{"!="}</code>
                    <span class="txt">不等于</span>
                </li>
                <li>
                    <code class="filter-op">{">"}</code>
                    <span class="txt">大于</span>
                </li>
                <li>
                    <code class="filter-op">{">="}</code>
                    <span class="txt">大于或等于</span>
                </li>
                <li>
                    <code class="filter-op">{"<"}</code>
                    <span class="txt">小于</span>
                </li>
                <li>
                    <code class="filter-op">{"<="}</code>
                    <span class="txt">小于或等于</span>
                </li>
                <li>
                    <code class="filter-op">{"~"}</code>
                    <span class="txt">
                        相似/包含（如果未指定，会自动在右侧字符串操作数前后添加"%"进行通配符匹配）
                    </span>
                </li>
                <li>
                    <code class="filter-op">{"!~"}</code>
                    <span class="txt">
                        不相似/不包含（如果未指定，会自动在右侧字符串操作数前后添加"%"进行通配符匹配）
                    </span>
                </li>

                <li>
                    <code class="filter-op">{"?="}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">等于</span>
                </li>
                <li>
                    <code class="filter-op">{"?!="}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">不等于</span>
                </li>
                <li>
                    <code class="filter-op">{"?>"}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">大于</span>
                </li>
                <li>
                    <code class="filter-op">{"?>="}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">大于或等于</span>
                </li>
                <li>
                    <code class="filter-op">{"?<"}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">小于</span>
                </li>
                <li>
                    <code class="filter-op">{"?<="}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">小于或等于</span>
                </li>
                <li>
                    <code class="filter-op">{"?~"}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">
                        相似/包含（如果未指定，会自动在右侧字符串操作数前后添加"%"进行通配符匹配）
                    </span>
                </li>
                <li>
                    <code class="filter-op">{"?!~"}</code>
                    <em class="txt-hint">任意/至少一个</em>
                    <span class="txt">
                        不相似/不包含（如果未指定，会自动在右侧字符串操作数前后添加"%"进行通配符匹配）
                    </span>
                </li>
            </ul>
        </li>
    </ul>

    <p>
        要分组和组合多个表达式，可以使用括号
        <code>(...)</code>、<code>&&</code>（与）和<code>||</code>（或）标记。
    </p>
{/if}

<style>
    .filter-op {
        display: inline-block;
        vertical-align: top;
        margin-right: 5px;
        width: 30px;
        text-align: center;
        padding-left: 0;
        padding-right: 0;
    }
</style>