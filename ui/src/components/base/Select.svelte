<script>
    import { onMount, createEventDispatcher } from "svelte";
    import tooltip from "@/actions/tooltip";
    import Toggler from "@/components/base/Toggler.svelte";
    import CommonHelper from "@/utils/CommonHelper";

    export let id = "";
    export let noOptionsText = "未找到选项";
    export let selectPlaceholder = "- 请选择 -";
    export let searchPlaceholder = "搜索...";
    export let items = [];
    export let multiple = false;
    export let disabled = false;
    export let readonly = false;
    export let upside = false;
    export let zeroFunc = () => (multiple ? [] : undefined);
    export let selected = zeroFunc();
    export let toggle = multiple; // 点击时切换选项
    export let closable = true; // 选择/取消选择选项时关闭下拉框
    export let labelComponent = undefined; // 用于每个选中选项标签的自定义组件
    export let labelComponentProps = {}; // 传递给自定义选项组件的属性
    export let optionComponent = undefined; // 用于每个下拉选项项的自定义组件
    export let optionComponentProps = {}; // 传递给自定义选项组件的属性
    export let searchable = false; // 是否显示下拉选项搜索输入框
    export let searchFunc = undefined; // 自定义搜索选项过滤器: `function(item, searchTerm):boolean`

    const dispatch = createEventDispatcher();

    let classes = "";
    export { classes as class }; // 导出保留关键字

    let toggler;
    let searchTerm = "";
    let container = undefined;
    let labelDiv = undefined;

    $: if (items) {
        ensureSelectedExist();
        resetSearch();
    }

    $: filteredItems = filterItems(items, searchTerm);

    $: isSelected = function (item) {
        const normalized = CommonHelper.toArray(selected);

        return CommonHelper.inArray(normalized, item);
    };

    // 选择处理函数
    // ---------------------------------------------------------------
    export function deselectItem(item) {
        if (CommonHelper.isEmpty(selected)) {
            return; // 没有可取消选择的项
        }

        let normalized = CommonHelper.toArray(selected);
        if (CommonHelper.inArray(normalized, item)) {
            CommonHelper.removeByValue(normalized, item);
            selected = multiple ? normalized : normalized?.[0] || zeroFunc();
        }

        dispatch("change", { selected });

        // 模拟原生change事件
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function selectItem(item) {
        if (multiple) {
            let normalized = CommonHelper.toArray(selected);
            if (!CommonHelper.inArray(normalized, item)) {
                selected = [...normalized, item];
            }
        } else {
            selected = item;
        }

        dispatch("change", { selected });

        // 模拟原生change事件
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function toggleItem(item) {
        return isSelected(item) ? deselectItem(item) : selectItem(item);
    }

    export function reset() {
        selected = zeroFunc();

        dispatch("change", { selected });

        // 模拟原生change事件
        container?.dispatchEvent(new CustomEvent("change", { detail: selected, bubbles: true }));
    }

    export function showDropdown() {
        toggler?.show && toggler?.show();
    }

    export function hideDropdown() {
        toggler?.hide && toggler?.hide();
    }

    function ensureSelectedExist() {
        if (CommonHelper.isEmpty(selected) || CommonHelper.isEmpty(items)) {
            return; // 无需检查
        }

        let selectedArray = CommonHelper.toArray(selected);
        let unselectedArray = [];

        // 查找缺失项
        for (const selectedItem of selectedArray) {
            if (!CommonHelper.inArray(items, selectedItem)) {
                unselectedArray.push(selectedItem);
            }
        }

        // 触发响应性
        if (unselectedArray.length) {
            for (const item of unselectedArray) {
                CommonHelper.removeByValue(selectedArray, item);
            }

            selected = multiple ? selectedArray : selectedArray[0];
        }
    }

    // 搜索处理函数
    // ---------------------------------------------------------------
    function defaultSearchFunc(item, search) {
        let normalizedSearch = ("" + search).replace(/\s+/g, "").toLowerCase();
        let normalizedItem = item;

        try {
            if (typeof item === "object" && item !== null) {
                normalizedItem = JSON.stringify(item);
            }
        } catch (e) {}

        return ("" + normalizedItem).replace(/\s+/g, "").toLowerCase().includes(normalizedSearch);
    }

    function resetSearch() {
        searchTerm = "";
    }

    function filterItems(items, search) {
        items = items || [];

        const filterFunc = searchFunc || defaultSearchFunc;

        return items.filter((item) => filterFunc(item, search)) || [];
    }

    // 选项操作
    // ---------------------------------------------------------------
    function handleOptionSelect(e, item) {
        e.preventDefault();

        if (toggle && multiple) {
            toggleItem(item);
        } else {
            selectItem(item);
        }
    }

    function handleOptionKeypress(e, item) {
        if (e.code === "Enter" || e.code === "Space") {
            handleOptionSelect(e, item);
            if (closable) {
                hideDropdown();
            }
        }
    }

    function onDropdownShow() {
        resetSearch();

        // 确保第一个选中的选项可见
        setTimeout(() => {
            const selected = container?.querySelector(".dropdown-item.option.selected");
            if (selected) {
                selected.focus();
                selected.scrollIntoView({ block: "nearest" });
            }
        }, 0);
    }

    // 标签激活
    // ---------------------------------------------------------------
    function onLabelClick(e) {
        e.stopPropagation();

        !readonly && !disabled && toggler?.toggle();
    }

    onMount(() => {
        const labels = document.querySelectorAll(`label[for="${id}"]`);

        for (const label of labels) {
            label.addEventListener("click", onLabelClick);
        }

        return () => {
            for (const label of labels) {
                label.removeEventListener("click", onLabelClick);
            }
        };
    });
</script>

<div bind:this={container} class="select {classes}" class:upside class:multiple class:disabled class:readonly>
    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
    <div
        bind:this={labelDiv}
        tabindex={disabled || readonly ? "-1" : "0"}
        class="selected-container"
        class:disabled
        class:readonly
        role="button"
    >
        {#each CommonHelper.toArray(selected) as item, i}
            <div class="option">
                {#if labelComponent}
                    <svelte:component this={labelComponent} {item} {...labelComponentProps} />
                {:else}
                    <span class="txt">{item}</span>
                {/if}

                {#if multiple || toggle}
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <!-- svelte-ignore a11y-no-static-element-interactions -->
                    <span
                        class="clear"
                        use:tooltip={"清除"}
                        on:click|preventDefault|stopPropagation={() => deselectItem(item)}
                    >
                        <i class="ri-close-line" />
                    </span>
                {/if}
            </div>
        {:else}
            <div class="block txt-placeholder" class:link-hint={!disabled && !readonly}>
                {selectPlaceholder}
            </div>
        {/each}
    </div>

    {#if !disabled && !readonly}
        <Toggler
            bind:this={toggler}
            class="dropdown dropdown-block options-dropdown dropdown-left {upside ? 'dropdown-upside' : ''}"
            trigger={labelDiv}
            on:show={onDropdownShow}
            on:hide
        >
            {#if searchable}
                <div class="form-field form-field-sm options-search">
                    <label class="input-group">
                        <div class="addon p-r-0">
                            <i class="ri-search-line" />
                        </div>
                        <!-- svelte-ignore a11y-autofocus -->
                        <input
                            autofocus
                            type="text"
                            placeholder={searchPlaceholder}
                            bind:value={searchTerm}
                        />

                        {#if searchTerm.length}
                            <div class="addon suffix p-r-5">
                                <button
                                    type="button"
                                    class="btn btn-sm btn-circle btn-transparent clear"
                                    on:click|preventDefault|stopPropagation={resetSearch}
                                >
                                    <i class="ri-close-line" />
                                </button>
                            </div>
                        {/if}
                    </label>
                </div>
            {/if}

            <slot name="beforeOptions" />

            <div class="options-list">
                {#each filteredItems as item}
                    <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
                    <!-- svelte-ignore a11y-no-static-element-interactions -->
                    <div
                        tabindex="0"
                        class="dropdown-item option"
                        role="menuitem"
                        class:closable
                        class:selected={isSelected(item)}
                        on:click={(e) => handleOptionSelect(e, item)}
                        on:keydown={(e) => handleOptionKeypress(e, item)}
                    >
                        {#if optionComponent}
                            <svelte:component this={optionComponent} {item} {...optionComponentProps} />
                        {:else}{item}{/if}
                    </div>
                {:else}
                    {#if noOptionsText}
                        <div class="txt-missing">{noOptionsText}</div>
                    {/if}
                {/each}
            </div>

            <slot name="afterOptions" />
        </Toggler>
    {/if}
</div>