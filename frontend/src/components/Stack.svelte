<script>
    // @ts-nocheck

    import {
        AddTask,
        EditTask,
        ReorderTasks,
        UpdateTaskDone,
        EditStack,
        DeleteStack,
        GetTags,
    } from "../../wailsjs/go/main/App"

    import * as ConfigStore from "../../wailsjs/go/wailsconfigstore/ConfigStore"

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
    import TagList from "./TagList.svelte"

    import { parseTaskText } from "../lib.js"
    import { tagSelections } from "../stores.js"

    export let stack

    overrideItemIdKeyNameBeforeInitialisingDndZones("ID")

    const dispatch = createEventDispatcher()
    const tabSet = writable(0)

    $: {
        stack
        $tabSet = 0
    }

    $: tags = stack.tags || []

    let filteredTags = []
    let filteredTasks = []

    $: if (tags.length && $tagSelections[stack.ID]) {
        const tagIds = tags.map((t) => t.ID.toString())
        $tagSelections[stack.ID] = Object.fromEntries(
            Object.entries($tagSelections[stack.ID]).filter(([id]) => {
                return tagIds.includes(id)
        }))
    }

    $: {
        const stackTags = new Set(tags.map((t) => t.ID))
        const taskTags = tagsForTasks(tasks.filter((t) => $tabSet == 0 ? !t.done : t.done) || [])
        const intersection = stackTags.intersection(taskTags)
        filteredTags = tags.filter((t) => intersection.has(t.ID)) || []
    }

    $: tasks = stack.tasks?.sort((a, b) => a.order - b.order) || []
    $: {
        $tabSet
        $tagSelections[stack.ID]
        filteredTasks = tasks.filter(filterTags)
    }
    $: tasksDone = filteredTasks.filter((t) => t.done)
    $: tasksTodo = filteredTasks.filter((t) => !t.done)

    $: {
        tasks
        dispatch("invalidate", stack.ID)
    }

    $: if ($tagSelections[stack.ID]) {
        if (Object.keys($tagSelections[stack.ID]).length) {
            ConfigStore.Set(`stack.tag_selections`, JSON.stringify($tagSelections))
        }
    }

    function tagsForTasks(tasks) {
        let set = new Set()
        for (let task of tasks) {
            set = set.union(new Set(task.tags?.map((t) => t.ID)))
        }
        return set
    }

    function filterTags(task) {
        let hasSelection = false
        const tagIds = filteredTags.map((t) => t.ID.toString())
        for (const [id, state] of Object.entries($tagSelections[stack.ID])) {
            if (tagIds.includes(id) && state) {
                hasSelection = true
                break
            }
        }
        if (!hasSelection) {
            return true
        }
        let selected = 0
        for (const tag of task.tags) {
            selected += $tagSelections[stack.ID][tag.ID] ? 1 : 0
        }
        return selected > 0
    }

    function updateTags() {
        GetTags(stack.ID)
            .then((ts) => tags = ts)
            .catch((error) => dispatch("error", error))
    }

    function addTask({ detail: text }) {
        const order =
            tasks.length > 0 ? tasks[tasks.length - 1].order + 1 : 0
        const parsed = parseTaskText(text)

        AddTask(stack.ID, parsed, order)
            .then((task) => {
                tasks = [
                    ...tasks,
                    {
                        ...task,
                        order,
                    },
                ]
            })
            .catch((error) => dispatch("error", error))
            .finally(updateTags)
    }

    function taskEdit({ detail: { task, parsed } }) {
        EditTask(task.ID, parsed)
            .then((editedTask) => {
                tasks = tasks.map((t) => {
                    if (t.ID === editedTask.ID) {
                        return editedTask
                    }
                    return t
                })
            })
            .catch((error) => dispatch("error", error))
            .finally(updateTags)
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
        updateTags()
    }

    function taskDone({ detail: { taskID, done } }) {
        UpdateTaskDone(taskID, done)
            .then(() => {
                stack.tasks = tasks.map((t) => {
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

<div class="flex flex-col gap-4 mx-8 max-h-screen pb-20 p-1">
    <button class="pl-7 py-1 flex items-center gap-2 hover:bg-surface-200 hover:text-surface-600" on:click={() => $tabSet = !$tabSet}>
        {#if $tabSet == 0}
            <span class="mt-[3px]"><TabTodoIcon /></span>
            <span class="text-xl">Todo</span>
        {:else if $tabSet == 1}
            <span class="mt-[3px]"><TabDoneIcon /></span>
            <span class="text-xl">Done</span>
        {/if}
    </button>
    <div class="pl-3">
        <TagList tags={filteredTags} />
    </div>
    <div class="flex flex-col gap-4 overflow-y-auto overflow-x-hidden">
        {#if $tabSet == 0}
            {#if tasksTodo.length}
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
</div>
