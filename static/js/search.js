
function search() {
    page = 0;
    if (searchText != '') {
        if (spells.length > 0) {
            spellSet = spells.filter(filterSearch);
            spellSet.sort(orderSpell);
            $('span.letter').removeClass('label-primary');
            $('span.letter').addClass('label-default');
        }
    } else {
        var span = $('span[data-let="' + letter + '"]');
        span.removeClass('label-default');
        span.addClass('label-primary');
        setSpellSet();
    }
    genResults(paginate(spellSet));
}

function filterSearch(doc) {
    var name = doc.name.toLowerCase();
    return name.indexOf(searchText) >= 0;
}
