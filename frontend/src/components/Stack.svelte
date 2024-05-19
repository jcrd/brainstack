<script>
    // @ts-nocheck

    import {
        AddTask,
        EditTask,
        ReorderTasks,
        EditStack,
        DeleteStack,
    } from "../../wailsjs/go/main/App"
    import { createEventDispatcher } from "svelte"
    import {
        dndzone,
        overrideItemIdKeyNameBeforeInitialisingDndZones,
    } from "svelte-dnd-action"
    import { popup } from "@skeletonlabs/skeleton"
    import AddIcon from "~icons/material-symbols/add"
    import MenuIcon from "~icons/material-symbols/menu-rounded"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import EditIcon from "~icons/material-symbols/edit"
    import Task from "./Task.svelte"

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
                AddTask(stack.ID, text, order).then((taskID) => {
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
                EditTask(task.ID, text).then(() => {
                    tasks = tasks.map((t) => {
                        if (t.ID == task.ID) {
                            t.text = text
                        }
                        return t
                    })
                })
            },
        })
    }

    function taskDeleted({ detail: taskID }) {
        tasks = tasks.filter((t) => t.ID != taskID)
    }

    function taskDone() {
        tasks = tasks
    }

    function handleDndConsider(e) {
        tasksTodo = e.detail.items
    }

    function handleDndFinalize(e) {
        const items = e.detail.items.map((task, i) => {
            task.order = i
            return task
        })
        ReorderTasks(items).then((result) => {
            tasks = [...result, ...tasksDone]
        })
    }

    const stackPopup = {
        event: "click",
        target: "stackPopup",
        placement: "left",
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
                EditStack(stack.ID, text).then(() => {
                    dispatch("edit", {
                        stackID: stack.ID,
                        name: text,
                    })
                })
            },
        })
    }

    function deleteStack() {
        modal.trigger({
            type: "confirm",
            title: "Are you sure you want to delete this stack?",
            response: (state) => {
                if (state) {
                    DeleteStack(stack.ID).then(() => {
                        dispatch("delete", stack.ID)
                    })
                }
            },
        })
    }
</script>

<div class="flex flex-col gap-4 mx-8 max-h-screen overflow-auto pb-20">
    <div data-popup="stackPopup">
        <div class="btn-group-vertical variant-filled">
            <button on:click={editStack}>
                <EditIcon />
                <span>Edit</span>
            </button>
            <button on:click={deleteStack}>
                <DeleteIcon />
                <span>Delete</span>
            </button>
        </div>
        <div class="arrow variant-filled"></div>
    </div>
    <div class="flex justify-end">
        <button use:popup={stackPopup} class="btn-icon variant-filled">
            <MenuIcon />
        </button>
    </div>
    {#if tasksTodo.length}
        <h1 class="badge variant-filled w-min">Todo</h1>
        <ul
            use:dndzone={{ items: tasksTodo, dropTargetStyle: {} }}
            on:consider={handleDndConsider}
            on:finalize={handleDndFinalize}
            class="flex flex-col gap-4"
        >
            {#each tasksTodo as task (task.ID)}
                <Task
                    on:delete={taskDeleted}
                    on:edit={taskEdit}
                    on:done={taskDone}
                    {task}
                />
            {/each}
        </ul>
    {/if}
    {#if tasksDone.length}
        <h1 class="badge variant-filled w-min">Done</h1>
        <ul class="flex flex-col gap-4">
            {#each tasksDone as task (task.ID)}
                <Task on:delete={taskDeleted} on:done={taskDone} {task} />
            {/each}
        </ul>
    {/if}
    <button
        on:click={addTask}
        type="button"
        class="btn-icon btn-icon-lg variant-filled fixed bottom-8 right-8"
    >
        <AddIcon style="font-size: 1.25em;" />
    </button>
</div>
