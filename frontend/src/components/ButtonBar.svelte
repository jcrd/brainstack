<script>
    // @ts-nocheck

    import { popup } from "@skeletonlabs/skeleton"
    import AddIcon from "~icons/material-symbols/add"
    import MenuIcon from "~icons/material-symbols/menu-rounded"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import EditIcon from "~icons/material-symbols/edit"
    import DoneIcon from '~icons/material-symbols/done-rounded'

    export let data

    let { addTask, nextTask, editStack, deleteStack } = data

    const stackPopup = {
        event: "click",
        target: "stackPopup",
        placement: "left",
    }
</script>

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
<div class="flex justify-between">
    <button use:popup={stackPopup} class="btn-icon btn-icon-lg variant-filled">
        <MenuIcon />
    </button>
    <button on:click={nextTask} class="btn-icon btn-icon-lg variant-filled">
        <DoneIcon style="font-size: 1.25em;" />
    </button>
    <button on:click={addTask} class="btn-icon btn-icon-lg variant-filled">
        <AddIcon style="font-size: 1.25em;" />
    </button>
</div>
