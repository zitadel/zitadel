<script>
    import { createEventDispatcher } from 'svelte';
    import { _ } from 'svelte-i18n';

    export let sections;
    export let slug;

    let filteredResults = [];

    let selectedIndex = 0;

    $: selectIndex(selectedIndex);

    function init(el){
        el.focus()
    }

    function handleKeydown(event) {
        console.log(event)
        if (event) {
            if (event.keyCode == 37 || event.keyCode == 38) {
                if (selectedIndex > 0) {
                    event.preventDefault();
                    selectedIndex --;
                }
            }
            else if (event.keyCode == 39 || event.keyCode == 40) {
                event.preventDefault();
                selectedIndex ++;
            }
        }
    }

    function selectIndex(index) {
        console.log(index);
        const el = document.getElementById(index);
        if (el) {
            console.log('focus: '+ el);
            el.focus();
        }
    }

	const dispatch = createEventDispatcher();
    let searchValue = '';

    $: executeQuery(searchValue);

    function executeQuery(value) {
        const toSearchFor = value.toLowerCase();
        const filteredSections = sections.filter(section => {
            const slugContainsValue = section.slug.toLowerCase().includes(toSearchFor);
            const htmlContainsValue = section.html.replace(/<[^>]*>?/gm, '').toLowerCase().includes(toSearchFor);
            return slugContainsValue || htmlContainsValue;
        }).map(section => {
            // const removedHtml = section.html.replace(/<[^>]*>?/gm, '');
            // const foundIndex = removedHtml.indexOf(toSearchFor);
            // const subhtml = section.html.substring(foundIndex, (removedHtml.length - 1) - foundIndex > 150 ? 150 : removedHtml.length - 1)
            return {
                title: section.slug,
                slug: section.slug,
            }
        });

        const filteredSubSections = sections.map(section => {
            return section.subsections.map(sub => {
                return {parent: section.slug, ...sub};
            });
        }).flat().filter(subsection => {
            if (subsection.slug) {
                const slugContainsValue = subsection.slug.toLowerCase().includes(toSearchFor);
                const titleContainsValue = subsection.title.toLowerCase().includes(toSearchFor);
                return slugContainsValue || titleContainsValue;
            }
        });

        filteredResults = filteredSections.concat(filteredSubSections);
        console.log(filteredResults);
    }

    function closeSearch() {
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
        background: #ffffff30;
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
        top: 20%;
        left: 50%;
        padding: 1rem;
        transform: translateX(-50%);
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

    .result-d {
        color: var(--grey);
        font-weight: 700;
        margin: 0;
        font-size: 1.3rem;
        margin: 5px 0;
    }

    .result-list {
        max-height: 400px;
        overflow-y: auto;
    }

    .result-list .result-item {
        align-items: center;
        display: flex;
        justify-content: space-between;
        color: white;
        padding: 1rem;
        border-radius: 8px;
        margin: 2px 0;
        border-bottom: none;
        height: 77px;
        max-height: 80px;
    }
 
    .result-list .result-item:hover {
        background-color: #556cd680;
        box-sizing: border-box;
    } 

    .result-list .result-item:focus {
        background-color: #556cd6;
    }

    .text {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .result-list .result-item .title{
        margin: 0;
        font-size: 1.3rem;
        font-weight: 700;
    }

    .result-list .result-item .title .second-param {
        color: var(--grey-text);
        margin-left: 2rem;
    }

    .result-list .result-item .desc{
        color: var(--grey-text);
        margin: 0;
        font-size: 1.3rem;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
</style>

<svelte:window on:keydown={handleKeydown}/>

<div on:click="{closeSearch}" class="overlay"></div>
<div class="search-field">
    <div class="search-line">
        <i class="las la-search"></i>
        <input placeholder="{$_('search_input_placeholder')}" bind:value={searchValue} use:init>
    </div>
        <p class="result-d">{$_('search_results')}: </p>
    <div tabindex="-1" class="result-list">
        {#each filteredResults as result, i}
        <a tabindex="0" class="result-item" href="{slug}#{result.slug}" on:click="{closeSearch}" id="{i}">
            <div class="text">
            {#if result.level > 2}
                <p class="title">{result.title}<span class="second-param">{result.parent}</span></p>
            {:else}
                <p class="title">{result.slug}</p>
            {/if}
                <p class="desc" style="color: #85d996;">{slug}#{result.slug}</p>
                <p class="desc">{result.html || result.slug}</p>
            </div>
            <i class="las la-link"></i>
        </a>
        {/each}
    </div>
</div>