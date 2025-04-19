<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let index;
    export let list = [];
    export let group = "默认";
    export let disabled = false;
    export let dragHandleClass = ""; // 默认整个元素都可拖动

    let dragging = false;
    let dragover = false;

    function onDrag(e, i) {
        if (!e || disabled) {
            return;
        }

        if (dragHandleClass && !e.target.classList.contains(dragHandleClass)) {
            // 不是拖动把手
            dragover = false;
            dragging = false;
            e.preventDefault();
            return;
        }

        dragging = true;

        e.dataTransfer.effectAllowed = "移动";
        e.dataTransfer.dropEffect = "移动";
        e.dataTransfer.setData(
            "text/plain",
            JSON.stringify({
                index: i,
                group: group,
            }),
        );

        dispatch("drag", e);
    }

    function onDrop(e, target) {
        dragover = false;
        dragging = false;

        if (!e || disabled) {
            return;
        }

        e.dataTransfer.dropEffect = "移动";

        let dragData = {};
        try {
            dragData = JSON.parse(e.dataTransfer.getData("text/plain"));
        } catch (_) {}

        if (dragData.group != group) {
            return; // 不同的可拖动组
        }

        const start = dragData.index << 0;

        if (start < target) {
            list.splice(target + 1, 0, list[start]);
            list.splice(start, 1);
        } else {
            list.splice(target, 0, list[start]);
            list.splice(start + 1, 1);
        }

        list = list;

        dispatch("sort", {
            oldIndex: start,
            newIndex: target,
            list: list,
        });
    }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    draggable={!disabled}
    class="draggable"
    class:dragging
    class:dragover
    on:dragover|preventDefault={() => {
        dragover = true;
    }}
    on:dragleave|preventDefault={() => {
        dragover = false;
    }}
    on:dragend={() => {
        dragover = false;
        dragging = false;
    }}
    on:dragstart={(e) => onDrag(e, index)}
    on:drop={(e) => onDrop(e, index)}
>
    <slot {dragging} {dragover} />
</div>

<style>
    .draggable {
        user-select: text;
        outline: 0;
        min-width: 0;
    }
</style>