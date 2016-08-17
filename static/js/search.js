
function search() {
    if (searchText != '') {
        if (spells.length > 0) {
            spellSet = spells.filter(filterSearch);
            spellSet.sort(orderSpell);
        }
    } else {
        setSpellSet();
    }
    genResults(paginate(spellSet));
}

function filterSearch(doc) {
    var name = doc.name.toLowerCase();
    return name.indexOf(searchText) >= 0;
}
