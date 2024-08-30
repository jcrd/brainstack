<script>
    // @ts-nocheck

    import {
        AddTask,
        EditTask,
        ReorderTasks,
        UpdateTaskDone,
        EditStack,
        DeleteStack,
    } from "../../wailsjs/go/main/App"

    import { createEventDispatcher } from "svelte"
    import { writable } from "svelte/store"

    import {
        dndzone,
        overrideItemIdKeyNameBeforeInitialisingDndZones,
    } from "svelte-dnd-action"

    import TabTodoIcon from '~icons/material-symbols/circle-outline'
    import TabDoneIcon from '~icons/material-symbols/check-circle'

    import Task from "./Task.svelte"
    import NewTask from "./NewTask.svelte"

    export let stack

    overrideItemIdKeyNameBeforeInitialisingDndZones("ID")

    const dispatch = createEventDispatcher()
    const tabSet = writable(0)

    $: {
        stack
        $tabSet = 0
    }

    $: tasks = stack.tasks?.sort((a, b) => a.order - b.order) || []
    $: tasksDone = tasks.filter((t) => t.done)
    $: tasksTodo = tasks.filter((t) => !t.done)

    $: {
        tasks
        dispatch("invalidate", stack.ID)
    }

    function addTask({ detail: text }) {
        const order =
            tasks.length > 0 ? tasks[tasks.length - 1].order + 1 : 0
        AddTask(stack.ID, text, order)
            .then((taskID) => {
                tasks = [
                    ...tasks,
                    {
                        ID: taskID,
                        stack_id: stack.ID,
                        order,
                        text,
                    },
                ]
            })
            .catch((error) => dispatch("error", error))
    }

    function taskEdit({ detail: task }) {
        EditTask(task.ID, task.text)
            .catch((error) => dispatch("error", error))
    }

    function reorderTasks(items) {
        items = items.map((task, i) => {
            task.order = i
            return task
        })
        ReorderTasks(items)
            .then((result) => {
                tasks = [...result, ...tasksDone]
            })
            .catch((error) => dispatch("error", error))
    }

    function taskDeleted({ detail: taskID }) {
        tasks = tasks.filter((t) => t.ID != taskID)
    }

    function taskDone({ detail: { taskID, done } }) {
        UpdateTaskDone(taskID, done)
            .then(() => {
                tasks = tasks.map((t) => {
                    if (t.ID === taskID) {
                        t.done = done
                    }
                    return t
                })
            })
            .catch((error) => dispatch("error", error))
    }

    function taskPromoted({ detail: task }) {
        const items = [task, ...tasksTodo.filter((t) => t.ID != task.ID)]
        reorderTasks(items)
    }

    function handleDndConsider(e) {
        tasksTodo = e.detail.items
    }

    function handleDndFinalize(e) {
        reorderTasks(e.detail.items)
    }
</script>

<div class="flex flex-col gap-4 mx-8 max-h-screen overflow-auto pb-20 p-1">
    <button class="pl-7 py-1 flex items-center gap-2 hover:bg-surface-200 hover:text-surface-600" on:click={() => $tabSet = !$tabSet}>
        {#if $tabSet == 0}
            <span class="mt-[3px]"><TabTodoIcon /></span>
            <span class="text-xl">Todo</span>
        {:else if $tabSet == 1}
            <span class="mt-[3px]"><TabDoneIcon /></span>
            <span class="text-xl">Done</span>
        {/if}
    </button>
    {#if $tabSet == 0 && tasksTodo.length}
        <ul
            use:dndzone={{ items: tasksTodo }}
            on:consider={handleDndConsider}
            on:finalize={handleDndFinalize}
            class="flex flex-col gap-4"
        >
            {#each tasksTodo as task (task.ID)}
                <Task
                    on:promote={taskPromoted}
                    on:delete={taskDeleted}
                    on:edit={taskEdit}
                    on:done={taskDone}
                    on:error={(error) => dispatch("error", error)}
                    {task}
                />
            {/each}
        </ul>
        <ul>
            <NewTask
                on:add={addTask}
            />
        </ul>
    {:else if $tabSet == 1 && tasksDone.length}
        <ul class="flex flex-col gap-4">
            {#each tasksDone as task (task.ID)}
                <Task
                    on:delete={taskDeleted}
                    on:done={taskDone}
                    on:error={(error) => dispatch("error", error)}
                    {task}
                />
            {/each}
        </ul>
    {/if}
</div>
