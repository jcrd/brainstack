<script>
    import { tagSelections } from "../stores"

    export let tag
</script>

<button
    on:click={() => $tagSelections[tag.ID] = !$tagSelections[tag.ID]}
    class="chip {$tagSelections[tag.ID] ? 'variant-filled' : 'variant-soft'}"
>
    #{tag.name}
</button>
