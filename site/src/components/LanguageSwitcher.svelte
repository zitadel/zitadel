<script>
    import Select from 'svelte-select';
    const items = [
		{value: 'en', label: 'English', group: 'lang'},
        {value: 'de', label: 'Deutsch', group: 'lang'},
	];
    import { locale } from 'svelte-i18n';

    let language = 'en';
    $: setLocale(language);

    function setLocale(lang) {
        if (lang && lang.value) {
            locale.set(lang.value);
        }
    }

    $: language = items.find(item => item.value == $locale);
</script>

<style>
    :root {
        --purple: rgb(15, 16, 34);
        --highlight: #222544;
    }

    .language-switcher {
        display: inline-block;
        position: relative;
        --height: 36px;
        --indicatorHeight: 18px;
        --indicatorTop: 8px;
        --background: var(--back);
        --itemColor: white;
         --itemHoverBG: #3c4361;
        /* --itemHoverColor: white; */
        --itemIsActiveBG: #ffffff20;
        --itemIsActiveColor: #ffffff80;
        --listBackground: #2a2f45;
        --border: 1px solid #ffffff20;
        --borderHoverColor: var(--grey-text);
        --borderFocusColor: var(--grey-text);
        --borderRadius: 6px;
        min-width: 100px;
        font-size: 14px;
    }
</style>

<div class="language-switcher">
    <Select {items} bind:selectedValue={language} isClearable={false} isSearchable={false}></Select>
</div>
