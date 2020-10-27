<style>
    button {
        height: 20px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        color: var(--grey-text);
        margin: 1rem 2rem;
        background-color: #2a2f45;
        border: 1px solid #ffffff20;
        height: 35px;
        padding: 1rem;
        width: calc(var(--sidebar-w) - 4rem);
        box-sizing: border-box;
    }

    button:hover {
        background-color: #373d58;
    }

    button i {
        margin-right: 1rem;
    }

    button span {
        font-size: 1.3rem;
        margin-right: 1rem;
    }
</style>

<button on:click>
    <i class="las la-search"></i>
    <span>Search for anything</span>
</button>