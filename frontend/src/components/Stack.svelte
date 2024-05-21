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
    import {
        dndzone,
        overrideItemIdKeyNameBeforeInitialisingDndZones,
    } from "svelte-dnd-action"
    import Task from "./Task.svelte"
    import ButtonBar from "./ButtonBar.svelte"

    export let stack
    export let modal

    overrideItemIdKeyNameBeforeInitialisingDndZones("ID")

    const dispatch = createEventDispatcher()

    $: tasks = stack.tasks?.sort((a, b) => a.order - b.order) || []
    $: tasksDone = tasks.filter((t) => t.done)
    $: tasksTodo = tasks.filter((t) => !t.done)

    $: {
        tasks
        dispatch("invalidate", stack.ID)
    }

    function addTask() {
        modal.trigger({
            type: "prompt",
            title: "Add task",
            valueAttr: { type: "text", required: true },
            response: (text) => {
                if (!text) {
                    return
                }
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
            },
        })
    }

    function taskEdit({ detail: task }) {
        modal.trigger({
            type: "prompt",
            title: "Edit task",
            value: task.text,
            valueAttr: { type: "text", required: true },
            response: (text) => {
                if (!text) {
                    return
                }
                EditTask(task.ID, text)
                    .then(() => {
                        tasks = tasks.map((t) => {
                            if (t.ID == task.ID) {
                                t.text = text
                            }
                            return t
                        })
                    })
                    .catch((error) => dispatch("error", error))
            },
        })
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

    function nextTask() {
        if (tasksTodo.length > 0) {
            taskDone({ detail: { taskID: tasksTodo[0].ID, done: true } })
        }
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

    function editStack() {
        modal.trigger({
            type: "prompt",
            title: "Edit stack",
            value: stack.name,
            valueAttr: { type: "text", required: true },
            response: (text) => {
                if (!text) {
                    return
                }
                EditStack(stack.ID, text)
                    .then(() => {
                        dispatch("edit", {
                            stackID: stack.ID,
                            name: text,
                        })
                    })
                    .catch((error) => dispatch("error", error))
            },
        })
    }

    function deleteStack() {
        modal.trigger({
            type: "confirm",
            title: "Are you sure you want to delete this stack?",
            response: (state) => {
                if (state) {
                    DeleteStack(stack.ID)
                        .then(() => {
                            dispatch("delete", stack.ID)
                        })
                        .catch((error) => dispatch("error", error))
                }
            },
        })
    }
</script>

<div class="flex flex-col gap-4 mx-8 max-h-screen overflow-auto pb-20 p-1">
    <ButtonBar data={{ addTask, nextTask, editStack, deleteStack }} />
    {#if tasksTodo.length}
        <h1 class="badge variant-filled w-min">Todo</h1>
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
    {/if}
    {#if tasksDone.length}
        <h1 class="badge variant-filled w-min">Done</h1>
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
