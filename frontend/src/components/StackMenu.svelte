<script>
    // @ts-nocheck

    import {
        EditStack,
        DeleteStack,
    } from "../../wailsjs/go/main/App"

    import { createEventDispatcher } from "svelte"

    import { popup } from "@skeletonlabs/skeleton"
    import MenuIcon from "~icons/material-symbols/menu-rounded"
    import DeleteIcon from "~icons/material-symbols/delete-forever-outline-rounded"
    import EditIcon from "~icons/material-symbols/edit"
    import PromoteIcon from '~icons/material-symbols/first-page'

    export let stack
    export let modal

    const dispatch = createEventDispatcher()

    const stackPopup = {
        event: "click",
        target: "stackPopup",
        placement: "bottom",
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

    function promoteStack() {
        if (stack.order !== 0) {
            dispatch("promote", stack.ID)
        }
    }
</script>

<div data-popup="stackPopup">
    <div class="btn-group-vertical variant-filled relative z-10">
        <button on:click={promoteStack} disabled={stack?.order === 0}>
            <PromoteIcon />
            <span>Promote</span>
        </button>
        <button on:click={editStack}>
            <div class="-ml-8">
                <EditIcon />
            </div>
            <span>Edit</span>
        </button>
        <button on:click={deleteStack}>
            <div class="-ml-4">
                <DeleteIcon />
            </div>
            <span>Delete</span>
        </button>
    </div>
    <div class="arrow variant-filled"></div>
</div>

<button use:popup={stackPopup} class="btn-icon">
    <MenuIcon style="font-size: 1.25em;" />
</button>
