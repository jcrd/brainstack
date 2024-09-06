<script>
    import DoneIcon from '~icons/material-symbols/check-circle'

    export let checked = false
    export let hover = true
    export let disabled = false
</script>

<div class="flex relative">
    <input
        type="checkbox"
        on:click
        bind:checked
        class="
            {hover ? 'hover:variant-soft' : ''}
            {disabled ? 'border-surface-400' : 'border-surface-700'}
            border-[1.5px]
            appearance-none peer w-4 h-4 rounded-xl"
    />
    <div class="absolute -ml-[1.5px] -mt-[1.5px] w-4 h-4 hidden peer-checked:block">
        <DoneIcon />
    </div>
</div>
