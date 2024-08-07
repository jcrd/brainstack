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

<li class="card flex items-center justify-between gap-2 px-4 py-1">
    <div class="flex items-center">
        <DragIcon style="color:gray" />
    </div>
    <div class="flex-1 flex gap-2 items-center">
        <input
            class="checkbox"
            type="checkbox"
            bind:checked={task.done}
            on:click={() => dispatch("done", { taskID: task.ID, done: !task.done })}
        />
        <div class="flex-1" contenteditable="true" bind:textContent={taskText}>{task.text}</div>
    </div>
    <div class="flex">
        <button on:click={deleteTask} class="btn-icon bg-initial">
            <DeleteIcon style="color:gray" />
        </button>
    </div>
</li>
