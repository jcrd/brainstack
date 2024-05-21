<script>
    import { DeleteTask } from "../../wailsjs/go/main/App"
    import { createEventDispatcher } from "svelte"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import EditIcon from "~icons/material-symbols/edit"
    import UndoIcon from "~icons/material-symbols/undo-rounded"
    import PromoteIcon from "~icons/material-symbols/arrow-upward"

    export let task

    const dispatch = createEventDispatcher()

    function deleteTask() {
        DeleteTask(task.ID)
            .then(() => {
                dispatch("delete", task.ID)
            })
            .catch((error) => dispatch("error", error))
    }
</script>

<li class="card flex items-center px-4 py-8">
    <span class="flex-1">
        {task.text}
    </span>
    <div class="flex">
        {#if task.done}
            <button
                on:click={() =>
                    dispatch("done", { taskID: task.ID, done: false })}
                class="btn-icon bg-initial"
            >
                <UndoIcon style="color:gray;font-size:1.25em" />
            </button>
        {:else}
            <button
                on:click={() => dispatch("promote", task)}
                class="btn-icon bg-initial"
            >
                <PromoteIcon style="color:gray" />
            </button>
            <button
                on:click={() => dispatch("edit", task)}
                class="btn-icon bg-initial"
            >
                <EditIcon style="color:gray" />
            </button>
        {/if}
        <button on:click={deleteTask} class="btn-icon bg-initial">
            <DeleteIcon style="color:gray" />
        </button>
    </div>
</li>
