<script>
    import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();
    let searchValue = '';
    // $: if (searchValue) {
    //     executeQuery(searchValue);
    // }

    $: executeQuery(searchValue);

    function executeQuery(value) {
        console.log(value);
    }

    function closeSearch() {
        console.log('clsoe search');
        dispatch('close', {
            closed: true,
        });
    }
</script>

<style>
    .overlay {
        position: fixed;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        background: #ffffff60;
        backdrop-filter: blur(10px);
    }

    .search-field {
        width: 100%;
        max-width: 500px;
        z-index: 100;
        box-shadow: 0 5px 10px rgba(0, 0, 0, .12);
        background-color: #2a2f45;
        border-radius: 8px;
        position: fixed;
        top: 30%;
        left: 50%;
        padding: 1rem;
        transform: translateX(-50%) translateY(-50%);
        display: relative;
    }

    .search-field .search-line {
        display: flex;
        align-items: center;
        border: 1px solid #ffffff10;
        border-radius: 8px;
    }

    .search-field .search-line i {
        margin: 0 1rem;
    }

    input {
        width: 100%;
        height: 45px;
        box-sizing: border-box;
        padding-inline-start: 10px;
        outline: none;
        display: inline-block;
        text-align: start;
        background-color: inherit;
        cursor: text;
        padding: 1px 20px;
        border-radius: 8px;
        margin: .5rem;
        transform: all .2 linear;
        font-size: 1.5rem;
        color: white;
        border: none;
    }

    input::placeholder {
        font-size: 14px;
        color: var(--grey-text);
        font-style: italic;
    }

    .result-list {
        display: flex;
        flex-direction: column;
    }

    .result-list .result-d {
        color: var(--grey);
        font-weight: 700;
        margin: 0;
        font-size: 1.3rem;
        margin: 5px 0;
    }

    .result-list .result-item {
        align-items: center;
        display: flex;
        justify-content: space-between;
        color: white;
        padding: 1rem;
        border-radius: 8px;
        margin: 2px 0;
    }

    .result-list .result-item:hover {
        background-color: #556cd6;
    }

    .result-list .result-item .title{
        margin: 0;
        font-size: 1.3rem;
        font-weight: 700;
    }

    .result-list .result-item .desc{
        color: var(--grey-text);
        margin: 0;
        font-size: 1.3rem;
    }
</style>

<div on:click="{closeSearch}" class="overlay"></div>
<div class="search-field">
    <div class="search-line">
        <i class="las la-search"></i>
        <input placeholder="Search for... " bind:value={searchValue}>
    </div>
    <div class="result-list">
    <p class="result-d">Found results: {searchValue}</p>
        <div class="result-item">
            <div>
                <p class="title">Title</p>
                <p class="desc">desc</p>
            </div>
            <i class="las la-link"></i>
        </div>
        <div class="result-item">
            <div>
                <p class="title">Title</p>
                <p class="desc">desc</p>
            </div>
            <i class="las la-link"></i>
        </div>
    </div>
</div>