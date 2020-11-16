<script>
    export let selected = false;
</script>
<style>
    button {
        background: inherit;
        border-radius: 8px;
        border: 1px solid var(--prime);
        box-shadow: 0 0 0 0 rgba(0,0,0,.2), 0 0 0 0 rgba(0,0,0,.14), 0 0 0 0 rgba(0,0,0,.12);
        padding: 0 15px;
        height: 40px;
        color: white;
        transition: background-color .2 ease;
        outline: none;
        min-width: 120px;
        cursor: pointer;
    }

    button[selected] {
        background-color: var(--prime);
        color: inherit;
    }

    button[selected]:hover {
        background-color: var(--prime);
        opacity: .9;
        color: inherit;
    }

    button[selected]:active {
        background-color: var(--prime);
        opacity: .8;
    }

    button:hover {
        background-color: #e91e6310;
    }
    button:active {
        background-color: #e91e6320;
    }

    button :global(span) {
        font-size: 14px;
        line-height: 14px;
        text-decoration: none;
    }
</style>

<button {selected}>            
    <slot />
</button> 
