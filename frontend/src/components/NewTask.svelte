<script>
    import { createEventDispatcher } from "svelte"
    import AddIcon from "~icons/material-symbols/add"

    import Checkbox from "./Checkbox.svelte"

    const dispatch = createEventDispatcher()

    let taskText = ""

    function addTask() {
        if (!taskText) {
            return
        }
        dispatch("add", taskText)
        taskText = ""
    }

    function handleKey(e) {
        if (e.key === "Enter") {
            addTask()
        }
    }
</script>

<li class="group flex items-center justify-between gap-2 my-1">
    <button class="w-5 h-5" on:click={addTask}>
        <AddIcon style="color:gray" />
    </button>
    <div class="flex-1 flex gap-2 items-center">
        <Checkbox hover={false} disabled />
        <input class="flex-1 bg-transparent" type="text" bind:value={taskText} on:keydown={handleKey} />
    </div>
    <div class="w-4 h-4"></div>
</li>
