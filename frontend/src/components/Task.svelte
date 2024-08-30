<script>
    import { DeleteTask } from "../../wailsjs/go/main/App"
    import { createEventDispatcher } from "svelte"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import DragIcon from '~icons/material-symbols/drag-indicator'

    export let task

    let taskText = task.text

    const dispatch = createEventDispatcher()

    $: if (taskText && taskText !== task.text) {
        task.text = taskText
        dispatch("edit", task)
    }

    function deleteTask() {
        DeleteTask(task.ID)
            .then(() => {
                dispatch("delete", task.ID)
            })
            .catch((error) => dispatch("error", error))
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
        {#if !task.done}
            <div class="flex-1" contenteditable bind:textContent={taskText}>{task.text}</div>
        {:else}
            <div class="flex-1">{task.text}</div>
        {/if}
    </div>
    <div class="flex items-center w-4 h-4">
        <button on:click={deleteTask} class="btn-icon bg-initial hidden group-hover:block">
            <DeleteIcon style="color:gray" />
        </button>
    </div>
</li>
