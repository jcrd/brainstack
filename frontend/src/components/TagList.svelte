<script>
    import { Accordion, AccordionItem } from "@skeletonlabs/skeleton"

    import TagIcon from '~icons/material-symbols/label-outline-rounded'

    import Tag from "./Tag.svelte"

    export let tags
</script>

<Accordion>
    <AccordionItem open>
        <svelte:fragment slot="lead">
            <TagIcon />
        </svelte:fragment>
        <svelte:fragment slot="summary">
            <span>Tags</span>
        </svelte:fragment>
        <svelte:fragment slot="content">
            <div class="flex flex-wrap gap-1">
                {#each tags as tag}
                    <Tag {tag} />
                {/each}
            </div>
        </svelte:fragment>
    </AccordionItem>
</Accordion>
