<script>
    // @ts-nocheck

    import { GetStacks, CreateStack } from "../wailsjs/go/main/App"
    import { writable } from "svelte/store"
    import {
        initializeStores,
        Modal,
        getModalStore,
        TabGroup,
        Tab,
        storePopup,
    } from "@skeletonlabs/skeleton"
    import {
        computePosition,
        autoUpdate,
        offset,
        shift,
        flip,
        arrow,
    } from "@floating-ui/dom"
    import AddIcon from "~icons/material-symbols/add"
    import Stack from "./components/Stack.svelte"

    initializeStores()
    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow })

    const tabSet = writable(0)
    const stacks = writable([])
    const modal = getModalStore()

    GetStacks().then((result) => {
        $stacks = result.sort((a, b) => a.order - b.order)
    })

    function addStack() {
        modal.trigger({
            type: "prompt",
            title: "Add stack",
            valueAttr: { type: "text", required: true },
            response: (text) => {
                if (!text) {
                    return
                }
                const order = $stacks[$stacks.length - 1].order + 1
                CreateStack(text).then((result) => {
                    $stacks = [
                        ...$stacks,
                        {
                            ID: result,
                            name: text,
                            order,
                        },
                    ]
                })
            },
        })
    }

    function stackEdited({ detail: { stackID, name } }) {
        $stacks = $stacks.map((s) => {
            if (s.ID == stackID) {
                s.name = name
            }
            return s
        })
    }

    function stackDeleted({ detail: stackID }) {
        $stacks = $stacks.filter((s) => s.ID != stackID)
        $tabSet = 0
    }
</script>

<Modal />
<main class="flex h-screen w-full flex-col overflow-hidden">
    <TabGroup>
        {#if $stacks.length}
            {#each $stacks as stack, i}
                <Tab bind:group={$tabSet} name={stack.name} value={i}>
                    <span>{stack.name}</span>
                </Tab>
            {/each}
        {/if}
        <button on:click={addStack} type="button" class="btn-icon bg-initial">
            <AddIcon style="font-size: 1.25em;" />
        </button>
        <svelte:fragment slot="panel">
            {#if $stacks.length}
                <Stack on:edit={stackEdited} on:delete={stackDeleted} stack={$stacks[$tabSet]} {modal} />
            {/if}
        </svelte:fragment>
    </TabGroup>
</main>
