<script>
    import { tagSelections } from "../stores.js"

    export let tag
</script>

<div>
    {#if $tagSelections[tag.stack_id]}
        <button
            on:click={() => $tagSelections[tag.stack_id][tag.ID] = !$tagSelections[tag.stack_id][tag.ID]}
            class="chip {$tagSelections[tag.stack_id][tag.ID] ? 'variant-filled' : 'variant-soft'}"
        >
            #{tag.name}
        </button>
    {/if}
</div>
