<script>
    import { UpdateTaskDone, DeleteTask } from "../../wailsjs/go/main/App"
    import { createEventDispatcher } from "svelte"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import EditIcon from "~icons/material-symbols/edit"

    export let task

    const dispatch = createEventDispatcher()

    function deleteTask() {
        DeleteTask(task.ID).then(() => {
            dispatch("delete", task.ID)
        }).catch((error) => dispatch("error", error))
    }

    function toggleDone() {
        UpdateTaskDone(task.ID, !task.done).then(() => {
            dispatch("done", task.ID)
        }).catch((error) => dispatch("error", error))
    }
</script>

<li class="card flex items-center justify-between px-4 py-8">
    <input
        on:click={toggleDone}
        class="checkbox"
        type="checkbox"
        bind:checked={task.done}
    />
    <span class="flex-1 flex justify-center">
        {task.text}
    </span>
    <div class="flex">
        <button
            class:invisible={task.done}
            on:click={() => dispatch("edit", task)}
            type="button"
            class="btn-icon bg-initial"
        >
            <EditIcon style="color:gray" />
        </button>
        <button on:click={deleteTask} type="button" class="btn-icon bg-initial">
            <DeleteIcon style="color:gray" />
        </button>
    </div>
</li>
