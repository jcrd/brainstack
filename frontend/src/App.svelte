<script>
    // @ts-nocheck

    import { GetStack, GetStacks, CreateStack } from "../wailsjs/go/main/App"

    import { writable } from "svelte/store"

    import {
        initializeStores,
        Modal,
        Toast,
        getToastStore,
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
    import StackMenu from "./components/StackMenu.svelte"

    initializeStores()
    storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow })

    const tabSet = writable(0)
    const stacks = writable([])
    const modal = getModalStore()
    const toast = getToastStore()

    function catchError(error) {
        toast.trigger({
            message: error,
            background: "variant-filled-error",
        })
    }

    GetStacks()
        .then((result) => {
            $stacks = result.sort((a, b) => a.order - b.order)
        })
        .catch(catchError)

    $: if ($tabSet !== null && $stacks.length > 0) {
        const stack = $stacks[$tabSet]
        if (!stack.tasks || stack.invalid) {
            GetStack(stack.ID)
                .then((s) => {
                    s.invalid = false
                    $stacks[$tabSet] = s
                })
                .catch(catchError)
        }
    }

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
                CreateStack(text)
                    .then((result) => {
                        $stacks = [
                            ...$stacks,
                            {
                                ID: result,
                                name: text,
                                order,
                            },
                        ]
                        $tabSet = $stacks.length - 1
                    })
                    .catch(catchError)
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
        let idx = 0
        $stacks = $stacks.filter((s, i) => {
            if (s.ID == stackID) {
                idx = i
                return false
            } else {
                return true
            }
        })
        $tabSet = Math.min(idx, $stacks.length - 1)
    }

    function stackInvalidated({ detail: stackID }) {
        $stacks.forEach((s) => {
            if (s.ID === stackID) {
                s.invalid = true
                return
            }
        })
    }
</script>

<Toast />
<Modal />
<main class="flex h-screen w-full flex-col overflow-hidden">
    <TabGroup>
        <StackMenu
            on:edit={stackEdited}
            on:delete={stackDeleted}
            stack={$stacks[$tabSet]}
            {modal}
        />
        {#if $stacks.length}
            {#each $stacks as stack, i}
                <Tab id={`tab-stack-${i}`} bind:group={$tabSet} name={stack.name} value={i}>
                    <span>{stack.name}</span>
                </Tab>
            {/each}
        {/if}
        <button on:click={addStack} type="button" class="btn-icon">
            <AddIcon style="font-size: 1.25em;" />
        </button>
        <svelte:fragment slot="panel">
            {#if $stacks.length}
                <Stack
                    on:invalidate={stackInvalidated}
                    on:error={({ detail: error }) => catchError(error)}
                    stack={$stacks[$tabSet]}
                />
            {/if}
        </svelte:fragment>
    </TabGroup>
</main>
