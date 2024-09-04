<script>
    import { DeleteTask } from "../../wailsjs/go/main/App"
    import { createEventDispatcher } from "svelte"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import DragIcon from '~icons/material-symbols/drag-indicator'

    import Tag from "./Tag.svelte"

    import { parseTaskText } from "../lib.js"

    export let task

    let taskText = task.text
    let editing = false

    $: if (editing) {
        editable?.focus()
    }

    let editable

    const dispatch = createEventDispatcher()

    function handleEditableFocus() {
        if (!editable) {
            return
        }
        let range = document.createRange()
        range.selectNodeContents(editable)
        range.collapse(false)
        let selection = window.getSelection()
        selection.removeAllRanges()
        selection.addRange(range)
    }

    function saveEdits() {
        editing = false
        if (taskText === task.text) {
            return
        }
        task.text = taskText
        const parsed = parseTaskText(task.text)
        dispatch("edit", { task, parsed })
        task = {
            ...task,
            ...parsed,
        }
        task.tags = parsed.tags.map((t) => { return {name: t} })
    }

    function deleteTask() {
        DeleteTask(task.ID)
            .then(() => {
                dispatch("delete", task.ID)
            })
            .catch((error) => dispatch("error", error))
    }

    function handleKey(e) {
        if (e.key === "Enter") {
            saveEdits()
            editable?.blur()
        }
    }
</script>

<li class="group flex items-center justify-between gap-2">
    <div class="flex items-center w-5 h-5">
        <DragIcon style="color:gray" />
    </div>
    <div class="flex-1 flex gap-2 items-center">
        <input
            class="checkbox"
            type="checkbox"
            bind:checked={task.done}
            on:click={() => dispatch("done", { taskID: task.ID, done: !task.done })}
        />
        {#if !task.done && editing}
            <div
                class="flex-1"
                contenteditable
                bind:this={editable}
                bind:textContent={taskText}
                on:focusin={handleEditableFocus}
                on:focusout={saveEdits}
                on:keydown={handleKey}
            >
            </div>
        {:else}
            <div class="flex gap-2">
                <button on:click={() => editing = true} class="flex-1">{task.parsed_text || task.text}</button>
                {#if task.tags}
                    <div class="flex gap-1">
                        {#each task.tags as tag}
                            <Tag {tag} />
                        {/each}
                    </div>
                {/if}
            </div>
        {/if}
    </div>
    <div class="flex items-center w-4 h-4">
        <button on:click={deleteTask} class="btn-icon bg-initial hidden group-hover:block">
            <DeleteIcon style="color:gray" />
        </button>
    </div>
</li>
